package pg

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/bloodsearchrequest"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/donorresponse"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/pet"
)

// EntDonorResponseRepository implements DonorResponseRepository using ENT
type EntDonorResponseRepository struct {
	db *ent.Client
}

// Внутренний хелпер для выбора клиента
func (r *EntDonorResponseRepository) client(ctx context.Context) *ent.Client {
	if tx := ent.TxFromContext(ctx); tx != nil {
		return tx.Client()
	}
	return r.db
}

// NewEntDonorResponseRepository creates a new ENT donor response repository
func NewEntDonorResponseRepository(client *ent.Client) *EntDonorResponseRepository {
	return &EntDonorResponseRepository{
		db: client,
	}
}

// toDomainModel converts ENT DonorResponse to domain DonorResponse
func (r *EntDonorResponseRepository) toDomainModel(entResp *ent.DonorResponse) *model.DonorResponse {
	if entResp == nil {
		return nil
	}

	return &model.DonorResponse{
		ID:         entResp.ID,
		RequestID:  entResp.Edges.Request.ID,
		DonorID:    entResp.Edges.Donor.ID,
		Conditions: entResp.Conditions,
		Status:     model.DonorResponseStatus(entResp.Status),
		CreatedAt:  entResp.CreatedAt,
		UpdatedAt:  entResp.UpdatedAt,
	}
}

func (r *EntDonorResponseRepository) CreateDonorResponse(ctx context.Context, reqID, donorID string, conditions []string) (*model.DonorResponse, error) {
	created, err := r.client(ctx).DonorResponse.
		Create().
		SetRequestID(reqID).
		SetDonorID(donorID).
		SetConditions(conditions).
		SetStatus("pending").
		Save(ctx)
	if err != nil {
		return nil, err
	}

	// Load with edges
	entResp, err := r.client(ctx).DonorResponse.Query().Where(donorresponse.ID(created.ID)).WithRequest().WithDonor().Only(ctx)
	if err != nil {
		return nil, err
	}

	return r.toDomainModel(entResp), nil
}

func (r *EntDonorResponseRepository) GetDonorResponseByID(ctx context.Context, id string) (*model.DonorResponse, error) {
	entResp, err := r.client(ctx).DonorResponse.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.toDomainModel(entResp), nil
}

func (r *EntDonorResponseRepository) UpdateDonorResponseStatus(ctx context.Context, id, status string) error {
	return r.client(ctx).DonorResponse.
		UpdateOneID(id).
		SetStatus(status).
		Exec(ctx)
}

func (r *EntDonorResponseRepository) DeleteDonorResponse(ctx context.Context, id string) error {
	return r.client(ctx).DonorResponse.DeleteOneID(id).Exec(ctx)
}

func (r *EntDonorResponseRepository) GetDonorResponsesByRequestID(ctx context.Context, reqID string) ([]*model.DonorResponse, error) {
	entResps, err := r.client(ctx).DonorResponse.
		Query().
		Where(donorresponse.HasRequestWith(bloodsearchrequest.ID(reqID))).
		All(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*model.DonorResponse, len(entResps))
	for i, entResp := range entResps {
		result[i] = r.toDomainModel(entResp)
	}
	return result, nil
}

func (r *EntDonorResponseRepository) GetDonorResponsesByDonorID(ctx context.Context, donorID string) ([]*model.DonorResponse, error) {
	entResps, err := r.client(ctx).DonorResponse.
		Query().
		Where(donorresponse.HasDonorWith(pet.ID(donorID))).
		WithRequest().
		WithDonor().
		All(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*model.DonorResponse, len(entResps))
	for i, entResp := range entResps {
		result[i] = r.toDomainModel(entResp)
	}
	return result, nil
}

func (r *EntDonorResponseRepository) ExistsByID(ctx context.Context, id string) (bool, error) {
	return r.client(ctx).DonorResponse.Query().
		Where(donorresponse.ID(id)).
		Exist(ctx)
}

func (r *EntDonorResponseRepository) ExistsByRequestID(ctx context.Context, reqID string) (bool, error) {
	return r.client(ctx).DonorResponse.Query().
		Where(donorresponse.HasRequestWith(bloodsearchrequest.ID(reqID))).
		Exist(ctx)
}

func (r *EntDonorResponseRepository) ExistsByDonorID(ctx context.Context, donorID string) (bool, error) {
	return r.client(ctx).DonorResponse.Query().
		Where(donorresponse.HasDonorWith(pet.ID(donorID))).
		Exist(ctx)
}

func (r *EntDonorResponseRepository) Count(ctx context.Context) (int, error) {
	return r.client(ctx).DonorResponse.Query().Count(ctx)
}
