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

func TestPostOrder(t *testing.T) {
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

	subUserID := uuid.New()
	subCookie, subCookieExp := GenerateCookie(subUserID)
	subHashedPassword, subBcrypteErr := bcrypt.GenerateFromPassword([]byte("123"), 4)
	if subBcrypteErr != nil {
		log.Println(subBcrypteErr)
	}
	subUser := models.User{
		ID:        subUserID,
		Username:  "test2",
		Password:  string(subHashedPassword),
		Cookie:    subCookie.String(),
		CookieExp: subCookieExp,
	}

	database.Repo.NewUser(&subUser)

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
		requestCookie          string
		want                   want
	}{
		{
			name:               fmt.Sprintf("%s positive #1", http.MethodPost),
			requestMethod:      http.MethodPost,
			requestContentType: "text/plain",
			requestBody:        "12345678903", // новый номер заказа принят в обработку
			requestPath:        "/api/user/orders",
			requestCookie:      cookie.String(),
			want: want{
				code: http.StatusAccepted,
			},
		},
		{
			name:               fmt.Sprintf("%s negative #1", http.MethodPost),
			requestMethod:      http.MethodPost,
			requestContentType: "text/plain",
			requestBody:        "12345678903", // номер заказа уже был загружен этим пользователем
			requestPath:        "/api/user/orders",
			requestCookie:      cookie.String(),
			want: want{
				code: http.StatusOK,
			},
		},
		{
			name:               fmt.Sprintf("%s negative #2", http.MethodPost),
			requestMethod:      http.MethodPost,
			requestContentType: "text/plain",
			requestBody:        "12345-678903", // неверный формат номера заказа
			requestPath:        "/api/user/orders",
			requestCookie:      cookie.String(),
			want: want{
				code: http.StatusUnprocessableEntity,
			},
		},
		{
			name:               fmt.Sprintf("%s negative #3", http.MethodPost),
			requestMethod:      http.MethodPost,
			requestContentType: "application/json", // неверный Content-Type неверный формат запроса
			requestBody:        "12345-678903",
			requestPath:        "/api/user/orders",
			requestCookie:      cookie.String(),
			want: want{
				code: http.StatusBadRequest,
			},
		},
		{
			name:               fmt.Sprintf("%s negative #4", http.MethodPost),
			requestMethod:      http.MethodPost,
			requestContentType: "text/plain",
			requestBody:        "1234567890", // не проходит проверку по Луну
			requestPath:        "/api/user/orders",
			requestCookie:      cookie.String(),
			want: want{
				code: http.StatusUnprocessableEntity,
			},
		},
		{
			name:               fmt.Sprintf("%s negative #5", http.MethodPost),
			requestMethod:      http.MethodPost,
			requestContentType: "text/plain",
			requestBody:        "1234567890", // без cookie
			requestPath:        "/api/user/orders",
			requestCookie:      "",
			want: want{
				code: http.StatusUnauthorized,
			},
		},
		{
			name:               fmt.Sprintf("%s negative #6", http.MethodPost),
			requestMethod:      http.MethodPost,
			requestContentType: "text/plain",
			requestBody:        "12345678903", // номер заказа уже был загружен другим пользователем
			requestPath:        "/api/user/orders",
			requestCookie:      subCookie.String(),
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
			request.Header.Set("Cookie", tt.requestCookie)
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
