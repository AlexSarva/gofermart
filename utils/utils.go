package utils

import (
	"errors"
	"strings"
)

// ErrNoCookie error that occurs when no cookie presents in Header
var ErrNoCookie = errors.New("no cookie")

// ParseCookie util that parse cookie string format into session id
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
