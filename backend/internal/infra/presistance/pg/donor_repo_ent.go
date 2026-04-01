package pg

import (
	"context"
	"fmt"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	petmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
	recipientmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/recipient/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/bloodsearchrequest"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/donorresponse"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/presistance/domainmapper"
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

func (r *EntDonorResponseRepository) CreateDonorResponse(ctx context.Context, resp *donormodel.DonorResponse) (*donormodel.DonorResponse, error) {
	created, err := r.client(ctx).DonorResponse.
		Create().
		SetRequestID(resp.RequestID).
		SetDonorID(resp.DonorID).
		SetAmount(resp.Amount).
		SetCompensationType(donorresponse.CompensationType(resp.CompensationType)).
		SetTaxiCompensation(resp.TaxiCompensation).
		SetStatus(donorresponse.Status(resp.Status)).
		Save(ctx)
	if err != nil {
		return nil, err
	}

	// Load with edges
	entResp, err := r.client(ctx).DonorResponse.Query().Where(donorresponse.ID(created.ID)).WithRequest().WithDonor().Only(ctx)
	if err != nil {
		return nil, err
	}

	return domainmapper.ApplicationToDomain(entResp), nil
}

func (r *EntDonorResponseRepository) GetDonorResponseByID(ctx context.Context, id string) (*donormodel.DonorResponse, error) {
	entResp, err := r.client(ctx).DonorResponse.Query().
		Where(donorresponse.ID(id)).
		WithRequest(func(q *ent.BloodSearchRequestQuery) { q.Select(bloodsearchrequest.FieldID) }).
		WithDonor(func(q *ent.PetQuery) { q.Select(pet.FieldID) }).
		Only(ctx)
	if err != nil {
		return nil, err
	}
	return domainmapper.ApplicationToDomain(entResp), nil
}

func (r *EntDonorResponseRepository) GetRecipient(ctx context.Context, id string) (*recipientmodel.Recipient, error) {
	blreq, err := r.db.BloodSearchRequest.Query().
		Where(bloodsearchrequest.IDEQ(id)).
		WithPet(func(pq *ent.PetQuery) {
			pq.WithBloodGroupRef()
			pq.WithOwner(
				func(uq *ent.UserQuery) {
					uq.WithDonorPreference()
				},
			)
		}).
		Only(ctx)
	if err != nil {
		return nil, err
	}

	recipient := &recipientmodel.Recipient{
		ID:                       blreq.ID,
		PetID:                    blreq.PetID,
		PetName:                  blreq.Edges.Pet.Name,
		SearchingBloodNames:      blreq.BloodGroupNames,
		PetType:                  petmodel.PetType(blreq.Edges.Pet.Type),
		SearchRegions:            blreq.Regions,
		BloodGroupName:           blreq.Edges.Pet.Edges.BloodGroupRef.BloodGroup,
		PhotoURLs:                blreq.Edges.Pet.PhotoUrls,
		BloodVolumeNeeded:        blreq.BloodVolumeNeeded,
		BloodVolumeReserved:      blreq.BloodVolumeReserved,
		PrioritySearch:           blreq.PrioritySearch,
		IncludeUnknownBloodGroup: blreq.IncludeUnknownBloodGroup,
		SmallPetsNotifyAllowed:   blreq.SmallPetsNotifyAllowed,
		OwnerName:                blreq.Edges.Pet.Edges.Owner.FullName,
		Status:                   string(blreq.Status),
		AdvancedInfo: &recipientmodel.AdvancedInfo{
			Description: blreq.Description,
			PhotoURLs:   blreq.PhotoUrls,
		},
	}

	return recipient, nil
}

func (r *EntDonorResponseRepository) GetByPetID(ctx context.Context, petID string) (*donormodel.DonorResponse, error) {
	entResp, err := r.client(ctx).DonorResponse.Query().
		Where(donorresponse.HasDonorWith(pet.ID(petID))).
		WithRequest(func(q *ent.BloodSearchRequestQuery) { q.Select(bloodsearchrequest.FieldID) }).
		WithDonor(func(q *ent.PetQuery) { q.Select(pet.FieldID) }).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperrors.ErrDonorResponseNotFound
		}
		return nil, err
	}

	return domainmapper.ApplicationToDomain(entResp), nil
}

func (r *EntDonorResponseRepository) UpdateDonorResponseStatus(ctx context.Context, id string, status donormodel.DonorResponseStatus) error {
	return r.client(ctx).DonorResponse.
		UpdateOneID(id).
		SetStatus(donorresponse.Status(status)).
		Exec(ctx)
}

func (r *EntDonorResponseRepository) DeleteDonorResponse(ctx context.Context, id string) error {
	return r.client(ctx).DonorResponse.DeleteOneID(id).Exec(ctx)
}

func (r *EntDonorResponseRepository) GetDonorResponsesByRequestID(ctx context.Context, reqID string) ([]*donormodel.DonorResponse, error) {
	entResps, err := r.client(ctx).DonorResponse.
		Query().
		Where(donorresponse.HasRequestWith(bloodsearchrequest.ID(reqID))).
		WithRequest(func(q *ent.BloodSearchRequestQuery) { q.Select(bloodsearchrequest.FieldID) }).
		WithDonor(func(q *ent.PetQuery) { q.Select(pet.FieldID) }).
		All(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*donormodel.DonorResponse, len(entResps))
	for i, entResp := range entResps {
		result[i] = domainmapper.ApplicationToDomain(entResp)
	}
	return result, nil
}

func (r *EntDonorResponseRepository) GetDonorResponsesByDonorID(ctx context.Context, donorID string) ([]*donormodel.DonorResponse, error) {
	entResps, err := r.client(ctx).DonorResponse.
		Query().
		Where(donorresponse.HasDonorWith(pet.ID(donorID))).
		WithRequest(func(q *ent.BloodSearchRequestQuery) { q.Select(bloodsearchrequest.FieldID) }).
		WithDonor(func(q *ent.PetQuery) { q.Select(pet.FieldID) }).
		All(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]*donormodel.DonorResponse, len(entResps))
	for i, entResp := range entResps {
		result[i] = domainmapper.ApplicationToDomain(entResp)
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

// Подтверждение донации реципиентом
func (r *EntDonorResponseRepository) Confirm(ctx context.Context, donorResponseID string, factAmount int32) error {
	update := r.client(ctx).DonorResponse.
		UpdateOneID(donorResponseID).
		SetIsConfirmed(true)

	if factAmount != 0 {
		update.SetAmount(factAmount)
	}

	if err := update.Exec(ctx); err != nil {
		return fmt.Errorf("failed to confirm blood request: %w", err)
	}

	return nil
}
