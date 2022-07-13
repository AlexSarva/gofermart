package handlers

import (
	"AlexSarva/gofermart/internal/app"
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserRegistration(t *testing.T) {
	database, dbErr := app.NewStorage("user=sarva password=77oFnWFF dbname=shortener sslmode=disable")
	if dbErr != nil {
		log.Fatal(dbErr)
	}
	type want struct {
		code            int
		location        string
		contentType     string
		contentEncoding string
		responseFormat  bool
		response        string
	}

	tests := []struct {
		name                   string
		request                string
		requestPath            string
		requestMethod          string
		requestBody            string
		requestCompressBody    []byte
		requestContentType     string
		requestAcceptEncoding  string
		requestContentEncoding string
		want                   want
	}{
		{
			name:               fmt.Sprintf("%s positive test #1", http.MethodPost),
			requestMethod:      http.MethodPost,
			requestContentType: "application/json",
			requestBody:        `{"login": "test", "password": "123"}`,
			requestPath:        "/api/user/register",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name:               fmt.Sprintf("%s negative test #1", http.MethodPost),
			requestMethod:      http.MethodPost,
			requestContentType: "application/json",
			requestBody:        `{"login": "test", "passord": "123"}`,
			requestPath:        "/api/user/register",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:               fmt.Sprintf("%s negative test #2", http.MethodPost),
			requestMethod:      http.MethodPost,
			requestContentType: "application/json",
			requestBody:        `{"login": "test", "password": "123", "asdqs" : 123}`,
			requestPath:        "/api/user/register",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:               fmt.Sprintf("%s negative test #3", http.MethodPost),
			requestMethod:      http.MethodPost,
			requestContentType: "application/json",
			requestBody:        `{"login": "test", "password": 123}`,
			requestPath:        "/api/user/register",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:               fmt.Sprintf("%s negative test #4", http.MethodPost),
			requestMethod:      http.MethodPost,
			requestContentType: "application/json",
			requestBody:        `{"login": "test", "password": "123"}`,
			requestPath:        "/api/user/register",
			want: want{
				code: http.StatusConflict,
			},
		},
	}
	Handler := *MyHandler(database)
	ts := httptest.NewServer(&Handler)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := []byte(tt.requestBody)
			reqURL := tt.requestPath + tt.request
			var request *http.Request
			request = httptest.NewRequest(tt.requestMethod, reqURL, bytes.NewBuffer(reqBody))
			request.Header.Set("Content-Type", tt.requestContentType)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			Handler.ServeHTTP(w, request)
			resp := w.Result()
			// Проверяем StatusCode
			respStatusCode := resp.StatusCode
			wantStatusCode := tt.want.code
			assert.Equal(t, respStatusCode, wantStatusCode, fmt.Errorf("expected StatusCode %d, got %d", wantStatusCode, respStatusCode))
		})
	}
}
