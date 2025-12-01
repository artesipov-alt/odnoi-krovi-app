package v1

import (
	"context"

	bloodpoolv1 "github.com/artesipov-alt/odnoi-krovi-app/microservices/blood-microservice/gen/api/bloodpool/v1"
)

type BloodPoolService struct{}

// Обрабатывает запрос добавления питомнца в пулл поиска крови
func (s *BloodPoolService) AddPet(
	ctx context.Context,
	req *bloodpoolv1.AddToSearchPoolReq,
) (*bloodpoolv1.AddToSearchPoolResp, error) {
	return &bloodpoolv1.AddToSearchPoolResp{
		Status: true,
	}, nil
}

// Обрабатывает запрос добавления питомнца в пулл поиска крови
func (s *BloodPoolService) GetPetsByCriterias(
	ctx context.Context,
	req *bloodpoolv1.GetPetsByCriteriaReq,
) (*bloodpoolv1.GetPetsByCriteriaResp, error) {
	return &bloodpoolv1.GetPetsByCriteriaResp{
		Pets: []*bloodpoolv1.PetSearchResult{},
	}, nil
}
