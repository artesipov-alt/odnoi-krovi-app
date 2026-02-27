package domain

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
)

// UserRepository определяет интерфейс для операций с данными пользователей
type UserRepository interface {
	// Create создает нового пользователя в базе данных
	Create(ctx context.Context, user *ent.CreateUserInput) (*ent.User, error)

	// GetByID возвращает пользователя по ID
	GetByID(ctx context.Context, id string, opts UserPreloadOptions) (*ent.User, error)

	// GetByTelegram возвращает пользователя по Telegram ID
	GetByTelegram(ctx context.Context, telegramID int64, opts UserPreloadOptions) (*ent.User, error)

	// Update обновляет существующего пользователя в базе данных
	Update(ctx context.Context, id string, input *ent.UpdateUserInput) error

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

// LocationRepository определяет интерфейс для операций с данными локаций
type LocationRepository interface {
	// GetByID получает локацию по её ID
	GetByID(ctx context.Context, id string) (*ent.Location, error)

	// GetAll получает все локации из базы данных
	GetAll(ctx context.Context) ([]*ent.Location, error)

	// Exists проверяет, существует ли локация с заданным ID
	Exists(ctx context.Context, id string) (bool, error)
}

type UserPreloadOptions struct {
	WithPets bool
}
