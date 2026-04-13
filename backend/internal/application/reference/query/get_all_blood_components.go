package query

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/common"
)

type GetAllBloodComponentsHandler struct{}

func NewGetAllBloodComponentsHandler() *GetAllBloodComponentsHandler {
	return &GetAllBloodComponentsHandler{}
}

func (h *GetAllBloodComponentsHandler) Handle(ctx context.Context) ([]common.BloodComponent, error) {
	return common.BloodComponents, nil
}
