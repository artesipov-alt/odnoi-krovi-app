package repositories

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/ent"
)

// PetRepository определяет интерфейс для операций с данными питомцев
type PetRepository interface {
	// Create создает нового питомца в базе данных
	Create(ctx context.Context, input *ent.CreatePetInput, healthInput *ent.CreatePetHealthInput, treatmentsInput *ent.CreatePetTreatmentInput, analysesInput []*ent.CreatePetAnalysisInput, bonusesInput *ent.CreatePetBonusInput) (*ent.Pet, error)

	// GetByID получает питомца по ID со связями
	GetByID(ctx context.Context, id string, preloads ...string) (*ent.Pet, error)

	// GetPetQuery возвращает query для eager loading
	GetPetQuery(ctx context.Context, id string) *ent.PetQuery

	// GetByUserID получает всех питомцев конкретного пользователя
	GetByUserID(ctx context.Context, userID string, preloads ...string) ([]*ent.Pet, error)

	// GetPetsQueryByUser возвращает query для eager loading питомцев пользователя
	GetPetsQueryByUser(ctx context.Context, userID string) *ent.PetQuery

	// Update обновляет питомца и его связанные сущности в одной транзакции
	// Все параметры (кроме id и pet) могут быть nil - тогда соответствующие данные не обновляются
	Update(ctx context.Context, id string, petInput *ent.UpdatePetInput, healthInput *ent.UpdatePetHealthInput, treatmentsInput *ent.UpdatePetTreatmentInput, analysesInput []*ent.UpdatePetAnalysisInput, bonusesInput *ent.UpdatePetBonusInput) (*ent.Pet, error)

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
