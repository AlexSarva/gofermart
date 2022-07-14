package handlers

import (
	"AlexSarva/gofermart/internal/app"
	"AlexSarva/gofermart/models"
	"AlexSarva/gofermart/utils/luhn"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func PostOrder(database *app.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		headerContentType := r.Header.Get("Content-Type")
		if !strings.Contains("text/plain", headerContentType) {
			messageResponse(w, "Content Type is not application/json or application/x-gzip", "application/json", http.StatusBadRequest)
			return
		}

		userID, cookieErr := GetCookie(r)
		if cookieErr != nil {
			messageResponse(w, "User unauthorized: "+cookieErr.Error(), "application/json", http.StatusUnauthorized)
			return
		}

		b, err := readBodyBytes(r)
		if err != nil {
			messageResponse(w, "Problem in body", "application/json", http.StatusBadRequest)
			return
		}

		body, bodyErr := io.ReadAll(b)
		if bodyErr != nil {
			messageResponse(w, "Problem in body", "application/json", http.StatusBadRequest)
			return
		}

		orderNum, convErr := strconv.Atoi(string(body))
		if convErr != nil {
			messageResponse(w, "invalid order number format", "application/json", http.StatusUnprocessableEntity)
			return
		}

		orderDB, orderDBErr := database.Repo.CheckOrder(orderNum)
		if orderDBErr != nil {
			messageResponse(w, "Internal Server Error: "+orderDBErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		if orderDB.OrderNum == orderNum && orderDB.UserID != userID {
			messageResponse(w, "the order number has already been uploaded by another user", "application/json", http.StatusConflict)
			return
		}

		if orderDB.OrderNum == orderNum && orderDB.UserID == userID {
			messageResponse(w, "order number has already been uploaded by this user", "application/json", http.StatusOK)
			return
		}

		if !luhn.Valid(orderNum) {
			messageResponse(w, "invalid order number format", "application/json", http.StatusUnprocessableEntity)
			return
		}

		var order models.Order
		order.UserID, order.OrderNum = userID, orderNum

		insertErr := database.Repo.NewOrder(&order)
		if insertErr != nil {
			messageResponse(w, "Internal Server Error: "+orderDBErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		messageResponse(w, "new order number accepted for processing", "application/json", http.StatusAccepted)
		return
	}
}
