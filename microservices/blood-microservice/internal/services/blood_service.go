package v1

import (
	"context"

	bloodv1 "github.com/artesipov-alt/odnoi-krovi-app/poolservice/gen/api/blood/v1"
)

type BloodServer struct{}

// Greet обрабатывает запрос приветствия
func (s *BloodServer) AddToSearchPool(
	ctx context.Context,
	req *bloodv1.AddToSearchPoolReq,
) (*bloodv1.AddToSearchPoolResp, error) {
	return &bloodv1.AddToSearchPoolResp{
		Status: true,
	}, nil
}
