package auth

import (
	"testing"
    "time"
    "net/http"
    "github.com/google/uuid"
)

func TestBearerToken(t *testing.T) {
    userID := uuid.New()
    token, _ := MakeJWT(userID, "12345", time.Duration(time.Second))

    req, err := http.NewRequest("POST", "http://example.com", nil)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
    req.Header.Add("authorization", "Bearer " + token)

    token_string, err := GetBearerToken(req.Header)
    if err != nil {
		t.Fatalf("err: %v", err)
	}   

    if token_string != token {
            t.Errorf("actual: '%s' vs expected: '%s'", token_string, token)
    }
}

