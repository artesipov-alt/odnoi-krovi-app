package repositories

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/models"
)

// PetRepository определяет интерфейс для операций с данными питомцев
type PetRepository interface {
	// Create создает нового питомца в базе данных
	Create(ctx context.Context, pet *models.Pet) error

	// GetByID получает питомца по его ID с возможностью предзагрузки связей
	GetByID(ctx context.Context, id string, preloads ...string) (*models.Pet, error)

	// GetByUserID получает всех питомцев конкретного пользователя с возможностью предзагрузки связей
	GetByUserID(ctx context.Context, userID string, preloads ...string) ([]*models.Pet, error)

	// Update обновляет существующего питомца в базе данных
	Update(ctx context.Context, pet *models.Pet) error

	// Delete удаляет питомца по его ID
	Delete(ctx context.Context, id string) error

	// ExistsByID проверяет, существует ли питомец с заданным ID
	ExistsByID(ctx context.Context, id string) (bool, error)
}
