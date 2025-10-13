package main
import (
    "fmt"
    "os"
    "errors"
    "context"
    "time"
    "log"
    "io"
    "encoding/xml"
    "html"
    "net/http"
    "github.com/google/uuid"
    "github.com/jjohnson-99/gator/internal/config"
    "github.com/jjohnson-99/gator/internal/database"
)

type state struct {
    db  *database.Queries
    cfg *config.Config
}

type command struct {
    name string
    args []string
}

type cliCommands struct {
    commands map[string]func(*state, command) error
}

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func (c *cliCommands) register(name string, f func(*state, command) error) {
    c.commands[name] = f
}

func (c *cliCommands) run(s *state, cmd command) error {
    f, ok := c.commands[cmd.name]
    if !ok {
        return errors.New("Unknown command.")
    }

    f(s, cmd)
    return nil
}

func handlerLogin(s *state, cmd command) error {
    if len(cmd.args) == 0 {
        return errors.New("Login command expects a single argument")
    }

    _, err := s.db.GetUser(context.Background(), cmd.args[0])
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    s.cfg.SetUser(cmd.args[0]) 
    fmt.Printf("User %s has been set.\n", cmd.args[0])
    return nil
}

func handlerRegister(s *state, cmd command) error {
    if len(cmd.args) == 0 { 
        return errors.New("Register command expects a single argument")
    }
    
    arg := database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: cmd.args[0]}  
    _, err := s.db.CreateUser(context.Background(), arg)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    s.cfg.SetUser(cmd.args[0])

    fmt.Printf("User %s was created.\n", cmd.args[0])
    return nil
}

func handlerReset(s *state, cmd command) error {
    err := s.db.ResetUsers(context.Background())
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    return nil
}

func handlerGetUsers(s *state, cmd command) error {
    users, err := s.db.GetUsers(context.Background()) 
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    for _, user := range users {
        if user == s.cfg.Current_user_name {
            fmt.Printf("* %s (current)\n", user)
        } else {
            fmt.Printf("* %s\n", user)
        }
    }

    return nil
}

func handlerGetFeeds(s *state, cmd command) error {
    feeds, err := s.db.GetFeeds(context.Background())
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    for _, feed := range feeds {
        fmt.Printf("Feed name: %s, Feed url: %s, User name: %s\n", feed.Name, feed.Url, feed.Name_2)
    }

    return nil
}

func handlerAggregator(s *state, cmd command) error {
    c := "https://www.wagslane.dev/index.xml"
    feed, err := fetchFeed(context.Background(), c)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(feed)
    return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error { 
    if len(cmd.args) < 2 { 
        os.Exit(1)
    }
    //user, err := s.db.GetUser(context.Background(), s.cfg.Current_user_name) 
    arg := database.CreateFeedParams{
        ID: uuid.New(),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        Name: cmd.args[0],
        Url: cmd.args[1],
        UserID: user.ID}  
        _, err := s.db.CreateFeed(context.Background(), arg)
        if err != nil {
            os.Exit(1)
        }

    fmt.Println("debug")
    handlerFollow(s, command{args: []string{cmd.args[1]}}, user)
    return nil
}

func handlerFollow (s *state, cmd command, user database.User) error {
    if len(cmd.args) < 1 {
        os.Exit(1)
    }
    feed, err := s.db.GetFeed(context.Background(), cmd.args[0])
    if err != nil {
        return err
    }
    arg := database.CreateFeedFollowParams{
        ID: uuid.New(),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
        UserID: user.ID,
        FeedID: feed.ID}
    _, err = s.db.CreateFeedFollow(context.Background(), arg) 
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    fmt.Printf("User name: %s, Feed name: %s\n", user.Name, feed.Name)
    return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
    follows, err := s.db.GetFeedFollowsForUser(context.Background(), user.Name)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("User %s is following:\n", user.Name)
    for _, follow := range follows {
        fmt.Println(follow.Name_2)
    }

    return nil
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
    var client http.Client

    req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
    if err != nil {
        log.Fatal(err)
    }
    req.Header.Set("User-Agent", "gator")

    res, err := client.Do(req)
    if err != nil {
        log.Fatal(err)
    }
    defer res.Body.Close()
    body, err := io.ReadAll(res.Body)

    if res.StatusCode > 299 {
        log.Fatalf("Response failed with status code: %d and \nbody: %s\n", res.StatusCode, body)
    }
    
    var feed RSSFeed
    if err := xml.Unmarshal(body, &feed); err != nil {
        return nil, err
    }

    feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
    feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
    for _, item := range feed.Channel.Item {
        item.Title = html.UnescapeString(item.Title)
        item.Description = html.UnescapeString(item.Description)
    }
            
    return &feed, nil
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
    return func(s *state, cmd command) error {
        user, err := s.db.GetUser(context.Background(), s.cfg.Current_user_name)
        if err != nil {
            log.Fatal(err)
        }
        handler(s, cmd, user)
        return nil
    }
}
