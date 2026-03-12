package user

import (
	"context"

	authmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/auth/model"
	usermodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user/model"
)

// Repository определяет интерфейс для операций с данными пользователей
type Repository interface {
	CreateUser(ctx context.Context, inputuser *usermodel.User) (*usermodel.User, error)

	UpsertUserIdentity(ctx context.Context, userID string, input *authmodel.Identity, metadata *authmodel.Metadata) error

	// CreateDonorPreference создает настройки донора для пользователя
	CreateDonorPreference(ctx context.Context, userID string, inputprefs *usermodel.DonorPreference) error

	// GetByID возвращает пользователя по ID
	GetByID(ctx context.Context, id string, opts UserPreloadOptions) (*usermodel.User, error)

	// GetByProviderID возвращает идентификатор пользователя по ID провайдера
	GetByProvider(ctx context.Context, providerID string, providerName authmodel.ProviderName) (*authmodel.Identity, error)

	// UpdateUserFields обновляет поля пользователя (атомарная операция)
	UpdateUserFields(ctx context.Context, id string, input *usermodel.User) error

	// TransferUserIdentity переносит все identity от одного пользователя к другому
	TransferUserIdentity(ctx context.Context, fromUserID, toUserID string) error

	// DeleteUserHard полностью удаляет пользователя (обход soft delete)
	DeleteUserHard(ctx context.Context, id string) error

	// UpsertDonorPreference создает или обновляет настройки донора
	UpsertDonorPreference(ctx context.Context, userID string, prefs *usermodel.DonorPreference) error

	// DeleteDonorPreferenceByUserID удаляет настройки донора по ID пользователя
	DeleteDonorPreferenceByUserID(ctx context.Context, userID string) error

	// DeleteUTMHistoryByUserID удаляет UTM-историю по ID пользователя
	DeleteUTMHistoryByUserID(ctx context.Context, userID string) error

	// TransferUTMHistory переносит UTM-историю от одного пользователя к другому
	TransferUTMHistory(ctx context.Context, fromUserID, toUserID string) error

	// Delete удаляет пользователя по его ID
	Delete(ctx context.Context, id string) error

	// ExistsByTelegramID проверяет, существует ли пользователь с заданным Telegram ID
	ExistsByProvider(ctx context.Context, providerID string, providerName authmodel.ProviderName) (bool, error)

	// ExistsByID проверяет, существует ли пользователь с заданным ID
	ExistsByID(ctx context.Context, id string) (bool, error)

	// ResetUser сбрасывает email и номер телефона пользователя по ID
	ResetUser(ctx context.Context, id string) error

	// RestoreUser восстанавливает пользователя по его ID
	RestoreUser(ctx context.Context, id string) error

	// GetDeletedUsers получает всех удаленных пользователей
	GetDeletedUsers(ctx context.Context) ([]*usermodel.User, error)

	// GetByPhone возвращает ID пользователя по номеру телефона
	GetByPhone(ctx context.Context, phone string) (string, error)

	UpsertUTM(ctx context.Context, userID string, metadata *authmodel.Metadata) error

	// AddPhotoURLs добавляет новые пути к фотографиям пользователя
	AddPhotoURLs(ctx context.Context, id string, paths []string) error
}

type UserPreloadOptions struct {
	WithPets            bool
	WithDonorPreference bool
}
