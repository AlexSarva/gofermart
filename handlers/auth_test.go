package handlers

import (
	"AlexSarva/gofermart/internal/app"
	"AlexSarva/gofermart/models"
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
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
	myChans := models.MyChans{InsertOrdersCh: make(chan models.Order)}
	Handler := *MyHandler(database, &myChans)
	ts := httptest.NewServer(&Handler)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := []byte(tt.requestBody)
			reqURL := tt.requestPath + tt.request
			request := httptest.NewRequest(tt.requestMethod, reqURL, bytes.NewBuffer(reqBody))
			request.Header.Set("Content-Type", tt.requestContentType)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			Handler.ServeHTTP(w, request)
			resp := w.Result()
			defer resp.Body.Close()
			// Проверяем StatusCode
			respStatusCode := resp.StatusCode
			wantStatusCode := tt.want.code
			assert.Equal(t, respStatusCode, wantStatusCode, fmt.Errorf("expected StatusCode %d, got %d", wantStatusCode, respStatusCode))
		})
	}
}

func TestUserAuthentication(t *testing.T) {
	database, dbErr := app.NewStorage("user=sarva password=77oFnWFF dbname=shortener sslmode=disable")

	userID := uuid.New()
	cookie, cookieExp := GenerateCookie(userID)
	hashedPassword, bcrypteErr := bcrypt.GenerateFromPassword([]byte("123"), 4)
	if bcrypteErr != nil {
		log.Println(bcrypteErr)
	}
	user := models.User{
		ID:        userID,
		Username:  "test",
		Password:  string(hashedPassword),
		Cookie:    cookie.String(),
		CookieExp: cookieExp,
	}

	database.Repo.NewUser(&user)

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
			requestPath:        "/api/user/login",
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name:               fmt.Sprintf("%s negative test #1", http.MethodPost),
			requestMethod:      http.MethodPost,
			requestContentType: "application/json",
			requestBody:        `{"login": "test", "passord": "123"}`,
			requestPath:        "/api/user/login",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:               fmt.Sprintf("%s negative test #2", http.MethodPost),
			requestMethod:      http.MethodPost,
			requestContentType: "application/json",
			requestBody:        `{"login": "test", "password": "123", "asdqs" : 123}`,
			requestPath:        "/api/user/login",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:               fmt.Sprintf("%s negative test #3", http.MethodPost),
			requestMethod:      http.MethodPost,
			requestContentType: "application/json",
			requestBody:        `{"login": "test", "password": 123}`,
			requestPath:        "/api/user/login",
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:               fmt.Sprintf("%s negative test #4", http.MethodPost),
			requestMethod:      http.MethodPost,
			requestContentType: "application/json",
			requestBody:        `{"login": "test100500", "password": "123"}`,
			requestPath:        "/api/user/login",
			want: want{
				code: http.StatusUnauthorized,
			},
		},
		{
			name:               fmt.Sprintf("%s negative test #5", http.MethodPost),
			requestMethod:      http.MethodPost,
			requestContentType: "application/json",
			requestBody:        `{"login": "test", "password": "1234"}`,
			requestPath:        "/api/user/login",
			want: want{
				code: http.StatusUnauthorized,
			},
		},
	}
	myChans := models.MyChans{InsertOrdersCh: make(chan models.Order)}
	Handler := *MyHandler(database, &myChans)
	ts := httptest.NewServer(&Handler)
	defer ts.Close()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody := []byte(tt.requestBody)
			reqURL := tt.requestPath + tt.request
			request := httptest.NewRequest(tt.requestMethod, reqURL, bytes.NewBuffer(reqBody))
			request.Header.Set("Content-Type", tt.requestContentType)
			// создаём новый Recorder
			w := httptest.NewRecorder()
			Handler.ServeHTTP(w, request)
			resp := w.Result()
			defer resp.Body.Close()
			// Проверяем StatusCode
			respStatusCode := resp.StatusCode
			wantStatusCode := tt.want.code
			assert.Equal(t, respStatusCode, wantStatusCode, fmt.Errorf("expected StatusCode %d, got %d", wantStatusCode, respStatusCode))
		})
	}
}
