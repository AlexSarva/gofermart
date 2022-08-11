package handlers

import (
	"AlexSarva/gofermart/crypto"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

// ErrNotValidCookie error while valid cookie doesn't contain in request Header
var ErrNotValidCookie = errors.New("valid cookie does not found")

// ErrNoAuth error while valid Bearer token doesn't contain in request Header
var ErrNoAuth = errors.New("no Bearer token")

// ErrNoCookie error that occurs when no cookie presents in Header
var ErrNoCookie = errors.New("no cookie")

// GenerateCookie function of generating cookie for user when he successfully registered and authenticated
// based at UserID (uuid format)
// returns Cookie format for respond and time of expiration
func GenerateCookie(userID uuid.UUID) (http.Cookie, time.Time) {
	session := crypto.Encrypt(userID, crypto.SecretKey)
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{Name: "session", Value: session, Expires: expiration, Path: "/"}
	return cookie, expiration
}

// GetCookie cookie selection function from Header
// returns UserID in uuid format
func GetCookie(r *http.Request) (uuid.UUID, error) {
	cookie, cookieErr := r.Cookie("session")
	if cookieErr != nil {
		return uuid.UUID{}, ErrNotValidCookie
	}
	userID, cookieDecryptErr := crypto.Decrypt(cookie.Value, crypto.SecretKey)
	if cookieDecryptErr != nil {
		return uuid.UUID{}, cookieDecryptErr
	}
	return userID, nil

}

// GetToken cookie selection function from Header
// returns UserID in uuid format
func GetToken(r *http.Request) (uuid.UUID, error) {
	auth := r.Header.Get("Authorization")
	if len(auth) == 0 {
		return uuid.UUID{}, ErrNoAuth
	}
	tokenValue := strings.Split(auth, "Bearer ")
	if len(tokenValue) < 2 {
		return uuid.UUID{}, ErrNoAuth
	}
	authToken := tokenValue[1]
	userID, tokenDecryptErr := crypto.Decrypt(authToken, crypto.SecretKey)
	if tokenDecryptErr != nil {
		return uuid.UUID{}, tokenDecryptErr
	}
	return userID, nil
}

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
