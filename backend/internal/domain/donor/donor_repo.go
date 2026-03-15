package donor

import (
	"context"

	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
)

// Repository определяет интерфейс для работы с откликами доноров
type Repository interface {
	CreateDonorResponse(ctx context.Context, resp *donormodel.DonorResponse) (*donormodel.DonorResponse, error)
	GetDonorResponseByID(ctx context.Context, id string) (*donormodel.DonorResponse, error)
	UpdateDonorResponseStatus(ctx context.Context, id, status string) error
	GetRecipient(ctx context.Context, id string) (*donormodel.Recipient, error)
	DeleteDonorResponse(ctx context.Context, id string) error
	GetDonorResponsesByRequestID(ctx context.Context, reqID string) ([]*donormodel.DonorResponse, error)
	GetDonorResponsesByDonorID(ctx context.Context, donorID string) ([]*donormodel.DonorResponse, error)
	ExistsByID(ctx context.Context, id string) (bool, error)
	ExistsByRequestID(ctx context.Context, reqID string) (bool, error)
	ExistsByDonorID(ctx context.Context, donorID string) (bool, error)
	Count(ctx context.Context) (int, error)
}
