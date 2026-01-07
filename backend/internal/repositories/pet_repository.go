package repositories

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/ent"
)

// PetRepository определяет интерфейс для операций с данными питомцев
type PetRepository interface {
	// Create создает нового питомца в базе данных
	Create(ctx context.Context, pet *ent.Pet) (*ent.Pet, error)

	// GetByID получает питомца по его ID со всеми связями
	GetByID(ctx context.Context, id string) (*ent.Pet, error)

	// GetByUserID получает всех питомцев конкретного пользователя со всеми связями
	GetByUserID(ctx context.Context, userID string) ([]*ent.Pet, error)

	// Update обновляет существующего питомца в базе данных
	Update(ctx context.Context, pet *ent.Pet) (*ent.Pet, error)

	// Delete удаляет питомца по его ID
	Delete(ctx context.Context, id string) error

	// ExistsByID проверяет, существует ли питомец с заданным ID
	ExistsByID(ctx context.Context, id string) (bool, error)
}
