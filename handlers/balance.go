package handlers

import (
	"AlexSarva/gofermart/internal/app"
	"AlexSarva/gofermart/models"
	"AlexSarva/gofermart/storage/storagepg"
	"AlexSarva/gofermart/utils/luhn"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

func GetBalance(database *app.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		headerContentType := r.Header.Get("Content-Length")
		if len(headerContentType) != 0 {
			messageResponse(w, "Content-Length is not equal 0", "application/json", http.StatusBadRequest)
			return
		}

		// Проверка авторизации по токену
		userID, tokenErr := GetToken(r)
		if tokenErr != nil {
			messageResponse(w, "User unauthorized: "+tokenErr.Error(), "application/json", http.StatusUnauthorized)
			return
		}

		//userID, cookieErr := GetCookie(r)
		//if cookieErr != nil {
		//	messageResponse(w, "User unauthorized: "+cookieErr.Error(), "application/json", http.StatusUnauthorized)
		//	return
		//}

		balance, balanceErr := database.Repo.GetBalance(userID)
		if balanceErr != nil {
			messageResponse(w, "Internal Server Error: "+balanceErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		// TODO обработать паники
		balanceJSON, balanceJSONErr := json.Marshal(balance)
		if balanceJSONErr != nil {
			panic(balanceJSONErr)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(balanceJSON)
	}
}

func Withdraw(database *app.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		headerContentType := r.Header.Get("Content-Type")
		if !strings.Contains("application/json, application/x-gzip", headerContentType) {
			messageResponse(w, "Content Type is not application/json or application/x-gzip", "application/json", http.StatusBadRequest)
			return
		}

		// Проверка авторизации по токену
		userID, tokenErr := GetToken(r)
		if tokenErr != nil {
			messageResponse(w, "User unauthorized: "+tokenErr.Error(), "application/json", http.StatusUnauthorized)
			return
		}

		//userID, cookieErr := GetCookie(r)
		//if cookieErr != nil {
		//	messageResponse(w, "User unauthorized: "+cookieErr.Error(), "application/json", http.StatusUnauthorized)
		//	return
		//}

		var withdraw models.Withdraw
		var unmarshalErr *json.UnmarshalTypeError

		b, err := readBodyBytes(r)
		if err != nil {
			messageResponse(w, "Problem in body", "application/json", http.StatusBadRequest)
			return
		}

		decoder := json.NewDecoder(b)
		decoder.DisallowUnknownFields()
		errDecode := decoder.Decode(&withdraw)

		if errDecode != nil {
			if errors.As(errDecode, &unmarshalErr) {
				messageResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, "application/json", http.StatusBadRequest)
			} else {
				messageResponse(w, "Bad Request. "+errDecode.Error(), "application/json", http.StatusBadRequest)
			}
			return
		}

		orderNum, convErr := strconv.Atoi(withdraw.OrderNum)
		if convErr != nil {
			messageResponse(w, "invalid order number format", "application/json", http.StatusUnprocessableEntity)
			return
		}

		if !luhn.Valid(orderNum) {
			messageResponse(w, "invalid order number format", "application/json", http.StatusUnprocessableEntity)
			return
		}

		balance, balanceErr := database.Repo.GetBalance(userID)
		if balanceErr != nil {
			messageResponse(w, "Internal Server Error: "+balanceErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		if (balance.Current - balance.Withdraw) < withdraw.Withdraw {
			messageResponse(w, "there are not enough funds on the account", "application/json", http.StatusPaymentRequired)
			return
		}

		withdraw.UserID = userID
		withdrawErr := database.Repo.NewWithdraw(&withdraw)
		if withdrawErr != nil {
			if withdrawErr == storagepg.ErrDuplicatePK {
				messageResponse(w, "order number has already been uploaded", "application/json", http.StatusConflict)
				return
			}
			messageResponse(w, "Internal Server Error "+withdrawErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		messageResponse(w, "successful request processing", "application/json", http.StatusOK)
	}
}

func GetAllWithdraws(database *app.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		headerContentType := r.Header.Get("Content-Length")
		if len(headerContentType) != 0 {
			messageResponse(w, "Content-Length is not equal 0", "application/json", http.StatusBadRequest)
			return
		}

		// Проверка авторизации по токену
		userID, tokenErr := GetToken(r)
		if tokenErr != nil {
			messageResponse(w, "User unauthorized: "+tokenErr.Error(), "application/json", http.StatusUnauthorized)
			return
		}

		//userID, cookieErr := GetCookie(r)
		//if cookieErr != nil {
		//	messageResponse(w, "User unauthorized: "+cookieErr.Error(), "application/json", http.StatusUnauthorized)
		//	return
		//}

		withdraws, withdrawsErr := database.Repo.GetAllWithdraw(userID)
		if withdrawsErr != nil {
			if withdrawsErr == storagepg.ErrNoValues {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNoContent)
				//messageResponse(w, "no data to answer: "+storagepg.ErrNoValues.Error(), "application/json", http.StatusNoContent)
				return
			}
			messageResponse(w, "Internal Server Error: "+withdrawsErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		withdrawsList, withdrawsListErr := json.Marshal(withdraws)
		if withdrawsListErr != nil {
			panic(withdrawsListErr)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(withdrawsList)
	}
}
