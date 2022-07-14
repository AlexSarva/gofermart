package handlers

import (
	"AlexSarva/gofermart/crypto"
	"errors"
	"github.com/google/uuid"
	"net/http"
	"time"
)

var ErrNotValidCookie = errors.New("valid cookie does not found")

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
