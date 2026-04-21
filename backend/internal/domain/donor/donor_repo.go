package donor

import (
	"context"

	bloodreqmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
)

type Repository interface {
	CreateDonorResponse(ctx context.Context, resp *donormodel.DonorResponse) (*donormodel.DonorResponse, error)
	GetDonorResponseByID(ctx context.Context, id string) (*donormodel.DonorResponse, error)
	UpdateDonorResponseStatus(ctx context.Context, id string, status donormodel.DonorResponseStatus) error
	GetRecipient(ctx context.Context, id string) (*bloodreqmodel.BloodRequestWithMatchingDonors, error)
	DeleteDonorResponse(ctx context.Context, id string) error
	GetDonorResponsesByRequestID(ctx context.Context, reqID string) ([]*donormodel.DonorResponse, error)
	GetDonorResponsesByDonorID(ctx context.Context, donorID string) ([]*donormodel.DonorResponse, error)
	ExistsByID(ctx context.Context, id string) (bool, error)
	ExistsByRequestID(ctx context.Context, reqID string) (bool, error)
	ExistsByDonorID(ctx context.Context, donorID string) (bool, error)
	Count(ctx context.Context) (int, error)
	GetByPetID(ctx context.Context, petID string) (*donormodel.DonorResponse, error)
	GetByPetIDs(ctx context.Context, petIDs []string, ignoreSoftDelete bool) (map[string][]*donormodel.DonorResponse, error)
	//=============================================
	Accept(ctx context.Context, id string) error
	Reject(ctx context.Context, res *donormodel.DonorResponse) error
	Complete(ctx context.Context, id string, amount float64) error
	Confirm(ctx context.Context, id string, amount float64) error
	Cancel(ctx context.Context, id string) error
}
