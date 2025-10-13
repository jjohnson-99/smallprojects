package main
import (
    "net/http"
    "log"
    "fmt"
    "time"
    "sync/atomic"
    "encoding/json"
    "strings"
    "slices"
    "os"
    "database/sql"
    "github.com/joho/godotenv"
    "github.com/google/uuid"
    "github.com/jjohnson-99/chirpy/internal/database"
    "github.com/jjohnson-99/chirpy/internal/auth"
    _ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
    db *database.Queries
    platform string
    token_secret string
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	    cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerVisits(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
    content := `<html>
                    <body>
                        <h1>Welcome, Chirpy Admin</h1>
                        <p>Chirpy has been visited %d times!</p>
                    </body>
                </html>`

    fmt.Fprintf(w, content, cfg.fileserverHits.Load())
}

func (cfg *apiConfig) handlerResetVisits(w http.ResponseWriter, req *http.Request) {
    cfg.fileserverHits.Store(0)
}

func (cfg *apiConfig) handlerResetUsers(w http.ResponseWriter, req *http.Request) {
    if cfg.platform != "dev" {
        w.WriteHeader(403)
        return
    }
    err := cfg.db.DeleteUsers(req.Context())
    if err != nil {
        log.Printf("Error resetting users: %s", err)
        w.WriteHeader(500)
        return
    }
    w.WriteHeader(200)
    return
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, req *http.Request) {
    type parameters struct {
        UserEmail string `json:"email"`
        Password string `json:"password"`
    }

    decoder := json.NewDecoder(req.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
        log.Printf("Error decoding request: %s", err)
        w.WriteHeader(500)
        return
    }

    hashed_password, err := auth.HashPassword(params.Password)
    if err != nil {
        log.Printf("Error hashing password: %s", err)
        w.WriteHeader(500)
        return
    }

    arg := database.CreateUserParams{Email: params.UserEmail, HashedPassword: hashed_password}
    user, err := cfg.db.CreateUser(req.Context(),  arg)
    if err != nil {
        log.Println(err)
        log.Printf("Error creating user with email: %s", params.UserEmail)
        w.WriteHeader(500)
        return
    }
    
    respBody := database.User{ID: user.ID, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, Email: user.Email}
    respondWithJSON(w, 200, respBody)
}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, req *http.Request) {
    token_string, err := auth.GetBearerToken(req.Header)
    if err != nil {
        log.Printf("Access token malformed or missing.")
        w.WriteHeader(401)
        return
    }
    userID, err := auth.ValidateJWT(token_string, cfg.token_secret)
    if err != nil {
        log.Printf("Failed to validate JWT")
        w.WriteHeader(401)
        return
    }

    type parameters struct {
        Email string `json:"email"`
        Password string `json:"password"`
    }
 
    decoder := json.NewDecoder(req.Body)
    params := parameters{}
    err = decoder.Decode(&params)
    if err != nil {
        log.Printf("Error decoding request: %s", err)
        w.WriteHeader(500)
        return
    }

    hashed_password, err := auth.HashPassword(params.Password)
    if err != nil {
        log.Printf("Error hashing password: %s", err)
        w.WriteHeader(500)
        return
    }
    
    arg := database.UpdateUserParams{ID: userID, Email: params.Email, HashedPassword: hashed_password}
    user, err := cfg.db.UpdateUser(req.Context(), arg)

    respBody := database.User{ID: user.ID, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, Email: user.Email}
    respondWithJSON(w, 200, respBody)
}

