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

	// GetByUserID получает всех питомцев конкретного пользователя со связями по запросу
	GetByUserID(ctx context.Context, userID string, preloads ...string) ([]*ent.Pet, error)

	// Update обновляет существующего питомца в базе данных
	Update(ctx context.Context, pet *ent.Pet, health *ent.PetHealth, treatments *ent.PetTreatment, analyses []*ent.PetAnalysis, bonuses *ent.PetBonus) (*ent.Pet, error)

	// Delete удаляет питомца по его ID
	Delete(ctx context.Context, id string) error

	// ExistsByID проверяет, существует ли питомец с заданным ID
	ExistsByID(ctx context.Context, id string) (bool, error)
}
