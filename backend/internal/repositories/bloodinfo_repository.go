package repositories

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/models"
)

// BloodInfoRepository определяет интерфейс для работы с типами крови
type BloodInfoRepository interface {
	AllComponents(ctx context.Context) ([]models.BloodComponent, error)
	ComponentByID(ctx context.Context, id int) (*models.BloodComponent, error)
	BloodGroupsByPetType(ctx context.Context, petType models.PetType) ([]*models.BloodGroup, error)
}
