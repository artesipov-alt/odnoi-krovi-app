package reference

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
)

// BloodInfoRepository определяет интерфейс для работы с типами крови
type BloodInfoRepository interface {
	AllComponents(ctx context.Context) ([]*ent.BloodComponent, error)
	ComponentByID(ctx context.Context, id string) (*ent.BloodComponent, error)
}
