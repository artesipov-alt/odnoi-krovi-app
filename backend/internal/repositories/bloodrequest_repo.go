package repositories

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/ent"
)

// BloodRequestRepository определяет интерфейс для работы с данными заявок на поиск крови питомцев
type BloodRequestRepository interface {
	// Create создает новую заявку на поиск крови
	Create(ctx context.Context, request *ent.BloodSearchRequest) (*ent.BloodSearchRequest, error)

	// CreateWithTx создает новую заявку на поиск крови в рамках транзакции
	CreateWithTx(ctx context.Context, tx *ent.Tx, request *ent.BloodSearchRequest) (*ent.BloodSearchRequest, error)

	// GetByID возвращает заявку по её идентификатору
	GetByID(ctx context.Context, id string) (*ent.BloodSearchRequest, error)

	// GetByPetID возвращает заявку по идентификатору питомца
	GetByPetID(ctx context.Context, petID string) (*ent.BloodSearchRequest, error)

	// Update обновляет информацию о заявке
	Update(ctx context.Context, request *ent.BloodSearchRequest) (*ent.BloodSearchRequest, error)

	// UpdateStatus обновляет статус заявки
	UpdateStatus(ctx context.Context, id string, status string) error

	// UpdateStatusWithTx обновляет статус заявки в рамках транзакции
	UpdateStatusWithTx(ctx context.Context, tx *ent.Tx, id string, status string) error

	// Delete удаляет заявку из хранилища
	Delete(ctx context.Context, id string) error

	// DeleteWithTx удаляет заявку из хранилища в рамках транзакции
	DeleteWithTx(ctx context.Context, tx *ent.Tx, id string) error

	// List возвращает список заявок с фильтрацией и пагинацией
	List(ctx context.Context, limit, offset int, filters map[string]any) ([]*ent.BloodSearchRequest, error)

	// ExistsByPetID проверяет существование активной заявки для питомца
	ExistsByPetID(ctx context.Context, petID string) (bool, error)

	// ExistsByID проверяет существование заявки по её идентификатору
	ExistsByID(ctx context.Context, id string) (bool, error)

	// Count возвращает общее количество заявок в хранилище
	Count(ctx context.Context) (int, error)

	// AddPhotoURLs добавляет новые пути к фотографиям заявки
	AddPhotoURLs(ctx context.Context, id string, paths []string) error
}
