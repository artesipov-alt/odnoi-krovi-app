package donor

import (
	"context"
	"time"

	bloodreqmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
)

type Repository interface {
	// Write methods

	// создает отклик донора
	CreateDonorResponse(ctx context.Context, resp *donormodel.DonorResponse) (*donormodel.DonorResponse, error)

	// обновляет статус отклика донора
	UpdateDonorResponseStatus(ctx context.Context, id string, status donormodel.DonorResponseStatus) error

	// удаляет отклик донора
	DeleteDonorResponse(ctx context.Context, id string) error

	// обновляет отклик донора
	Update(ctx context.Context, resp *donormodel.DonorResponse) error

	// Read methods

	// возвращает отклик донора по ID
	GetDonorResponseByID(ctx context.Context, id string) (*donormodel.DonorResponse, error)

	// возвращает заявку реципиента с подходящими донорами
	GetRecipient(ctx context.Context, id string) (*bloodreqmodel.BloodRequestWithMatchingDonors, error)

	// возвращает все отклики по ID заявки
	GetDonorResponsesByRequestID(ctx context.Context, reqID string) ([]*donormodel.DonorResponse, error)

	// возвращает все отклики по ID донора
	GetDonorResponsesByDonorID(ctx context.Context, donorID string) ([]*donormodel.DonorResponse, error)

	// проверяет, существует ли отклик с заданным ID
	ExistsByID(ctx context.Context, id string) (bool, error)

	// проверяет, существуют ли отклики по ID заявки
	ExistsByRequestID(ctx context.Context, reqID string) (bool, error)

	// проверяет, существуют ли отклики по ID донора
	ExistsByDonorID(ctx context.Context, donorID string) (bool, error)

	// возвращает количество откликов
	Count(ctx context.Context) (int, error)

	// возвращает отклик по ID питомца
	GetByPetID(ctx context.Context, petID string) (*donormodel.DonorResponse, error)

	// возвращает отклики по списку ID питомцев
	GetByPetIDs(ctx context.Context, petIDs []string, ignoreSoftDelete bool) (map[string][]*donormodel.DonorResponse, error)

	// GetLatestByPetIDs возвращает мапу "pet ID → последний по created_at отклик донора".
	// Не фильтрует по статусу — решение об актуальности отклика принимает
	// DonorResponse.IsActiveForDonation() на стороне вызывающего кода.
	GetLatestByPetIDs(ctx context.Context, petIDs []string) (map[string]*donormodel.DonorResponse, error)

	// CountFullyCompletedByOwnerID возвращает количество полностью завершённых
	// (подтверждённых реципиентом) донаций по всем питомцам владельца, включая
	// мягко удалённых — статистика не должна теряться при удалении питомца.
	// Критерий соответствует DonorResponse.IsFullyCompleted().
	CountFullyCompletedByOwnerID(ctx context.Context, ownerID string) (int, error)

	// возвращает неподтвержденные отклики старше cutoffTime
	FindNotConfirmed(ctx context.Context, cutoffTime time.Time) ([]*donormodel.DonorResponse, error)

	// FindAcceptedForAutoConfirm возвращает accepted-отклики старше cutoffTime,
	// для которых донация так и не была отмечена. Используется автоподтверждением
	// для кейса «донор принят, но реципиент пропал и не подтвердил донацию».
	FindAcceptedForAutoConfirm(ctx context.Context, cutoffTime time.Time) ([]*donormodel.DonorResponse, error)
}
