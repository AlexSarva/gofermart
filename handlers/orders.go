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

// PostOrder - Load order number
//
// Handler POST /api/user/orders
//
// The handler is available only to authenticated users.
// The order number is a sequence of digits of arbitrary length.
// The order number is checked for correct input using the Luhn algorithm.
// Request format:
// 12345678903
//
// Possible response codes:
// 200 - the order number has already been uploaded by this user;
// 202 - new order number accepted for processing;
// 400 - invalid request format;
// 401 - user not authenticated;
// 409 - the order number has already been uploaded by another user;
// 422 - invalid order number format;
// 500 - an internal server error.
func PostOrder(database *app.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		headerContentType := r.Header.Get("Content-Type")
		if !strings.Contains("text/plain", headerContentType) {
			messageResponse(w, "Content Type is not application/json or application/x-gzip", "application/json", http.StatusBadRequest)
			return
		}

		// Проверка авторизации по токену
		userID, tokenErr := GetToken(r)
		if tokenErr != nil {
			messageResponse(w, "User unauthorized: "+tokenErr.Error(), "application/json", http.StatusUnauthorized)
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

// GetOrders Getting a list of loaded order numbers
//
// Handler: GET /api/user/orders.
//
// The handler is available only to an authorized user. Order numbers in the search results should be sorted by download time from oldest to newest. The date format is RFC3339.
// Available settlement processing statuses:
// NEW — the order was loaded into the system, but was not processed;
// PROCESSING - the reward for the order is calculated;
// INVALID — the system for calculating remuneration refused to calculate;
// PROCESSED — order data has been checked and billing information has been successfully received.
//
// Possible response codes:
// 200 - successful processing of the request.
// Response format:
//   [
//       {
//           "number": "9278923470",
//           "status": "PROCESSED",
//           accrual: 500
//           "uploaded_at": "2020-12-10T15:15:45+03:00"
//       },
//       {
//           "number": "12345678903",
//           "status": "PROCESSING",
//           "uploaded_at": "2020-12-10T15:12:01+03:00"
//       },
//       {
//           "number": "346436439",
//           "status": "INVALID",
//           "uploaded_at": "2020-12-09T16:09:53+03:00"
//       }
//   ]
//
// 204 - No response data.
// 401 - The user is not authorized.
// 500 - an internal server error.
func GetOrders(database *app.Database) http.HandlerFunc {
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

		orders, ordersErr := database.Repo.GetOrders(userID)
		if ordersErr != nil {
			if ordersErr == storagepg.ErrNoValues {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNoContent)
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
