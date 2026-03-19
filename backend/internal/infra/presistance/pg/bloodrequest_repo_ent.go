package pg

import (
	"context"
	"errors"
	"fmt"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	bloodreqmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	petmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/bloodsearchrequest"
)

// EntBloodRequestRepository implements BloodRequestRepository using ENT
type EntBloodRequestRepository struct {
	db *ent.Client
}

// Внутренний хелпер для выбора клиента
func (r *EntBloodRequestRepository) client(ctx context.Context) *ent.Client {
	if tx := ent.TxFromContext(ctx); tx != nil {
		return tx.Client()
	}
	return r.db
}

// NewEntBloodRequestRepository creates a new ENT blood request repository
func NewEntBloodRequestRepository(client *ent.Client) *EntBloodRequestRepository {
	return &EntBloodRequestRepository{
		db: client,
	}
}

// mapToRecipient maps ent.BloodSearchRequest to donormodel.Recipient
func (r *EntBloodRequestRepository) mapToRecipient(req *ent.BloodSearchRequest, donors []*petmodel.Pet) *donormodel.Recipient {
	if req == nil {
		return nil
	}

	recipient := &donormodel.Recipient{
		ID:                       req.ID,
		PetID:                    req.PetID,
		BloodVolumeRemaining:     req.BloodVolumeNeeded - req.BloodVolumeReserved,
		PrioritySearch:           req.PrioritySearch,
		IncludeUnknownBloodGroup: req.IncludeUnknownBloodGroup,
		SearchingBloodNames:      req.BloodGroupNames,
		Status:                   string(req.Status),
	}

	if req.Edges.Pet != nil {
		recipient.PetName = req.Edges.Pet.Name
		recipient.PetType = petmodel.PetType(req.Edges.Pet.Type)
		recipient.PhotoURLs = req.Edges.Pet.PhotoUrls
		if req.Edges.Pet.Edges.BloodGroupRef != nil {
			recipient.BloodGroupName = req.Edges.Pet.Edges.BloodGroupRef.BloodGroup
		}
	}

	// Find matching donors
	for _, donor := range donors {
		recipient.AddMatchingDonor(donor)
	}

	return recipient
}

// bloodReqToDomainModel converts ENT BloodSearchRequest to domain BloodRequest
func (r *EntBloodRequestRepository) bloodReqToDomainModel(entReq *ent.BloodSearchRequest) *bloodreqmodel.BloodRequest {
	if entReq == nil {
		return nil
	}

	// Map responses to DonorApplications
	var donorApps []bloodreqmodel.DonorApplication
	if entReq.Edges.Responses != nil {
		donorApps = make([]bloodreqmodel.DonorApplication, len(entReq.Edges.Responses))
		for i, resp := range entReq.Edges.Responses {
			app := bloodreqmodel.DonorApplication{
				ID:               resp.ID,
				RequestID:        entReq.ID,
				DonorID:          resp.Edges.Donor.ID,
				DonorName:        resp.Edges.Donor.Name,
				DonorPhotos:      resp.Edges.Donor.PhotoUrls,
				DonorBloodGroup:  resp.Edges.Donor.Edges.BloodGroupRef.BloodGroup,
				Amount:           0,
				WarnFactors:      []string{},
				CompensationType: string(resp.CompensationType),
				TaxiCompensation: resp.TaxiCompensation,
			}
			donorApps[i] = app
		}
	}

	return &bloodreqmodel.BloodRequest{
		ID:                       entReq.ID,
		PetID:                    entReq.PetID,
		BloodVolumeNeeded:        entReq.BloodVolumeNeeded,
		BloodVolumeReserved:      entReq.BloodVolumeReserved,
		Regions:                  entReq.Regions,
		SmallPetsNotifyAllowed:   entReq.SmallPetsNotifyAllowed,
		Status:                   bloodreqmodel.BloodRequestStatus(entReq.Status),
		Description:              entReq.Description,
		PhotoURLs:                entReq.PhotoUrls,
		BloodGroupNames:          entReq.BloodGroupNames,
		BloodComponentIDs:        entReq.BloodComponentIds,
		OnBoarding:               entReq.OnBoarding,
		DonorApplications:        donorApps,
		PrioritySearch:           entReq.PrioritySearch,
		IncludeUnknownBloodGroup: entReq.IncludeUnknownBloodGroup,
		CreatedAt:                &entReq.CreatedAt,
		UpdatedAt:                &entReq.UpdatedAt,
		DeletedAt:                entReq.DeletedAt,
	}
}

