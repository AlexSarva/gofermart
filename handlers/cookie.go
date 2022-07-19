package handlers

import (
	"AlexSarva/gofermart/crypto"
	"errors"
	"github.com/google/uuid"
	"net/http"
	"strings"
	"time"
)

var ErrNotValidCookie = errors.New("valid cookie does not found")
var ErrNoAuth = errors.New("no Bearer token")

func GenerateCookie(userID uuid.UUID) (http.Cookie, time.Time) {
	session := crypto.Encrypt(userID, crypto.SecretKey)
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{Name: "session", Value: session, Expires: expiration, Path: "/"}
	return cookie, expiration
}

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
