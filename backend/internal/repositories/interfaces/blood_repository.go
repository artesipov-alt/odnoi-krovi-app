package repositories

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/models"
)

// BloodRepository определяет интерфейс для работы с типами крови
type BloodRepository interface {
	GetAllComponents(ctx context.Context) ([]models.BloodComponent, error)
	GetComponentByID(ctx context.Context, id int) (*models.BloodComponent, error)
	GetBloodGroupsByPetType(ctx context.Context, petType models.PetType) ([]*models.BloodGroup, error)
}
