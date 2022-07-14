package handlers

import (
	"AlexSarva/gofermart/internal/app"
	"encoding/json"
	"net/http"
)

func GetBalance(database *app.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		headerContentType := r.Header.Get("Content-Length")
		if len(headerContentType) != 0 {
			messageResponse(w, "Content-Length is not equal 0", "application/json", http.StatusBadRequest)
			return
		}

		userID, cookieErr := GetCookie(r)
		if cookieErr != nil {
			messageResponse(w, "User unauthorized: "+cookieErr.Error(), "application/json", http.StatusUnauthorized)
			return
		}

		balance, balanceErr := database.Repo.GetBalance(userID)
		if balanceErr != nil {
			messageResponse(w, "Internal Server Error: "+balanceErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		balanceList, balanceListErr := json.Marshal(balance)
		if balanceListErr != nil {
			panic(balanceListErr)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(balanceList)
	}
}
