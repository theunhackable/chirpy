package auth

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var secret = "my-secret"
var randUserId, _ = uuid.NewRandom()

func TestValidJwt(t *testing.T) {
	token, err := MakeJWT(randUserId, secret, 5*time.Second)

	if err != nil {
		t.Errorf("something went wrong while making jwt token: %v", err)
		return
	}

	println("token: ", token)

	uuid, err := ValidateJWT(token, secret)

	if err != nil {
		t.Errorf("something went wrong while validating jwt token: %v", err)
		return
	}

	if uuid == randUserId {
		fmt.Println("valid jwt is success")
	}

}

func TestExpiredToken(t *testing.T) {
	token, err := MakeJWT(randUserId, secret, time.Microsecond)

	if err != nil {
		t.Errorf("something went wrong while making jwt token: %v", err)
		return
	}

	_, validJWTErr := ValidateJWT(token, secret)

	if !errors.Is(validJWTErr, jwt.ErrTokenExpired) {
		t.Error("Token should be expired but haven't")
		return
	}

}
