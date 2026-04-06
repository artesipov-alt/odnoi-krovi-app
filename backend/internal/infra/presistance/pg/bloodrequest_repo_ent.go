package pg

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	bloodreqmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/bloodsearchrequest"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/donorresponse"
	entuser "github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/user"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/presistance/domainmapper"
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

// Create создает новую заявку на поиск крови
func (r *EntBloodRequestRepository) Create(ctx context.Context, req *bloodreqmodel.BloodRequest) (*bloodreqmodel.BloodRequestWithApplications, error) {
	newBloodReq, err := r.client(ctx).BloodSearchRequest.
		Create().
		SetPetID(req.PetID).
		SetBloodVolumeNeeded(req.BloodVolumeNeeded).
		SetRegions(req.Regions).
		SetSmallPetsNotifyAllowed(req.SmallPetsNotifyAllowed).
		SetStatus(bloodsearchrequest.Status(req.Status)).
		SetDescription(req.AdvancedInfo.Description).
		SetPhotoUrls(req.AdvancedInfo.PhotoURLs).
		SetBloodGroupNames(req.BloodGroupNames).
		SetBloodComponentIds(req.BloodComponentIDs).
		SetOnBoarding(req.OnBoarding).
		SetPrioritySearch(req.PrioritySearch).
		SetIncludeUnknownBloodGroup(req.IncludeUnknownBloodGroup).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	return domainmapper.BloodReqToDomain(newBloodReq), nil
}

// GetByID возвращает заявку по её идентификатору
func (r *EntBloodRequestRepository) GetByID(ctx context.Context, id string) (*bloodreqmodel.BloodRequestWithApplications, error) {
	reqQuery := r.client(ctx).BloodSearchRequest.Query().
		Where(bloodsearchrequest.ID(id)).
		WithResponses(func(drq *ent.DonorResponseQuery) {
			drq.WithDonor(func(pq *ent.PetQuery) {
				pq.WithBloodGroupRef()
				pq.WithOwner(func(uq *ent.UserQuery) {
					uq.Select(entuser.FieldFullName)
				})
			})
		})

	req, err := reqQuery.Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperrors.ErrBloodRequestNotFound
		}
		return nil, apperrors.Internal(err, "failed to execute blood request query")
	}
	return domainmapper.BloodReqToDomain(req), nil
}

// GetByPetID возвращает заявку по идентификатору питомца
func (r *EntBloodRequestRepository) GetByPetID(ctx context.Context, petID string) (*bloodreqmodel.BloodRequestWithApplications, error) {
	req, err := r.client(ctx).BloodSearchRequest.Query().
		Where(
			bloodsearchrequest.PetID(petID),
		).
		Order(bloodsearchrequest.ByCreatedAt(sql.OrderDesc())).
		WithResponses(func(drq *ent.DonorResponseQuery) {
			drq.WithDonor(func(pq *ent.PetQuery) {
				//Возвращаем полного донора, чтобы пересчитать warn-факторы.
				pq.WithBloodGroupRef()
				pq.WithHealth()
				pq.WithTreatments()
				pq.WithAnalyses()
				pq.WithOwner(func(uq *ent.UserQuery) {
					uq.Select(entuser.FieldFullName)
				})
			})
		}).
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperrors.ErrBloodRequestNotFound
		}
		return nil, apperrors.Internal(err, "failed to execute blood request query by pet ID")
	}
	return domainmapper.BloodReqToDomain(req), nil
}

// GetByApplicationID возвращает заявку по id отклика на эту заявку
func (r *EntBloodRequestRepository) GetByApplicationID(ctx context.Context, id string) (*bloodreqmodel.BloodRequestWithApplications, error) {
	req, err := r.client(ctx).BloodSearchRequest.Query().
		Where(bloodsearchrequest.HasResponsesWith(donorresponse.IDEQ(id))).
		WithResponses().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get blood request by application ID: %w", err)
	}
	return domainmapper.BloodReqToDomain(req), nil
}

// Update обновляет информацию о заявке
func (r *EntBloodRequestRepository) Update(ctx context.Context, id string, req *bloodreqmodel.BloodRequest) (*bloodreqmodel.BloodRequestWithApplications, error) {
	if req == nil {
		return nil, errors.New("blood request cannot be nil")
	}

	if id == "" {
		return nil, errors.New("invalid blood request ID")
	}

	updater := r.client(ctx).BloodSearchRequest.UpdateOneID(id).
		SetBloodVolumeNeeded(req.BloodVolumeNeeded).
		SetRegions(req.Regions).
		SetSmallPetsNotifyAllowed(req.SmallPetsNotifyAllowed).
		SetStatus(bloodsearchrequest.Status(req.Status)).
		SetDescription(req.AdvancedInfo.Description).
		SetPhotoUrls(req.AdvancedInfo.PhotoURLs).
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

	return domainmapper.BloodReqToDomain(updatedBloodReq), nil
}

// UpdateStatus обновляет статус заявки
func (r *EntBloodRequestRepository) UpdateStatus(ctx context.Context, id string, status bloodreqmodel.BloodRequestStatus) error {
	return r.client(ctx).BloodSearchRequest.UpdateOneID(id).
		SetStatus(bloodsearchrequest.Status(status)).
		Exec(ctx)
}

// Delete удаляет заявку из хранилища (soft delete)
func (r *EntBloodRequestRepository) Delete(ctx context.Context, id string) error {
	return r.client(ctx).BloodSearchRequest.DeleteOneID(id).Exec(ctx)
}

// List возвращает список заявок с фильтрацией и пагинацией
func (r *EntBloodRequestRepository) List(ctx context.Context, filters donormodel.DonorPreloadFilter) ([]*bloodreqmodel.BloodRequestWithApplications, error) {
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

	result := make([]*bloodreqmodel.BloodRequestWithApplications, len(requests))
	for i, req := range requests {
		result[i] = domainmapper.BloodReqToDomain(req)
	}

	return result, nil
}

// AdptiveList возвращает список заявок с фильтрацией и пагинацией
func (r *EntBloodRequestRepository) AdaptiveList(ctx context.Context, filters donormodel.DonorPreloadFilter) ([]*bloodreqmodel.BloodRequestWithMatchingDonors, error) {
	requests, err := r.client(ctx).BloodSearchRequest.Query().
		Where(
			bloodsearchrequest.StatusEQ(bloodsearchrequest.Status(filters.Status)),
		).
		WithPet(func(pq *ent.PetQuery) {
			pq.WithBloodGroupRef()
		}).
		Limit(filters.Limit).
		Offset(filters.Offset).
		All(ctx)

	if err != nil {
		return nil, err
	}

	result := make([]*bloodreqmodel.BloodRequestWithMatchingDonors, len(requests))
	for i, req := range requests {
		result[i] = domainmapper.RecipientToDomain(req)
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
