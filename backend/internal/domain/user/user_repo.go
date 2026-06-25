package user

import (
	"context"

	authmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/auth/model"
	usermodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user/model"
)

// Repository определяет интерфейс для операций с данными пользователей
type Repository interface {
	// Write methods

	// создает нового пользователя
	CreateUser(ctx context.Context, inputuser *usermodel.User) (*usermodel.User, error)

	// создает или обновляет identity пользователя
	UpsertUserIdentity(ctx context.Context, userID string, input *authmodel.Identity, metadata *authmodel.Metadata) error

	// создает настройки донора для пользователя
	CreateDonorPreference(ctx context.Context, userID string, inputprefs *usermodel.DonorPreference) error

	// обновляет поля пользователя (атомарная операция, не включает phone)
	UpdateUserFields(ctx context.Context, id string, input *usermodel.User) error

	// обновляет номер телефона пользователя по его ID
	UpdatePhone(ctx context.Context, id string, phone string) error

	// переносит все identity от одного пользователя к другому
	TransferUserIdentity(ctx context.Context, fromUserID, toUserID string) error

	// полностью удаляет пользователя (обход soft delete)
	DeleteUserHard(ctx context.Context, id string) error

	// создает или обновляет настройки донора
	UpsertDonorPreference(ctx context.Context, userID string, prefs *usermodel.DonorPreference) error

	// удаляет настройки донора по ID пользователя
	DeleteDonorPreferenceByUserID(ctx context.Context, userID string) error

	// удаляет UTM-историю по ID пользователя
	DeleteUTMHistoryByUserID(ctx context.Context, userID string) error

	// переносит UTM-историю от одного пользователя к другому
	TransferUTMHistory(ctx context.Context, fromUserID, toUserID string) error

	// удаляет пользователя по его ID
	Delete(ctx context.Context, id string) error

	// сбрасывает email и номер телефона пользователя по ID
	ResetUser(ctx context.Context, id string) error

	// восстанавливает пользователя по его ID
	RestoreUser(ctx context.Context, id string) error

	// создает или обновляет UTM-метаданные пользователя
	UpsertUTM(ctx context.Context, userID string, metadata *authmodel.Metadata) error

	// добавляет новые пути к фотографиям пользователя
	AddPhotoURLs(ctx context.Context, id string, paths []string) error

	// Read methods

	// возвращает пользователя по ID
	GetByID(ctx context.Context, id string, opts UserPreloadOptions) (*usermodel.User, error)

	// возвращает identity по ID провайдера
	GetByProvider(ctx context.Context, providerID string, providerName authmodel.ProviderName) (*authmodel.Identity, error)

	// проверяет, существует ли пользователь с заданным providerID
	ExistsByProvider(ctx context.Context, providerID string, providerName authmodel.ProviderName) (bool, error)

	// проверяет, существует ли пользователь с заданным ID
	ExistsByID(ctx context.Context, id string) (bool, error)

	// получает всех удаленных пользователей
	GetDeletedUsers(ctx context.Context) ([]*usermodel.User, error)

	// возвращает ID пользователя по номеру телефона
	GetByPhone(ctx context.Context, phone string) (string, error)
}

type UserPreloadOptions struct {
	WithPets            bool
	WithDonorPreference bool
	WithIdentities      bool
}
