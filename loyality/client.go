package loyality

import (
	"AlexSarva/gofermart/models"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/eapache/go-resiliency.v1/retrier"
	"gopkg.in/h2non/gentleman-retry.v2"
	"gopkg.in/h2non/gentleman.v2"
	"gopkg.in/h2non/gentleman.v2/plugins/timeout"
	"gopkg.in/h2non/gentleman.v2/plugins/url"
	"log"
	"net/http"
	"time"
)

var ErrInternalServer = errors.New("ErrInternalServer")

type ProcessingClient struct {
	Client *gentleman.Client
}

func NewProcessingClient(serviceAddress string) *ProcessingClient {
	cli := gentleman.New()
	cli.Use(timeout.Request(60 * time.Second))
	cli.Use(retry.New(retrier.New(retrier.ExponentialBackoff(5, 100*time.Millisecond), nil)))
	cli.Use(url.BaseURL(serviceAddress))
	return &ProcessingClient{
		Client: cli,
	}
}

func (pc *ProcessingClient) GetOrder(orderNum string) (models.ProcessingOrder, error) {
	req := pc.Client.Request()
	req.Path(fmt.Sprintf("/%s", orderNum))
	res, err := req.Send()
	var order models.ProcessingOrder
	if err != nil {
		log.Printf("Request error: %s\n", err)
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
	return order, nil
}
