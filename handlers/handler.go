package handlers

import (
	"AlexSarva/gofermart/crypto"
	"AlexSarva/gofermart/internal/app"
	"bytes"
	"compress/gzip"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var ErrNotValidCookie = errors.New("valid cookie does not found")

// Обработка сжатых запросов
func readBodyBytes(r *http.Request) (io.ReadCloser, error) {
	// GZIP decode
	if len(r.Header["Content-Encoding"]) > 0 && r.Header["Content-Encoding"][0] == "gzip" {
		// Read body
		bodyBytes, readErr := ioutil.ReadAll(r.Body)
		if readErr != nil {
			return nil, readErr
		}
		defer r.Body.Close()

		log.Println("compressed request")

		newR, gzErr := gzip.NewReader(ioutil.NopCloser(bytes.NewBuffer(bodyBytes)))
		if gzErr != nil {
			log.Println(gzErr)
			return nil, gzErr
		}
		defer newR.Close()

		return newR, nil
	} else {
		log.Println("no compressed request")
		return r.Body, nil
	}
}

func PingDB(database *app.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ping := database.Repo.Ping()
		if ping {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

	}
}

var gzipContentTypes = "application/x-gzip, application/javascript, application/json, text/css, text/html, text/plain, text/xml"

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

func CookieHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		_, userIDErr := getCookie(r)
		if userIDErr != nil {
			log.Println(userIDErr)
			userCookie := GenerateCookie(uuid.New())
			log.Println(userCookie)
			r.AddCookie(&userCookie)
			http.SetCookie(w, &userCookie)
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func MyHandler(database *app.Database) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(CookieHandler)
	r.Use(middleware.AllowContentEncoding("gzip"))
	r.Use(middleware.AllowContentType("application/json", "text/plain", "application/x-gzip"))
	r.Use(middleware.Compress(5, gzipContentTypes))
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})
	r.Get("/ping", PingDB(database))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, nfErr := w.Write([]byte("route does not exist"))
		if nfErr != nil {
			log.Println(nfErr)
		}
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, naErr := w.Write([]byte("sorry, only GET, POST and DELETE methods are supported."))
		if naErr != nil {
			log.Println(naErr)
		}
	})
	return r
}
