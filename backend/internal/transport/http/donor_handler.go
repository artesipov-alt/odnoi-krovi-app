package http

import (
	"context"
	"net/http"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/transport/http/dto"

	"github.com/danielgtaylor/huma/v2"
)

type DonorHandler struct {
	recipientsListHandler *any
}

// NewBloodRequestHandler создает новый обработчик для заявок на поиск крови
func NewDonorHandler() *DonorHandler {
	return &DonorHandler{}
}

// Register регистрирует маршруты заявок на поиск крови в Huma API
func (h *DonorHandler) Register(api huma.API) {

	// Получить список заявок на поиск крови
	huma.Register(api, huma.Operation{
		OperationID: "get-recipients",
		Method:      http.MethodGet,
		Path:        "/v1/donor/{id}/recipients-list",
		Summary:     "Получить список заявок на поиск крови",
		Description: "Возвращает список реципиентов по фильтрам",
		Tags:        []string{"donor-v1"},
	}, h.GetRecipientsList)

}

func (h *DonorHandler) GetRecipientsList(ctx context.Context, input *dto.GetRecipientsListInput) (*dto.ListRecipientsOutput, error) {
	return &dto.ListRecipientsOutput{Body: dto.RecipientsList{}}, nil
}
