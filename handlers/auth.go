package handlers

import (
	"AlexSarva/gofermart/internal/app"
	"AlexSarva/gofermart/models"
	"AlexSarva/gofermart/storage/storagepg"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
)

// UserRegistration регистрация нового пользователя
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
		userCookie := GenerateCookie(userID)
		hashedPassword, bcrypteErr := bcrypt.GenerateFromPassword([]byte(user.Password), 4)
		if bcrypteErr != nil {
			log.Println(bcrypteErr)
		}

		user.ID, user.Password, user.Cookie = userID, string(hashedPassword), userCookie.String()

		newUserErr := database.Repo.NewUser(&user)
		if newUserErr != nil {
			if newUserErr == storagepg.ErrDuplicatePK {
				messageResponse(w, "login is already busy", "application/json", http.StatusConflict)
				return
			}
			messageResponse(w, "Internal Server Error "+newUserErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &userCookie)
		messageResponse(w, "user successfully registered and authenticated", "application/json", http.StatusOK)

		log.Printf("user registered \n %+v\n", user)

	}
}

// UserAuthentication - аутентификация пользователя
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
			if userDBErr == sql.ErrNoRows {
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

		w.Header().Add("Set-Cookie", userDB.Cookie)
		messageResponse(w, "user successfully authenticated", "application/json", http.StatusOK)
	}
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
