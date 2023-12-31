package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetApiKey (headers http.Header) (string, error) {
	key := headers.Get("Authorization")
	
	if key == "" {
		return "", errors.New("no API Key detected")
	}

	keys := strings.Split(key, " ")
	if len(keys) != 2 {
		return "", errors.New("malformed auth header")
	}

	if keys[0] != "RSS_KEY" {
		return "", errors.New("malformed first part of auth header")
	}

	return keys[1], nil

}