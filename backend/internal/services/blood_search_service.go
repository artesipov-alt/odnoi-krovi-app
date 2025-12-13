package services

import (
	"context"
	"net/http"
	"time"

	bloodrequestv1 "github.com/artesipov-alt/odnoi-krovi-app/microservices/blood-microservice/gen/api/bloodrequest/v1"
	"github.com/artesipov-alt/odnoi-krovi-app/microservices/blood-microservice/gen/api/bloodrequest/v1/bloodrequestv1connect"
)

// BloodRequestClient - клиент для взаимодействия с blood-microservice
type BloodRequestClient struct {
	client bloodrequestv1connect.BloodRequestPoolClient
}

// NewBloodRequestClient создает новый клиент с настроенными таймаутами, указывая URL микросервиса, например "http://localhost:8081"
func NewBloodRequestClient(baseURL string) *BloodRequestClient {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	client := bloodrequestv1connect.NewBloodRequestPoolClient(
		httpClient,
		baseURL,
	)
	return &BloodRequestClient{client: client}
}

// AddPet добавляет питомца в пул поиска крови
func (c *BloodRequestClient) AddPet(ctx context.Context, pet *bloodrequestv1.BloodRequest) (*bloodrequestv1.BloodRequestStatus, error) {
	resp, err := c.client.AddPet(ctx, pet)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetPets получает питомцев из пула поиска крови по фильтрам
func (c *BloodRequestClient) GetPets(ctx context.Context, filter *bloodrequestv1.GetBloodRequests) (*bloodrequestv1.BloodRequests, error) {
	resp, err := c.client.GetPets(ctx, filter)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
