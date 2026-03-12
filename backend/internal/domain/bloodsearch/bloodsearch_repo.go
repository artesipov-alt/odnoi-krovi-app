package bloodsearch

import (
	"context"

	bloodreqmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	petmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
)

// BloodRequestRepository определяет интерфейс для работы с данными заявок на поиск крови питомцев
type BloodRequestRepository interface {
	// Create создает новую заявку на поиск крови
	Create(ctx context.Context, req *bloodreqmodel.BloodRequest) (*bloodreqmodel.BloodRequest, error)

	// GetByID возвращает заявку по её идентификатору
	GetByID(ctx context.Context, id string) (*bloodreqmodel.BloodRequest, error)

	// GetByPetID возвращает заявку по идентификатору питомца
	GetByPetID(ctx context.Context, petID string) (*bloodreqmodel.BloodRequest, error)

	// Update обновляет информацию о заявке
	Update(ctx context.Context, id string, req *bloodreqmodel.BloodRequest) (*bloodreqmodel.BloodRequest, error)

	// UpdateStatus обновляет статус заявки
	UpdateStatus(ctx context.Context, id string, status string) error

	// Delete удаляет заявку из хранилища (soft delete)
	Delete(ctx context.Context, id string) error

	// List возвращает список заявок с фильтрацией и пагинацией
	List(ctx context.Context, filters donormodel.DonorPreloadFilter) ([]*bloodreqmodel.BloodRequest, error)

	AdptiveList(ctx context.Context, donors []*petmodel.Pet, filters donormodel.DonorPreloadFilter) ([]*donormodel.Recipient, error)

	// ExistsByPetID проверяет существование активной заявки для питомца
	ExistsByPetID(ctx context.Context, petID string) (bool, error)

	// ExistsByID проверяет существование заявки по её идентификатору
	ExistsByID(ctx context.Context, id string) (bool, error)

	// Count возвращает общее количество заявок в хранилище
	Count(ctx context.Context) (int, error)

	// AddPhotoURLs добавляет новые пути к фотографиям заявки
	AddPhotoURLs(ctx context.Context, id string, paths []string) error
}

// DonorResponseRepository определяет интерфейс для работы с откликами доноров
type DonorResponseRepository interface {
	CreateDonorResponse(ctx context.Context, resp *bloodreqmodel.DonorResponse) (*bloodreqmodel.DonorResponse, error)
	GetDonorResponseByID(ctx context.Context, id string) (*bloodreqmodel.DonorResponse, error)
	UpdateDonorResponseStatus(ctx context.Context, id, status string) error
	DeleteDonorResponse(ctx context.Context, id string) error
	GetDonorResponsesByRequestID(ctx context.Context, reqID string) ([]*bloodreqmodel.DonorResponse, error)
	GetDonorResponsesByDonorID(ctx context.Context, donorID string) ([]*bloodreqmodel.DonorResponse, error)
	ExistsByID(ctx context.Context, id string) (bool, error)
	ExistsByRequestID(ctx context.Context, reqID string) (bool, error)
	ExistsByDonorID(ctx context.Context, donorID string) (bool, error)
	Count(ctx context.Context) (int, error)
}
