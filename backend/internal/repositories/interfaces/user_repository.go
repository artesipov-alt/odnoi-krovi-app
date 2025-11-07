package repositories

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/models"
)

// UserRepository определяет интерфейс для операций с данными пользователей
type UserRepository interface {
	// Create создает нового пользователя в базе данных
	Create(ctx context.Context, user *models.User) error

	// GetByID получает пользователя по его ID
	GetByID(ctx context.Context, id int) (*models.User, error)

	// GetByTelegramID получает пользователя по его Telegram ID
	GetByTelegramID(ctx context.Context, telegramID int64) (*models.User, error)

	// Update обновляет существующего пользователя в базе данных
	Update(ctx context.Context, user *models.User) error

	// Delete удаляет пользователя по его ID
	Delete(ctx context.Context, id int) error

	// ExistsByTelegramID проверяет, существует ли пользователь с заданным Telegram ID
	ExistsByTelegramID(ctx context.Context, telegramID int64) (bool, error)

	// ResetUser сбрасывает email и номер телефона пользователя по ID
	ResetUser(ctx context.Context, id int) error

	// RestoreUser восстанавливает пользователя по его ID
	RestoreUser(ctx context.Context, id int) error

	// GetDeletedUsers получает всех удаленных пользователей
	GetDeletedUsers(ctx context.Context) ([]*models.User, error)
}
