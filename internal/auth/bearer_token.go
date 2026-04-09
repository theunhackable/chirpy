package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	bearer := headers.Get("Authorization")

	if bearer == "" {
		return "", errors.New("Bearer token not found.")
	}
	splits := strings.Split(bearer, " ")

	if len(splits) != 2 {

		return "", errors.New("Invalid bearer token.")
	}

	return splits[1], nil

}
