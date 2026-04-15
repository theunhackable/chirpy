package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	bearer := headers.Get("Authorization")

	if bearer == "" {
		return "", errors.New("API token not found.")
	}
	splits := strings.Split(bearer, " ")

	if len(splits) != 2 {

		return "", errors.New("Invalid API token.")
	}

	return splits[1], nil

}