// Create создает новую заявку на поиск крови
func (r *EntBloodRequestRepository) Create(ctx context.Context, req *bloodreqmodel.BloodRequest) (*bloodreqmodel.BloodRequest, error) {
	newBloodReq, err := r.client(ctx).BloodSearchRequest.
		Create().
		SetPetID(req.PetID).
		SetBloodVolumeNeeded(req.BloodVolumeNeeded).
		SetBloodVolumeReserved(req.BloodVolumeReserved).
		SetRegions(req.Regions).
		SetSmallPetsNotifyAllowed(req.SmallPetsNotifyAllowed).
		SetStatus(bloodsearchrequest.Status(req.Status)).
		SetDescription(req.Description).
		SetPhotoUrls(req.PhotoURLs).
		SetBloodGroupNames(req.BloodGroupNames).
		SetBloodComponentIds(req.BloodComponentIDs).
		SetOnBoarding(req.OnBoarding).
		SetPrioritySearch(req.PrioritySearch).
		SetIncludeUnknownBloodGroup(req.IncludeUnknownBloodGroup).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return r.bloodReqToDomainModel(newBloodReq), nil
}

// GetByID возвращает заявку по её идентификатору
func (r *EntBloodRequestRepository) GetByID(ctx context.Context, id string) (*bloodreqmodel.BloodRequest, error) {
	reqQuery := r.client(ctx).BloodSearchRequest.Query().
		Where(bloodsearchrequest.ID(id)).
		WithResponses(func(drq *ent.DonorResponseQuery) {
			drq.WithDonor(func(pq *ent.PetQuery) {
				pq.WithBloodGroupRef()
			})
		})

	req, err := reqQuery.Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperrors.ErrBloodRequestNotFound
		}
		return nil, apperrors.Internal(err, "failed to execute blood request query")
	}
	return r.bloodReqToDomainModel(req), nil
}

// GetByPetID возвращает заявку по идентификатору питомца
func (r *EntBloodRequestRepository) GetByPetID(ctx context.Context, petID string) (*bloodreqmodel.BloodRequest, error) {
	req, err := r.client(ctx).BloodSearchRequest.Query().
		Where(bloodsearchrequest.PetID(petID)).
		WithResponses(func(drq *ent.DonorResponseQuery) {
			drq.WithDonor(func(pq *ent.PetQuery) {
				pq.WithBloodGroupRef()
			})
		}).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperrors.ErrBloodRequestNotFound
		}
		return nil, apperrors.Internal(err, "failed to execute blood request query by pet ID")
	}
	return r.bloodReqToDomainModel(req), nil
}

// Update обновляет информацию о заявке
func (r *EntBloodRequestRepository) Update(ctx context.Context, id string, req *bloodreqmodel.BloodRequest) (*bloodreqmodel.BloodRequest, error) {
	if req == nil {
		return nil, errors.New("blood request cannot be nil")
	}

	if id == "" {
		return nil, errors.New("invalid blood request ID")
	}

	updater := r.client(ctx).BloodSearchRequest.UpdateOneID(id).
		SetBloodVolumeNeeded(req.BloodVolumeNeeded).
		SetBloodVolumeReserved(req.BloodVolumeReserved).
		SetRegions(req.Regions).
		SetSmallPetsNotifyAllowed(req.SmallPetsNotifyAllowed).
		SetStatus(bloodsearchrequest.Status(req.Status)).
		SetDescription(req.Description).
		SetPhotoUrls(req.PhotoURLs).
		SetBloodGroupNames(req.BloodGroupNames).
		SetBloodComponentIds(req.BloodComponentIDs).
		SetOnBoarding(req.OnBoarding).
		SetPrioritySearch(req.PrioritySearch).
		SetIncludeUnknownBloodGroup(req.IncludeUnknownBloodGroup)

	updatedBloodReq, err := updater.Save(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperrors.ErrBloodRequestNotFound
		}
		return nil, apperrors.Internal(err, "failed to update blood request")
	}

	return r.bloodReqToDomainModel(updatedBloodReq), nil
}

