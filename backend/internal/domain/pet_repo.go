package domain

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/pet/model"
)

// PetRepository определяет интерфейс для операций с данными питомцев
type PetRepository interface {
	// Create создает нового питомца в базе данных
	Create(ctx context.Context, petDomain *model.Pet) (*model.Pet, error)

	// GetPetQuery возвращает query для eager loading
	GetPet(ctx context.Context, id string, opts PetPreloadOptions) (*model.Pet, error)

	// GetPetsQueryByUser возвращает query для eager loading питомцев пользователя
	GetPetsByUser(ctx context.Context, userID string, opts PetPreloadOptions) ([]*model.Pet, error)

	// Update обновляет питомца и его связанные сущности в одной транзакции
	// Все параметры (кроме id и pet) могут быть nil - тогда соответствующие данные не обновляются
	Update(ctx context.Context, id string, petDomain *model.Pet) (*model.Pet, error)

	// Delete удаляет питомца по его ID
	Delete(ctx context.Context, id string) error

	// ExistsByID проверяет, существует ли питомец с заданным ID
	ExistsByID(ctx context.Context, id string) (bool, error)

	CountSuitableDonors(ctx context.Context, bloodGroups []string) (int, error)

	// // UpdateStatus обновляет статус питомца по его ID
	// UpdateStatus(ctx context.Context, id string, status string) error

	// AddPhotoURLs добавляет новые пути к фотографиям питомца
	AddPhotoURLs(ctx context.Context, id string, paths []string) error
}

// PetPreloadOptions определяет опции для preload связанных данных питомца
type PetPreloadOptions struct {
	WithHealth     bool
	WithTreatments bool
	WithAnalyses   bool
	WithBonuses    bool
	WithAll        bool
}
