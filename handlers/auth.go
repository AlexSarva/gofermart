package handlers

import (
	"AlexSarva/gofermart/internal/app"
	"AlexSarva/gofermart/models"
	"AlexSarva/gofermart/storage/storagepg"
	"AlexSarva/gofermart/utils"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const timeLayout = "2006-01-02 15:04:05"

// UserRegistration - user registration method
//
// Handler POST /api/user/register
//
// Registration is performed by a pair of login/password.
// Each login must be set.
// After successful registration, automatic user authentication is required.
// post message should contain such body:
//    "login": "<login>",
//    "password": "<password>"
//
// Possible response codes:
// 200 - user successfully registered and authenticated;
// 400 - invalid request format;
// 409 - login is already taken;
// 500 - an internal server error.
func UserRegistration(database *app.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		headerContentType := r.Header.Get("Content-Type")
		if !strings.Contains("application/json, application/x-gzip", headerContentType) {
			messageResponse(w, "Content Type is not application/json or application/x-gzip", "application/json", http.StatusBadRequest)
			return
		}

		var user models.User
		var unmarshalErr *json.UnmarshalTypeError

		b, err := readBodyBytes(r)
		if err != nil {
			messageResponse(w, "Problem in body", "application/json", http.StatusBadRequest)
			return
		}

		decoder := json.NewDecoder(b)
		decoder.DisallowUnknownFields()
		errDecode := decoder.Decode(&user)

		if errDecode != nil {
			if errors.As(errDecode, &unmarshalErr) {
				messageResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, "application/json", http.StatusBadRequest)
			} else {
				messageResponse(w, "Bad Request. "+errDecode.Error(), "application/json", http.StatusBadRequest)
			}
			return
		}

		userID := uuid.New()
		userCookie, userCookieExp := GenerateCookie(userID)
		hashedPassword, bcrypteErr := bcrypt.GenerateFromPassword([]byte(user.Password), 4)
		if bcrypteErr != nil {
			log.Println(bcrypteErr)
		}

		user.ID, user.Password, user.Cookie, user.CookieExp = userID, string(hashedPassword), userCookie.String(), userCookieExp

		newUserErr := database.Repo.NewUser(&user)
		if newUserErr != nil {
			if newUserErr == storagepg.ErrDuplicatePK {
				messageResponse(w, "login is already busy", "application/json", http.StatusConflict)
				return
			}
			messageResponse(w, "Internal Server Error "+newUserErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		generatedAt := time.Now().Format(timeLayout)
		expiresAt := userCookie.Expires.Format(timeLayout)
		tokenDetails := models.Token{
			TokenType:   "Bearer",
			AuthToken:   userCookie.Value,
			GeneratedAt: generatedAt,
			ExpiresAt:   expiresAt,
		}
		jsonResp, _ := json.Marshal(tokenDetails)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Authorization", tokenDetails.TokenType+" "+tokenDetails.AuthToken)
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResp)
	}
}

// UserAuthentication - user authentication method
//
// Handler POST /api/user/login
//
// Authentication is performed by a login/password pair.
// Request format:
//    {"login": "<login>",
//    "password": "<password>"}
//
// Possible response codes:
// 200 - user successfully authenticated;
// 400 - invalid request format;
// 401 - invalid login/password pair;
// 500 - an internal server error.
func UserAuthentication(database *app.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		headerContentType := r.Header.Get("Content-Type")
		if !strings.Contains("application/json, application/x-gzip", headerContentType) {
			messageResponse(w, "Content Type is not application/json or application/x-gzip", "application/json", http.StatusBadRequest)
			return
		}

		var user models.User
		var unmarshalErr *json.UnmarshalTypeError

		b, err := readBodyBytes(r)
		if err != nil {
			messageResponse(w, "Problem in body", "application/json", http.StatusBadRequest)
			return
		}

		decoder := json.NewDecoder(b)
		decoder.DisallowUnknownFields()
		errDecode := decoder.Decode(&user)

		if errDecode != nil {
			if errors.As(errDecode, &unmarshalErr) {
				messageResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, "application/json", http.StatusBadRequest)
			} else {
				messageResponse(w, "Bad Request. "+errDecode.Error(), "application/json", http.StatusBadRequest)
			}
			return
		}

		userDB, userDBErr := database.Repo.GetUser(user.Username)
		if userDBErr != nil {
			if errors.Is(userDBErr, sql.ErrNoRows) {
				messageResponse(w, "User unauthorized: "+user.Username+", please register at /api/user/register", "application/json", http.StatusUnauthorized)
				return
			}
			messageResponse(w, "Internal Server Error: "+userDBErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		cryptErr := bcrypt.CompareHashAndPassword([]byte(userDB.Password), []byte(user.Password))
		if cryptErr != nil {
			messageResponse(w, "User unauthorized: "+cryptErr.Error(), "application/json", http.StatusUnauthorized)
			return
		}
		// TODO Предусмотреть обновление куки
		if userDB.CookieExp.Before(time.Now()) {
			log.Println("cookie expired")
		}

		// Авторизация по токену
		generatedAt := time.Now().Format(timeLayout)
		expiresAt := userDB.CookieExp.Format(timeLayout)
		cookieSession, cookieSessionErr := utils.ParseCookie(userDB.Cookie)
		if cookieSessionErr != nil {
			log.Println(cookieSessionErr)
		}

		tokenDetails := models.Token{
			TokenType:   "Bearer",
			AuthToken:   cookieSession,
			GeneratedAt: generatedAt,
			ExpiresAt:   expiresAt,
		}
		jsonResp, _ := json.Marshal(tokenDetails)
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Authorization", tokenDetails.TokenType+" "+tokenDetails.AuthToken)
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResp)
	}
}
