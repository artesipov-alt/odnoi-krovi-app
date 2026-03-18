package pet

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
)

// PetReadRepository определяет операции чтения для питомцев
type PetReadRepository interface {
	// GetByID получает питомца по ID с опциями загрузки связанных данных
	GetByID(ctx context.Context, id string, opts PetPreloadOptions) (*model.Pet, error)

	// GetByUserID получает всех питомцев пользователя
	GetByUserID(ctx context.Context, userID string, opts PetPreloadOptions) ([]*model.Pet, error)

	// Exists проверяет существование питомца по ID
	Exists(ctx context.Context, id string) (bool, error)

	// ExistsByID алиас для Exists (для обратной совместимости)
	// Deprecated: используйте Exists
	ExistsByID(ctx context.Context, id string) (bool, error)
}

// PetWriteRepository определяет операции записи для питомцев
type PetWriteRepository interface {
	// Create создает нового питомца
	Create(ctx context.Context, pet *model.Pet) (*model.Pet, error)

	// Update обновляет существующего питомца
	Update(ctx context.Context, id string, pet *model.Pet) (*model.Pet, error)

	// Delete удаляет питомца (soft delete)
	Delete(ctx context.Context, id string) error
}

// PetStatsRepository определяет операции для статистики и поиска доноров
type PetStatsRepository interface {
	// CountSuitableDonors считает количество подходящих доноров для групп крови
	CountSuitableDonors(ctx context.Context, bloodGroups []string) (int, error)
}

// PetPhotoRepository определяет операции для работы с фото питомцев
type PetPhotoRepository interface {
	// AddPhotoURLs добавляет URL фотографий к питомцу
	AddPhotoURLs(ctx context.Context, id string, paths []string) error
}

// PetPreloadOptions определяет опции для загрузки связанных данных
type PetPreloadOptions struct {
	WithHealth           bool
	WithTreatments       bool
	WithAnalyses         bool
	WithBonuses          bool
	WithBloodReq         bool
	WithDonorApplication bool
	WithAll              bool
}

// Repository объединяет все интерфейсы для обратной совместимости
// Deprecated: используйте специализированные интерфейсы
type Repository interface {
	PetReadRepository
	PetWriteRepository
	PetStatsRepository
	PetPhotoRepository
}
