package handlers

import (
	"AlexSarva/gofermart/internal/app"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// messageResponse additional respond generator
// useful in case of error handling in outputting results to respond
func messageResponse(w http.ResponseWriter, message, ContentType string, httpStatusCode int) {
	w.Header().Set("Content-Type", ContentType)
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, jsonRespErr := json.Marshal(resp)
	if jsonRespErr != nil {
		log.Println(jsonRespErr)
	}
	w.Write(jsonResp)
}

// readBodyBytes compressed request processing function
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

// gzipContentTypes request types that support data compression
var gzipContentTypes = "application/x-gzip, application/javascript, application/json, text/css, text/html, text/plain, text/xml"

// MyHandler - the main handler of the server
// contains middlewares and all routes
func MyHandler(database *app.Database) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.AllowContentEncoding("gzip"))
	r.Use(middleware.AllowContentType("application/json", "text/plain", "application/x-gzip"))
	r.Use(middleware.Compress(5, gzipContentTypes))
	r.Mount("/debug", middleware.Profiler())
	r.Post("/api/user/register", UserRegistration(database))
	r.Post("/api/user/login", UserAuthentication(database))
	r.Post("/api/user/orders", PostOrder(database))
	r.Get("/api/user/orders", GetOrders(database))
	r.Get("/api/user/balance", GetBalance(database))
	r.Post("/api/user/balance/withdraw", Withdraw(database))
	r.Get("/api/user/withdrawals", GetAllWithdraws(database))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, nfErr := w.Write([]byte("route does not exist"))
		if nfErr != nil {
			log.Println(nfErr)
		}
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, naErr := w.Write([]byte("sorry, only GET and POST methods are supported."))
		if naErr != nil {
			log.Println(naErr)
		}
	})
	return r
}
