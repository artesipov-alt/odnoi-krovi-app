package repositories

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/models"
)

// BloodTypeRepository определяет интерфейс для работы с типами крови
type BloodTypeRepository interface {
	GetAll(ctx context.Context) ([]models.BloodType, error)
	GetByID(ctx context.Context, id int) (*models.BloodType, error)
	Create(ctx context.Context, bloodType *models.BloodType) error
	Update(ctx context.Context, bloodType *models.BloodType) error
	Delete(ctx context.Context, id int) error
}
