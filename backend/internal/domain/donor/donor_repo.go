package donor

import (
	"context"
	"time"

	bloodreqmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
)

type Repository interface {
	// Write Methods
	CreateDonorResponse(ctx context.Context, resp *donormodel.DonorResponse) (*donormodel.DonorResponse, error)
	UpdateDonorResponseStatus(ctx context.Context, id string, status donormodel.DonorResponseStatus) error
	DeleteDonorResponse(ctx context.Context, id string) error
	Update(ctx context.Context, resp *donormodel.DonorResponse) error

	// Read Methods
	GetDonorResponseByID(ctx context.Context, id string) (*donormodel.DonorResponse, error)
	GetRecipient(ctx context.Context, id string) (*bloodreqmodel.BloodRequestWithMatchingDonors, error)
	GetDonorResponsesByRequestID(ctx context.Context, reqID string) ([]*donormodel.DonorResponse, error)
	GetDonorResponsesByDonorID(ctx context.Context, donorID string) ([]*donormodel.DonorResponse, error)
	ExistsByID(ctx context.Context, id string) (bool, error)
	ExistsByRequestID(ctx context.Context, reqID string) (bool, error)
	ExistsByDonorID(ctx context.Context, donorID string) (bool, error)
	Count(ctx context.Context) (int, error)
	GetByPetID(ctx context.Context, petID string) (*donormodel.DonorResponse, error)
	GetByPetIDs(ctx context.Context, petIDs []string, ignoreSoftDelete bool) (map[string][]*donormodel.DonorResponse, error)
	FindNotConfirmed(ctx context.Context, cutoffTime time.Time) ([]*donormodel.DonorResponse, error)
}
