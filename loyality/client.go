package loyality

import (
	"AlexSarva/gofermart/logger"
	"AlexSarva/gofermart/models"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"gopkg.in/eapache/go-resiliency.v1/retrier"
	"gopkg.in/h2non/gentleman-retry.v2"
	"gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/plugins/timeout"
)

// ErrInternalServer error that occurs when loyalty service doesn't work
var ErrInternalServer = errors.New("ErrInternalServer")

// ErrEmptyOrder error that occurs when get empty order
var ErrEmptyOrder = errors.New("empty order")

// ProcessingClient client for check processing in loyalty system
type ProcessingClient struct {
	Client *gentleman.Client
}

// NewProcessingClient generate new client for check processing in loyalty system
func NewProcessingClient(serviceAddress, basicURL string) *ProcessingClient {
	log.Println("LoyalityServer: ", serviceAddress+basicURL)
	cli := gentleman.New()
	cli.Use(logger.New(os.Stdout))
	cli.Use(timeout.Request(60 * time.Second))
	cli.Use(retry.New(retrier.New(retrier.ExponentialBackoff(5, 100*time.Millisecond), nil)))
	cli.URL(serviceAddress + basicURL)
	return &ProcessingClient{
		Client: cli,
	}
}

// GetOrder method that check order number in loyalty system and returns result of processing
func (pc *ProcessingClient) GetOrder(orderNum string) (models.ProcessingOrder, error) {
	req := pc.Client.Request()
	req.Method("GET")
	req.AddPath(fmt.Sprintf("/%s", orderNum))
	res, err := req.Send()
	var order models.ProcessingOrder
	if err != nil {
		return order, err
	}

	switch res.StatusCode {
	case http.StatusInternalServerError:
		log.Printf("Internas server error: %d\n", res.StatusCode)
		return order, ErrInternalServer
	case http.StatusTooManyRequests:
		log.Printf("Too Many Requests: %d\n", res.StatusCode)
		time.Sleep(time.Second * 60)
	case http.StatusOK:
		if UnmarshErr := json.Unmarshal(res.Bytes(), &order); UnmarshErr != nil {
			return order, UnmarshErr
		}
	}

	emptyOrder := models.ProcessingOrder{
		OrderNum: "",
		Status:   "",
		Accrual:  nil,
	}

	if order == emptyOrder {
		return order, ErrEmptyOrder
	}

	if order.OrderNum == "" {
		return order, ErrEmptyOrder
	}

	return order, nil
}
