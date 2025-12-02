package services

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
	bloodsearchv1 "github.com/artesipov-alt/odnoi-krovi-app/microservices/blood-microservice/gen/api/bloodsearch/v1"
	"github.com/artesipov-alt/odnoi-krovi-app/microservices/blood-microservice/gen/api/bloodsearch/v1/bloodsearchv1connect"
)

type BloodSearchService struct {
	client bloodsearchv1connect.BloodSearchPoolClient
}

// Создание нового сервиса с клиентом, указывая URL микросервиса, например "http://localhost:8081"
func NewBloodSearchService(baseURL string) *BloodSearchService {
	client := bloodsearchv1connect.NewBloodSearchPoolClient(
		http.DefaultClient,
		baseURL,
		connect.WithGRPC(),
	)
	return &BloodSearchService{client: client}
}

// Метод вызова gRPC RPC AddPet
func (s *BloodSearchService) AddPet(ctx context.Context, pet *bloodsearchv1.PetRow) (*bloodsearchv1.PetRowStatus, error) {
	resp, err := s.client.AddPet(ctx, pet)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Метод вызова gRPC RPC GetPets
func (s *BloodSearchService) GetPets(ctx context.Context, filter *bloodsearchv1.GetPetRows) (*bloodsearchv1.PetRows, error) {
	resp, err := s.client.GetPets(ctx, filter)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
