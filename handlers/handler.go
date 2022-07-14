package handlers

import (
	"AlexSarva/gofermart/internal/app"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// Дополнительный обработчик ответа
func messageResponse(w http.ResponseWriter, message, ContentType string, httpStatusCode int) {
	w.Header().Set("Content-Type", ContentType)
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

// Обработка сжатых запросов
// TODO вынести в middleware
func readBodyBytes(r *http.Request) (io.ReadCloser, error) {
	// GZIP decode
	if len(r.Header["Content-Encoding"]) > 0 && r.Header["Content-Encoding"][0] == "gzip" {
		// Read body
		bodyBytes, readErr := ioutil.ReadAll(r.Body)
		if readErr != nil {
			return nil, readErr
		}
		defer r.Body.Close()

		newR, gzErr := gzip.NewReader(ioutil.NopCloser(bytes.NewBuffer(bodyBytes)))
		if gzErr != nil {
			log.Println(gzErr)
			return nil, gzErr
		}
		defer newR.Close()

		return newR, nil
	} else {
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

func MyHandler(database *app.Database) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	//r.Use(CookieHandler)
	r.Use(middleware.AllowContentEncoding("gzip"))
	r.Use(middleware.AllowContentType("application/json", "text/plain", "application/x-gzip"))
	r.Use(middleware.Compress(5, gzipContentTypes))
	r.Post("/api/user/register", UserRegistration(database))
	r.Post("/api/user/login", UserAuthentication(database))
	r.Post("/api/user/orders", PostOrder(database))
	r.Get("/api/user/orders", GetOrders(database))
	r.Get("/api/user/balance", GetBalance(database))
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
