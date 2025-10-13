package auth

import (
	"testing"
    "time"
    "github.com/google/uuid"
)

func TestValidatePassword(t *testing.T) {
    userID := uuid.New()
    token, _ := MakeJWT(userID, "12345", time.Duration(time.Second))
    type input struct {
        id uuid.UUID
        tokenString string
        tokenSecret string
    }
    type casesStruct []struct {
		input input
		expected uuid.UUID
	}
    cases := casesStruct{
		{
            input: input{id: userID, tokenString: token, tokenSecret: "12345"},
			expected: userID,
		},
	}

	for _, c := range cases {
		actual, _ := ValidateJWT(token, c.input.tokenSecret)
		if actual != c.expected {
            t.Errorf("actual: '%s' vs expected: '%s'", actual, c.expected)
			continue
		}
	}
}

func TestValidateExpiredToken(t *testing.T) {
    userID := uuid.New()
    token, _ := MakeJWT(userID, "12345", time.Duration(time.Millisecond))
    type input struct {
        id uuid.UUID
        tokenString string
        tokenSecret string
    }
    type casesStruct []struct {
		input input
		expected uuid.UUID
	}
    cases := casesStruct{
		{
            input: input{id: userID, tokenString: token, tokenSecret: "12345"},
			expected: userID,
		},
	}

	for _, c := range cases {
        time.Sleep(time.Second)
		_, err := ValidateJWT(token, c.input.tokenSecret)
		if err == nil {
            t.Error("expected error")
			continue
		}
	}
}

func TestValidatePasswordError(t *testing.T) {
    userID := uuid.New()
    token, _ := MakeJWT(userID, "12345", time.Duration(time.Second))
    type input struct {
        id uuid.UUID
        tokenString string
        tokenSecret string
    }
    type casesStruct []struct {
		input input
		expected uuid.UUID
	}
    cases := casesStruct{
		{
            input: input{id: userID, tokenString: token, tokenSecret: "12354"},
			expected: userID,
		},
	}

	for _, c := range cases {
		_, err := ValidateJWT(token, c.input.tokenSecret)
		if err == nil {
            t.Error("expected error")
			continue
		}
	}
}

