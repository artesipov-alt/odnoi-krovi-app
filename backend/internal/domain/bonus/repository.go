package bonus

import "context"

// Repository defines the interface for bonus data access.
type Repository interface {
	// CreateBatch creates multiple bonuses in a single transaction.
	CreateBatch(ctx context.Context, bonuses []*Bonus) error
	// ExistsByPromoCodes returns a list of promo codes that already exist in the database.
	ExistsByPromoCodes(ctx context.Context, codes []string) ([]string, error)
}
