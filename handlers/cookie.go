package handlers

import (
	"AlexSarva/gofermart/crypto"
	"github.com/google/uuid"
	"log"
	"net/http"
	"time"
)

func GenerateCookie(userID uuid.UUID) http.Cookie {
	session := crypto.Encrypt(userID, crypto.SecretKey)
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{Name: "session", Value: session, Expires: expiration, Path: "/"}
	return cookie
}

func getCookie(r *http.Request) (uuid.UUID, error) {
	cookie, cookieErr := r.Cookie("session")
	if cookieErr != nil {
		log.Println(cookieErr)
		return uuid.UUID{}, ErrNotValidCookie
	}
	userID, cookieDecryptErr := crypto.Decrypt(cookie.Value, crypto.SecretKey)
	if cookieDecryptErr != nil {
		return uuid.UUID{}, cookieDecryptErr
	}
	return userID, nil

}
