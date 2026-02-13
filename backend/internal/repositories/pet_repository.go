package repositories

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/ent"
)

// PetRepository определяет интерфейс для операций с данными питомцев
type PetRepository interface {
	// Create создает нового питомца в базе данных
	Create(ctx context.Context, pet *ent.Pet, health *ent.PetHealth, treatments *ent.PetTreatment, analyses []*ent.PetAnalysis, bonuses *ent.PetBonus) (*ent.Pet, error)

	// GetByID получает питомца по его ID со связями по запросу
	GetByID(ctx context.Context, id string, preloads ...string) (*ent.Pet, error)

	// GetPetQuery возвращает query для eager loading
	GetPetQuery(ctx context.Context, id string) *ent.PetQuery

	// GetByUserID получает всех питомцев конкретного пользователя со связями по запросу
	GetByUserID(ctx context.Context, userID string, preloads ...string) ([]*ent.Pet, error)

	// GetPetsQueryByUser возвращает query для eager loading питомцев пользователя
	GetPetsQueryByUser(ctx context.Context, userID string) *ent.PetQuery

	// Update обновляет существующего питомца в базе данных
	Update(ctx context.Context, pet *ent.Pet, health *ent.PetHealth, treatments *ent.PetTreatment, analyses []*ent.PetAnalysis, bonuses *ent.PetBonus) (*ent.Pet, error)

	// Delete удаляет питомца по его ID
	Delete(ctx context.Context, id string) error

	// ExistsByID проверяет, существует ли питомец с заданным ID
	ExistsByID(ctx context.Context, id string) (bool, error)

	// UpdateStatus обновляет статус питомца по его ID
	UpdateStatus(ctx context.Context, id string, status string) error

	// UpdateStatusWithTx обновляет статус питомца по его ID в рамках транзакции
	UpdateStatusWithTx(ctx context.Context, tx *ent.Tx, id string, status string) error

	// AddPhotoURLs добавляет новые пути к фотографиям питомца
	AddPhotoURLs(ctx context.Context, id string, paths []string) error
}