func (cfg *apiConfig) handlerUpgradeUser(w http.ResponseWriter, req *http.Request) {
    type parameters struct {
        Event string `json:"event"`
        Data struct {
            UserID uuid.UUID `json:"user_id"`
        } `json:"data"`
    }
 
    decoder := json.NewDecoder(req.Body)
    params := parameters{}
    err := decoder.Decode(&params)
    if err != nil {
        log.Printf("Error decoding request: %s", err)
        w.WriteHeader(500)
        return
    }

    if params.Event != "user.upgrade" {
        log.Printf("Event must be user.upgrade, got %s", params.Event)
        w.WriteHeader(204)
        return
    }

    err = cfg.db.UpgradeChirpy(req.Context(), params.Data.UserID)
    if err != nil {
        log.Print("User could not be found")
        w.WriteHeader(404)
        return
    }

    w.WriteHeader(204)
    return
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
    type parameters struct {
        Email string `json:"email"`
        Password string `json:"password"`
    }

    decoder := json.NewDecoder(req.Body)
    params := parameters{}  
    err := decoder.Decode(&params)
    if err != nil {
        log.Printf("Error decoding request: %s", err)
        w.WriteHeader(500)
        return
    }

    user, err := cfg.db.GetUserFromEmail(req.Context(), params.Email)
    if err != nil {
        w.WriteHeader(401)
        w.Write([]byte("Incorrect email or password"))
        return
    }

    err = auth.CheckPasswordHash(user.HashedPassword, params.Password)
    if err != nil {
        w.WriteHeader(401)
        w.Write([]byte("Incorrect email or password"))
        return
    }  

    token, err := auth.MakeJWT(user.ID, cfg.token_secret, time.Hour)
    if err != nil {
        log.Printf("Error creating new JWT for user: %s", user.ID)
    }

    refresh_token, err := auth.MakeRefreshToken()
    arg := database.CreateRefreshTokenParams{Token: refresh_token, UserID: user.ID}
    _, err = cfg.db.CreateRefreshToken(req.Context(), arg)
    if err != nil {
        log.Println(err)
        log.Printf("Error creating refresh_token with UserID: %s", user.ID)
        w.WriteHeader(500)
        return
    }
    
    respBody := struct {
                    ID uuid.UUID
                    CreatedAt time.Time 
                    UpdatedAt time.Time
                    Email string
                    Token string
                    Refresh_Token string
                }{
                    ID: user.ID,
                    CreatedAt: user.CreatedAt,
                    UpdatedAt: user.UpdatedAt,
                    Email: user.Email,
                    Token: token,
                    Refresh_Token: refresh_token,
                }

    respondWithJSON(w, 200, respBody)
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, req *http.Request) {
    token_string, err := auth.GetBearerToken(req.Header)
    if err != nil {
        log.Printf("No token is present in request Header")
        w.WriteHeader(401)
        return
    }
    refresh_token, err := cfg.db.GetRefreshToken(req.Context(), token_string)
    if err != nil {
        log.Println(err)
        log.Printf("No User with token: %s exists", token_string)
        w.WriteHeader(401)
        return
    }
    if time.Now().After(refresh_token.ExpiresAt) {
        log.Println(err)
        log.Printf("Token: %s has expired", refresh_token)
        w.WriteHeader(401)
        return
    }

    token, err := auth.MakeJWT(refresh_token.UserID, cfg.token_secret, time.Hour)
    
    respBody := struct{token string}{token: token}
    respondWithJSON(w, 200, respBody)
}


