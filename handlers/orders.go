package handlers

import (
	"AlexSarva/gofermart/internal/app"
	"AlexSarva/gofermart/models"
	"AlexSarva/gofermart/storage/storagepg"
	"AlexSarva/gofermart/utils/luhn"
	"encoding/json"
	"fmt"
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

		if !luhn.Valid(orderNum) {
			messageResponse(w, "invalid order number format", "application/json", http.StatusUnprocessableEntity)
			return
		}

		orderNumStr := fmt.Sprintf("%d", orderNum)

		orderDB, orderDBErr := database.Repo.CheckOrder(orderNumStr)
		if orderDBErr != nil {
			messageResponse(w, "Internal Server Error: "+orderDBErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		if orderDB.OrderNum == orderNumStr && orderDB.UserID != userID {
			messageResponse(w, "the order number has already been uploaded by another user", "application/json", http.StatusConflict)
			return
		}

		if orderDB.OrderNum == orderNumStr && orderDB.UserID == userID {
			messageResponse(w, "order number has already been uploaded by this user", "application/json", http.StatusOK)
			return
		}

		var order models.Order
		order.UserID, order.OrderNum = userID, orderNumStr

		insertErr := database.Repo.NewOrder(&order)
		if insertErr != nil {
			messageResponse(w, "Internal Server Error: "+insertErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		messageResponse(w, "new order number accepted for processing", "application/json", http.StatusAccepted)
	}
}

func GetOrders(database *app.Database) http.HandlerFunc {
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

		orders, ordersErr := database.Repo.GetOrders(userID)
		if ordersErr != nil {
			if ordersErr == storagepg.ErrNoValues {
				messageResponse(w, "no data to answer: "+storagepg.ErrNoValues.Error(), "application/json", http.StatusNoContent)
				return
			}
			messageResponse(w, "Internal Server Error: "+ordersErr.Error(), "application/json", http.StatusInternalServerError)
			return
		}

		ordersList, ordersListErr := json.Marshal(orders)
		if ordersListErr != nil {
			panic(ordersListErr)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(ordersList)
	}
}
