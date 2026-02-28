package pet

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/pet/model"
)

// Repository определяет интерфейс для операций с данными питомцев
type Repository interface {
	Create(ctx context.Context, petDomain *model.Pet) (*model.Pet, error)
	GetPet(ctx context.Context, id string, opts PetPreloadOptions) (*model.Pet, error)
	GetPetsByUser(ctx context.Context, userID string, opts PetPreloadOptions) ([]*model.Pet, error)
	Update(ctx context.Context, id string, petDomain *model.Pet) (*model.Pet, error)
	Delete(ctx context.Context, id string) error
	ExistsByID(ctx context.Context, id string) (bool, error)
	CountSuitableDonors(ctx context.Context, bloodGroups []string) (int, error)
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
