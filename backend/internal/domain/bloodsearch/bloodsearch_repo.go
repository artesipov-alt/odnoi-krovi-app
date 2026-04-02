package bloodsearch

import (
	"context"

	bloodreqmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
)

// BloodRequestRepository определяет интерфейс для работы с данными заявок на поиск крови питомцев
type BloodRequestRepository interface {
	// Create создает новую заявку на поиск крови
	Create(ctx context.Context, req *bloodreqmodel.BloodRequest) (*bloodreqmodel.BloodRequestWithApplications, error)

	// GetByID возвращает заявку по её идентификатору
	GetByID(ctx context.Context, id string) (*bloodreqmodel.BloodRequestWithApplications, error)

	// GetByApplicationID возвращает заявку по id отклика на эту заявку
	GetByApplicationID(ctx context.Context, id string) (*bloodreqmodel.BloodRequestWithApplications, error)

	// GetByPetID возвращает заявку по идентификатору питомца
	GetByPetID(ctx context.Context, petID string) (*bloodreqmodel.BloodRequestWithApplications, error)

	// Update обновляет информацию о заявке
	Update(ctx context.Context, id string, req *bloodreqmodel.BloodRequest) (*bloodreqmodel.BloodRequestWithApplications, error)

	// UpdateStatus обновляет статус заявки
	UpdateStatus(ctx context.Context, id string, status bloodreqmodel.BloodRequestStatus) error

	// Delete удаляет заявку из хранилища (soft delete)
	Delete(ctx context.Context, id string) error

	// List возвращает список заявок с фильтрацией и пагинацией
	List(ctx context.Context, filters donormodel.DonorPreloadFilter) ([]*bloodreqmodel.BloodRequestWithApplications, error)

	AdaptiveList(ctx context.Context, filters donormodel.DonorPreloadFilter) ([]*bloodreqmodel.BloodRequestWithMatchingDonors, error)

	// ExistsByPetID проверяет существование активной заявки для питомца
	ExistsByPetID(ctx context.Context, petID string) (bool, error)

	// ExistsByID проверяет существование заявки по её идентификатору
	ExistsByID(ctx context.Context, id string) (bool, error)

	// Count возвращает общее количество заявок в хранилище
	Count(ctx context.Context) (int, error)

	// AddPhotoURLs добавляет новые пути к фотографиям заявки
	AddPhotoURLs(ctx context.Context, id string, paths []string) error
}
