package pg

import (
	"context"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	bloodreqmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/common"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
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

func (r *EntDonorResponseRepository) GetRecipient(ctx context.Context, id string) (*bloodreqmodel.BloodRequestWithMatchingDonors, error) {
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

	recipient := &bloodreqmodel.BloodRequestWithMatchingDonors{
		BloodRequest: bloodreqmodel.BloodRequest{
			ID:                       blreq.ID,
			PetID:                    blreq.PetID,
			BloodGroupNames:          blreq.BloodGroupNames,
			Regions:                  blreq.Regions,
			BloodVolumeNeeded:        blreq.BloodVolumeNeeded,
			PrioritySearch:           blreq.PrioritySearch,
			IncludeUnknownBloodGroup: blreq.IncludeUnknownBloodGroup,
			SmallPetsNotifyAllowed:   blreq.SmallPetsNotifyAllowed,
			Status:                   bloodreqmodel.BloodRequestStatus(blreq.Status),
			AdvancedInfo: bloodreqmodel.AdvancedInfo{
				Description: blreq.Description,
				PhotoURLs:   blreq.PhotoUrls,
			},
		},
		RecipientData: bloodreqmodel.RecipientData{
			PetName:        blreq.Edges.Pet.Name,
			PetType:        common.PetType(blreq.Edges.Pet.Type),
			BloodGroupName: blreq.Edges.Pet.Edges.BloodGroupRef.BloodGroup,
			OwnerName:      blreq.Edges.Pet.Edges.Owner.FullName,
			PhotoURLs:      blreq.Edges.Pet.PhotoUrls,
		},
	}

	return recipient, nil
}

func (r *EntDonorResponseRepository) GetByPetID(ctx context.Context, petID string) (*donormodel.DonorResponse, error) {
	entResp, err := r.client(ctx).DonorResponse.Query().
		Where(donorresponse.HasDonorWith(pet.ID(petID))).
		Order(donorresponse.ByCreatedAt(sql.OrderDesc())).
		WithRequest(func(q *ent.BloodSearchRequestQuery) { q.Select(bloodsearchrequest.FieldID) }).
		WithDonor(func(q *ent.PetQuery) { q.Select(pet.FieldID) }).
		First(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperrors.ErrDonorResponseNotFound
		}
		return nil, err
	}

	return domainmapper.ApplicationToDomain(entResp), nil
}

// GetByPetIDs возвращает мапу слайсов откликов доноров по идентификаторам питомцев (все отклики, отсортированные по дате создания DESC)
func (r *EntDonorResponseRepository) GetByPetIDs(ctx context.Context, petIDs []string) (map[string][]*donormodel.DonorResponse, error) {
	if len(petIDs) == 0 {
		return make(map[string][]*donormodel.DonorResponse), nil
	}

	entResps, err := r.client(ctx).DonorResponse.Query().
		Where(donorresponse.HasDonorWith(pet.IDIn(petIDs...))).
		Order(donorresponse.ByCreatedAt(sql.OrderDesc())).
		WithRequest(func(q *ent.BloodSearchRequestQuery) { q.Select(bloodsearchrequest.FieldID) }).
		WithDonor(func(q *ent.PetQuery) { q.Select(pet.FieldID) }).
		All(ctx)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to execute donor response query by pet IDs")
	}

	result := make(map[string][]*donormodel.DonorResponse)
	for _, entResp := range entResps {
		petID := entResp.Edges.Donor.ID
		result[petID] = append(result[petID], domainmapper.ApplicationToDomain(entResp))
	}
	return result, nil
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
func (r *EntDonorResponseRepository) Confirm(ctx context.Context, donorResponseID string, factAmount float64) error {
	update := r.client(ctx).DonorResponse.
		UpdateOneID(donorResponseID).
		SetStatus(donorresponse.StatusCompleted).
		SetIsConfirmed(true)

	if factAmount != 0 {
		update.SetAmount(math.Round(factAmount*10) / 10)
	}

	if err := update.Exec(ctx); err != nil {
		return fmt.Errorf("failed to confirm blood request: %w", err)
	}

	return nil
}

// Завершение донации донором
func (r *EntDonorResponseRepository) Complete(ctx context.Context, donorResponseID string, factAmount float64) error {
	update := r.client(ctx).DonorResponse.
		UpdateOneID(donorResponseID).
		SetStatus(donorresponse.StatusCompleted)

	if factAmount != 0 {
		update.SetAmount(math.Round(factAmount*10) / 10)
	}

	if err := update.Exec(ctx); err != nil {
		return fmt.Errorf("failed to complete donation: %w", err)
	}

	return nil
}

// Reject отклоняет отклик донора с причиной
func (r *EntDonorResponseRepository) Reject(ctx context.Context, req *donormodel.DonorResponse) error {
	update := r.client(ctx).DonorResponse.
		UpdateOneID(req.ID).
		SetStatus(donorresponse.Status(req.Status)).
		SetRejectedReason(req.RejectedReason)

	if err := update.Exec(ctx); err != nil {
		return fmt.Errorf("failed to reject donor response: %w", err)
	}

	return nil
}

// Cancel отменяет отклик донора
func (r *EntDonorResponseRepository) Cancel(ctx context.Context, donorResponseID string) error {
	update := r.client(ctx).DonorResponse.
		UpdateOneID(donorResponseID).
		SetStatus(donorresponse.StatusCancelled)

	if err := update.Exec(ctx); err != nil {
		return fmt.Errorf("failed to cancel donor response: %w", err)
	}

	return nil
}

// Accept accepts a donor response
func (r *EntDonorResponseRepository) Accept(ctx context.Context, donorResponseID string) error {
	update := r.client(ctx).DonorResponse.
		UpdateOneID(donorResponseID).
		SetStatus(donorresponse.StatusAccepted)

	if err := update.Exec(ctx); err != nil {
		return fmt.Errorf("failed to accept donor response: %w", err)
	}

	return nil
}