func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, req *http.Request) {
    token_string, err := auth.GetBearerToken(req.Header)
    if err != nil {
        log.Printf("No token is present in request Header")
        w.WriteHeader(401)
        return
    }
    refresh_token, err := cfg.db.GetRefreshToken(req.Context(), token_string) 
    if err != nil {
        log.Println(err)
        log.Printf("Refresh token: %s does not exist", token_string)
        w.WriteHeader(401)
        return
    }
    
    err = cfg.db.RevokeRefreshToken(req.Context(), refresh_token.Token) 
    w.WriteHeader(204)
    return
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, req *http.Request) {
    type parameters struct {
        Body string `json:"body"`
        UserID  uuid.UUID `json:"user_id"`
    }

    token_string, err := auth.GetBearerToken(req.Header)
    if err != nil {
        log.Printf("Unauthorized")
        w.WriteHeader(401)
        return
    }
    _, err = auth.ValidateJWT(token_string, cfg.token_secret)
    if err != nil {
        log.Printf("Unauthorized")
        w.WriteHeader(401)
        return
    }

    decoder := json.NewDecoder(req.Body)
    params := parameters{}
    err = decoder.Decode(&params)
    if err != nil {
        log.Printf("Error decoding request: %s", err)
        w.WriteHeader(500)
        return
    }

    if len(params.Body) > 140 {respondWithError(w, 400, "Chrip is too long"); return}

    arg := database.CreateChirpParams{Body: cleanWords(params.Body), UserID: params.UserID}
    chirp, err := cfg.db.CreateChirp(req.Context(), arg)
    if err != nil {
        log.Println(err)
        log.Printf("Error creating chirp: %s for user: %s", params.Body, params.UserID)
        w.WriteHeader(500)
        return
    }

    respondWithJSON(w, 200, chirp)
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, req *http.Request) {
    chirps, err := cfg.db.GetChirps(req.Context()) 
    if err != nil {
        log.Println(err)
        log.Printf("Error trying to get chirps")
        w.WriteHeader(500)
        return
    }

    respondWithJSON(w, 200, chirps)

}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, req *http.Request) {
    id, err := uuid.Parse(req.PathValue("id"))
    if err != nil {
        log.Println(err)
    }

    chirp, err := cfg.db.GetChirp(req.Context(), id)
    if err != nil {
        log.Println(err)
        log.Printf("Error trying to get chirp with ID: %d", id)
        w.WriteHeader(500)
        return
    }

    respondWithJSON(w, 200, chirp)
}

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, req *http.Request) { 
    token_string, err := auth.GetBearerToken(req.Header)
    if err != nil {
        log.Printf("Error getting token")
        w.WriteHeader(401)
        return
    }
    user_id, err := auth.ValidateJWT(token_string, cfg.token_secret)
    if err != nil {
        log.Printf("Error validating JWT")
        w.WriteHeader(401)
        return
    }

    chirp_id, err := uuid.Parse(req.PathValue("chirpID"))
    if err != nil {
        log.Println(err)
    }

    chirp, err := cfg.db.GetChirp(req.Context(), chirp_id)
    if err != nil {
        log.Println(err)
        log.Printf("Could not find chirp with ID: %d", chirp_id)
        w.WriteHeader(404)
        return
    }

    if chirp.UserID != user_id {
        log.Printf("User: %s is not authorized to delete this chirp", user_id)
        w.WriteHeader(403)
        return
    }
     
    err = cfg.db.DeleteChirp(req.Context(), chirp.ID) 
    w.WriteHeader(204)
    return
}

func cleanWords(msg string) string {
    profaneList := []string{"kerfuffle", "sharbert", "fornax"}
    words := strings.Split(msg, " ")   
    //var sb strings.Builder
    for i, word := range words {
        if slices.Contains(profaneList, strings.ToLower(word)) {
            words[i] = "****"
        }
    }
    //return sb.String()
    return strings.Join(words, " ")
}

func respondWithError(w http.ResponseWriter, code int, msg string) { 
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    w.WriteHeader(code)
    w.Write([]byte(msg))
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) { 
    dat, err := json.Marshal(payload)
    if err != nil {
        log.Printf("Error marshalling JSON: %s", err)
        w.WriteHeader(500)
        return
    }   
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(dat)
}

func main() {
    godotenv.Load() 
    port := "8080"
    var cfg apiConfig 
    cfg.platform = os.Getenv("PLATFORM")
    cfg.token_secret = os.Getenv("SECRET")

    dbURL := os.Getenv("DB_URL")
    db, err := sql.Open("postgres", dbURL)
    if err != nil {
        log.Fatal(err)
    }
    dbQueries := database.New(db)
    cfg.db = dbQueries


    mux := http.NewServeMux()
    mux.Handle("/app/", cfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(".")))))

    mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)   
        w.Write([]byte("OK"))
    })

    mux.HandleFunc("GET /admin/metrics", cfg.handlerVisits)
    mux.HandleFunc("GET /api/chirps", cfg.handlerGetChirps)
    mux.HandleFunc("GET /api/chirps/{id}", cfg.handlerGetChirp)
    mux.HandleFunc("POST /admin/reset", cfg.handlerResetUsers)
    mux.HandleFunc("POST /api/chirps", cfg.handlerCreateChirp)
    mux.HandleFunc("POST /api/users", cfg.handlerCreateUser)
    mux.HandleFunc("POST /api/login", cfg.handlerLogin)
    mux.HandleFunc("POST /api/refresh", cfg.handlerRefresh)
    mux.HandleFunc("POST /api/revoke", cfg.handlerRefresh)
    mux.HandleFunc("POST /api/revoke", cfg.handlerRefresh)
    mux.HandleFunc("POST /api/polka/webhooks", cfg.handlerUpgradeUser)
    mux.HandleFunc("PUT /api/users", cfg.handlerUpdateUser)
    mux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.handlerDeleteChirp)

    srv := &http.Server{
        Addr:    ":" + port,
        Handler: mux,
    }

    log.Fatal(srv.ListenAndServe())
}
