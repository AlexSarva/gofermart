package utils

import (
	"errors"
	"strings"
)

var ErrNoCookie = errors.New("no cookie")

func ParseCookie(cookieStr string) (string, error) {
	cookieInfo := strings.Split(cookieStr, "; ")
	for _, pairs := range cookieInfo {
		elements := strings.Split(pairs, "=")
		if elements[0] == "session" {
			return elements[1], nil
		}
	}
	return "", ErrNoCookie
}
