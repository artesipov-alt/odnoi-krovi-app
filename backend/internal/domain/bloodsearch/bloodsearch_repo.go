package bloodsearch

import (
	"context"

	bloodreqmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
)

// BloodRequestRepository определяет интерфейс для работы с данными заявок на поиск крови питомцев
type BloodRequestRepository interface {
	// Write methods

	// создает новую заявку на поиск крови
	Create(ctx context.Context, req *bloodreqmodel.BloodRequest) (*bloodreqmodel.BloodRequestWithApplications, error)

	// обновляет информацию о заявке
	Update(ctx context.Context, id string, req *bloodreqmodel.BloodRequest) (*bloodreqmodel.BloodRequestWithApplications, error)

	// обновляет статус заявки
	UpdateStatus(ctx context.Context, id string, status bloodreqmodel.BloodRequestStatus) error

	// удаляет заявку (soft delete)
	Delete(ctx context.Context, id string) error

	// добавляет пути к фотографиям заявки
	AddPhotoURLs(ctx context.Context, id string, paths []string) error

	// Read methods

	// возвращает заявку по ID
	GetByID(ctx context.Context, id string) (*bloodreqmodel.BloodRequestWithApplications, error)

	// возвращает заявку по ID отклика на нее
	GetByApplicationID(ctx context.Context, id string, ignoreSoftDelete bool) (*bloodreqmodel.BloodRequestWithApplications, error)

	// возвращает заявку по ID питомца
	GetByPetID(ctx context.Context, petID string) (*bloodreqmodel.BloodRequestWithApplications, error)

	// возвращает заявки по списку ID питомцев
	GetByPetIDs(ctx context.Context, petIDs []string) (map[string]*bloodreqmodel.BloodRequestWithApplications, error)

	// возвращает список заявок с фильтрацией и пагинацией
	List(ctx context.Context, filters donormodel.DonorPreloadFilter) ([]*bloodreqmodel.BloodRequestWithApplications, error)

	// возвращает список заявок с подходящими донорами
	AdaptiveList(ctx context.Context, filters donormodel.DonorPreloadFilter) ([]*bloodreqmodel.BloodRequestWithMatchingDonors, error)

	// проверяет, существует ли активная заявка для питомца
	ExistsByPetID(ctx context.Context, petID string) (bool, error)

	// проверяет, существует ли заявка по ID
	ExistsByID(ctx context.Context, id string) (bool, error)

	// возвращает общее количество заявок
	Count(ctx context.Context) (int, error)
}
