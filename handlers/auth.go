package handlers

import (
	"AlexSarva/gofermart/internal/app"
	"AlexSarva/gofermart/models"
	"AlexSarva/gofermart/storage/storagepg"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strings"
)

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
