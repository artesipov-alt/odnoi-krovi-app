package services

import (
	"context"
	"net/http"
	"time"

	bloodsearchv1 "github.com/artesipov-alt/odnoi-krovi-app/microservices/blood-microservice/gen/api/bloodsearch/v1"
	"github.com/artesipov-alt/odnoi-krovi-app/microservices/blood-microservice/gen/api/bloodsearch/v1/bloodsearchv1connect"
)

// BloodSearchClient - клиент для взаимодействия с blood-microservice
type BloodSearchClient struct {
	client bloodsearchv1connect.BloodSearchPoolClient
}

// NewBloodSearchClient создает новый клиент с настроенными таймаутами, указывая URL микросервиса, например "http://localhost:8081"
func NewBloodSearchClient(baseURL string) *BloodSearchClient {
	httpClient := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	client := bloodsearchv1connect.NewBloodSearchPoolClient(
		httpClient,
		baseURL,
	)
	return &BloodSearchClient{client: client}
}

// AddPet добавляет питомца в пул поиска крови
func (c *BloodSearchClient) AddPet(ctx context.Context, pet *bloodsearchv1.PetRow) (*bloodsearchv1.PetRowStatus, error) {
	resp, err := c.client.AddPet(ctx, pet)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetPets получает питомцев из пула поиска крови по фильтрам
func (c *BloodSearchClient) GetPets(ctx context.Context, filter *bloodsearchv1.GetPetRows) (*bloodsearchv1.PetRows, error) {
	resp, err := c.client.GetPets(ctx, filter)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
