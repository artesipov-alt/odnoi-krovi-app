package repositories

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/ent"
)

// UserRepository определяет интерфейс для операций с данными пользователей
type UserRepository interface {
	// Create создает нового пользователя в базе данных
	Create(ctx context.Context, user *ent.User) (*ent.User, error)

	// GetByID получает пользователя по его ID
	GetByID(ctx context.Context, id string, preloads ...string) (*ent.User, error)

	// GetByTelegramID получает пользователя по его Telegram ID
	GetByTelegramID(ctx context.Context, telegramID int64) (*ent.User, error)

	// Update обновляет существующего пользователя в базе данных
	Update(ctx context.Context, user *ent.User) (*ent.User, error)

	// Delete удаляет пользователя по его ID
	Delete(ctx context.Context, id string) error

	// ExistsByTelegramID проверяет, существует ли пользователь с заданным Telegram ID
	ExistsByTelegramID(ctx context.Context, telegramID int64) (bool, error)

	// ExistsByID проверяет, существует ли пользователь с заданным ID
	ExistsByID(ctx context.Context, id string) (bool, error)

	// ResetUser сбрасывает email и номер телефона пользователя по ID
	ResetUser(ctx context.Context, id string) error

	// RestoreUser восстанавливает пользователя по его ID
	RestoreUser(ctx context.Context, id string) error

	// GetDeletedUsers получает всех удаленных пользователей
	GetDeletedUsers(ctx context.Context) ([]*ent.User, error)

	// AddPhotoURLs добавляет новые пути к фотографиям пользователя
	AddPhotoURLs(ctx context.Context, id string, paths []string) error
}
