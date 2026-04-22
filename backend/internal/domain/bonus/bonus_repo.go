package bonus

import (
	"context"
	"time"

	bonusmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bonus/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/common"
)

// Repository defines the interface for bonus data access.
type Repository interface {
	// CreateBatch creates multiple bonuses in a single transaction.
	CreateBatch(ctx context.Context, bonuses []*bonusmodel.Bonus) error
	// ExistsByPromoCodes returns a list of promo codes that already exist in the database.
	ExistsByPromoCodes(ctx context.Context, codes []string) ([]string, error)

	// GetAvailableBonuses retrieves available (unassigned) bonuses based on filters (petType: common.PetType; stage: unused).
	GetAvailableBonuses(ctx context.Context, petType common.PetType) ([]*bonusmodel.Bonus, error)

	// GetLastDonation gets the last donation date for a user
	GetLastDonation(ctx context.Context, userID string) (*time.Time, error)

	// UnassignBonuses unassigns bonuses from a user for a specific pet type by setting UserID to nil and stage to unused.
	UnassignBonuses(ctx context.Context, userID string, petType common.PetType) error

	// AssignBonuses assigns bonuses to a user by updating their UserID and DonorResponseID, and setting stage to reserved.
	AssignBonuses(ctx context.Context, bonusIDs []string, userID string, donorResponseID string) error

	// ConfirmBonuses confirms bonuses for a user by setting stage to unused without clearing UserID.
	ConfirmBonuses(ctx context.Context, userID string, petType common.PetType) error

	// MarkBonusesAsUsed marks reserved bonuses for a user as used by setting stage to used.
	MarkBonusesAsUsed(ctx context.Context, userID string, petType common.PetType) error

	// SetLastDonation sets the last donation date for a user
	SetLastDonation(ctx context.Context, userID string, donationDate time.Time) error

	// AddPrioritySearch increments the priority search count for a user by 1
	AddPrioritySearch(ctx context.Context, id string) error

	// SubtractPrioritySearch decrements the priority search count for a user by 1
	SubtractPrioritySearch(ctx context.Context, id string) error
}
