package auth
import (
    "fmt"
    "time"
    "github.com/google/uuid"
    "golang.org/x/crypto/bcrypt"
    "github.com/golang-jwt/jwt/v5"
)

func HashPassword(password string) (string, error) {
    passwordByte := []byte(password)

    // Hashing the password with the default cost of 10
    hashedPassword, err := bcrypt.GenerateFromPassword(passwordByte, bcrypt.DefaultCost)
    if err != nil {
        panic(err)
    }
    
    return string(hashedPassword), err
}

func CheckPasswordHash(hash, password string) error {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
    claims := jwt.RegisteredClaims{
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
        IssuedAt:  jwt.NewNumericDate(time.Now()),
        Issuer:    "chirpy",
        Subject:   userID.String(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    ss, err := token.SignedString([]byte(tokenSecret))
    return ss, err
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
    var id uuid.UUID
    claims := jwt.RegisteredClaims{}
    token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
	    return []byte(tokenSecret), nil
    })
    if err != nil {
        fmt.Print(err)
        return id, err
    } 

    s, err := token.Claims.GetSubject()
    if err != nil {
        fmt.Print(err)
        return id, err
    }

    id, err = uuid.Parse(s)
    return id, err
}
