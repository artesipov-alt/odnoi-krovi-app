package bonus

import (
	"context"

	bonusmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bonus/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/common"
)

// Repository defines the interface for bonus data access.
type Repository interface {
	// CreateBatch creates multiple bonuses in a single transaction.
	CreateBatch(ctx context.Context, bonuses []*bonusmodel.Bonus) error
	// ExistsByPromoCodes returns a list of promo codes that already exist in the database.
	ExistsByPromoCodes(ctx context.Context, codes []string) ([]string, error)

	// GetAvailableBonuses retrieves available (unassigned) bonuses based on filters (petType: common.PetType; isActive: true/false).
	GetAvailableBonuses(ctx context.Context, petType common.PetType, isActive bool) ([]*bonusmodel.Bonus, error)
}
