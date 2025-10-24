package repositories

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/models"
)

// PetRepository определяет интерфейс для операций с данными питомцев
type PetRepository interface {
	// Create создает нового питомца в базе данных
	Create(ctx context.Context, pet *models.Pet) error

	// GetByID получает питомца по его ID
	GetByID(ctx context.Context, id int) (*models.Pet, error)

	// GetByUserID получает всех питомцев конкретного пользователя
	GetByUserID(ctx context.Context, userID int) ([]*models.Pet, error)

	// Update обновляет существующего питомца в базе данных
	Update(ctx context.Context, pet *models.Pet) error

	// Delete удаляет питомца по его ID
	Delete(ctx context.Context, id int) error

	// ExistsByID проверяет, существует ли питомец с заданным ID
	ExistsByID(ctx context.Context, id int) (bool, error)
}