// UpdateStatus обновляет статус заявки
func (r *EntBloodRequestRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	return r.client(ctx).BloodSearchRequest.UpdateOneID(id).
		SetStatus(bloodsearchrequest.Status(status)).
		Exec(ctx)
}

// UpdateReservedVolume обновляет зарезервированный объём и статус заявки
func (r *EntBloodRequestRepository) UpdateReservedVolume(ctx context.Context, id string, reservedVolume int32, status string) error {
	return r.client(ctx).BloodSearchRequest.UpdateOneID(id).
		SetBloodVolumeReserved(reservedVolume).
		SetStatus(bloodsearchrequest.Status(status)).
		Exec(ctx)
}

// Delete удаляет заявку из хранилища (soft delete)
func (r *EntBloodRequestRepository) Delete(ctx context.Context, id string) error {
	return r.client(ctx).BloodSearchRequest.DeleteOneID(id).Exec(ctx)
}

// List возвращает список заявок с фильтрацией и пагинацией
func (r *EntBloodRequestRepository) List(ctx context.Context, filters donormodel.DonorPreloadFilter) ([]*bloodreqmodel.BloodRequest, error) {
	requests, err := r.client(ctx).BloodSearchRequest.Query().
		Where(
			bloodsearchrequest.StatusEQ(bloodsearchrequest.Status(filters.Status)),
		).
		WithResponses().
		Limit(filters.Limit).
		Offset(filters.Offset).
		All(ctx)

	if err != nil {
		return nil, err
	}

	result := make([]*bloodreqmodel.BloodRequest, len(requests))
	for i, req := range requests {
		result[i] = r.bloodReqToDomainModel(req)
	}

	return result, nil
}

// AdptiveList возвращает список заявок с фильтрацией и пагинацией
func (r *EntBloodRequestRepository) AdptiveList(ctx context.Context, donors []*petmodel.Pet, filters donormodel.DonorPreloadFilter) ([]*donormodel.Recipient, error) {
	requests, err := r.client(ctx).BloodSearchRequest.Query().
		Where(bloodsearchrequest.StatusEQ(bloodsearchrequest.Status(filters.Status))).
		WithPet(func(pq *ent.PetQuery) {
			pq.WithBloodGroupRef()
		}).
		Limit(filters.Limit).
		Offset(filters.Offset).
		All(ctx)

	if err != nil {
		return nil, err
	}

	result := make([]*donormodel.Recipient, len(requests))
	for i, req := range requests {
		result[i] = r.mapToRecipient(req, donors)
	}

	return result, nil
}

// ExistsByPetID проверяет существование активной заявки для питомца
func (r *EntBloodRequestRepository) ExistsByPetID(ctx context.Context, petID string) (bool, error) {
	return r.client(ctx).BloodSearchRequest.Query().
		Where(
			bloodsearchrequest.PetID(petID),
			bloodsearchrequest.StatusEQ(bloodsearchrequest.StatusActive),
		).
		Exist(ctx)
}

// ExistsByID проверяет существование заявки по её идентификатору
func (r *EntBloodRequestRepository) ExistsByID(ctx context.Context, id string) (bool, error) {
	return r.client(ctx).BloodSearchRequest.Query().
		Where(bloodsearchrequest.ID(id)).
		Exist(ctx)
}

// Count возвращает общее количество заявок в хранилище
func (r *EntBloodRequestRepository) Count(ctx context.Context) (int, error) {
	return r.client(ctx).BloodSearchRequest.Query().Count(ctx)
}

// AddPhotoURLs adds new photo paths to the blood request's PhotoUrls array
func (r *EntBloodRequestRepository) AddPhotoURLs(ctx context.Context, id string, paths []string) error {
	if id == "" {
		return errors.New("invalid blood request ID")
	}

	err := r.client(ctx).BloodSearchRequest.UpdateOneID(id).
		SetPhotoUrls(paths).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to update blood request photo URLs: %w", err)
	}

	return nil
}
