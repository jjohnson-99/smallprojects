package main
import (
    "fmt"
    "os"
    "log"
    "database/sql"
    "github.com/jjohnson-99/gator/internal/config"
    "github.com/jjohnson-99/gator/internal/database"
    _ "github.com/lib/pq"
)

func main() {
    var cfg config.Config
    cfg = config.Read()
    s := state{cfg: &cfg}

    dbURL := cfg.DB_URL
    db, err := sql.Open("postgres", dbURL) 
    if err != nil {
        log.Fatal(err)
    }
    dbQueries := database.New(db)
    s.db = dbQueries

    commands := cliCommands{commands: make(map[string]func(*state, command) error)}
    commands.register("login", handlerLogin)
    commands.register("register", handlerRegister)
    commands.register("reset", handlerReset)
    commands.register("users", handlerGetUsers)
    commands.register("agg", handlerAggregator)
    commands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
    commands.register("feeds", handlerGetFeeds)
    commands.register("follow", middlewareLoggedIn(handlerFollow))
    commands.register("following", middlewareLoggedIn(handlerFollowing))

    args := os.Args
    if len(args) < 2 {
        fmt.Println("At least two arguments expected.")
        os.Exit(1)
    }
    cmd := command{name: args[1], args: args[2:]}
    commands.run(&s, cmd)
}
