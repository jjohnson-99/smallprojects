package auth
import (
    "net/http"
    "errors"
    "strings"
    "crypto/rand"
    "encoding/hex"
    //"fmt"
    //"time"
    //"github.com/google/uuid"
    //"golang.org/x/crypto/bcrypt"
    //"github.com/golang-jwt/jwt/v5"
)

func GetBearerToken(headers http.Header) (string, error) {
   token_string := headers.Get("authorization")
    if token_string == "" {
        return "", errors.New("The authorization header does not exist")
    }
    token_string = strings.Split(token_string, " ")[1]

    return token_string, nil
}

func MakeRefreshToken() (string, error) {
    var src []byte
    rand.Read(src)
    refresh_token := hex.EncodeToString(src)

    return refresh_token, nil
}
