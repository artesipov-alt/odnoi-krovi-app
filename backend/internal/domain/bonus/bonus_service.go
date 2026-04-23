package bonus

import (
	"context"
	"time"

	bonusmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bonus/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/common"
)

// BonusService handles bonus aggregation logic.
type BonusService struct {
	repo Repository
}

// NewBonusService creates a new bonus service.
func NewBonusService(repo Repository) *BonusService {
	return &BonusService{repo: repo}
}

// GetAggregatedBonuses retrieves and aggregates bonuses: one per subcategory per partner, prioritized by expiration date.
func (s *BonusService) GetAggregatedBonuses(ctx context.Context, petType common.PetType, userID string) ([]*bonusmodel.Bonus, error) {
	// Get the last bonus for the user to check if locked
	lastBonus, err := s.repo.GetLastBonus(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Check if the last bonus was updated within the last 2 months, if so, return only lock bonus
	isLocked := lastBonus != nil && time.Since(lastBonus.UpdatedAt) < 2*30*24*time.Hour
	if isLocked {
		return []*bonusmodel.Bonus{bonusmodel.NewLockBonus()}, nil
	}

	bonuses, err := s.repo.GetAvailableBonuses(ctx, petType)
	if err != nil {
		return nil, err
	}

	// Group by partner_name -> subcategory -> list of bonuses
	grouped := make(map[string]map[string][]*bonusmodel.Bonus)
	for _, b := range bonuses {
		subcat := ""
		if b.Subcategory != nil {
			subcat = *b.Subcategory
		}
		if grouped[b.PartnerName] == nil {
			grouped[b.PartnerName] = make(map[string][]*bonusmodel.Bonus)
		}
		grouped[b.PartnerName][subcat] = append(grouped[b.PartnerName][subcat], b)
	}

	// For each partner and subcategory, select the bonus with the latest expiration
	var result []*bonusmodel.Bonus
	for _, subcats := range grouped {
		for _, bs := range subcats {
			if len(bs) == 0 {
				continue
			}
			result = append(result, bs[0])
		}
	}

	return result, nil
}

// AssignBonuses assigns available bonuses for a pet type to a user.
func (s *BonusService) AssignBonuses(ctx context.Context, userID string, petType common.PetType, donorResponseID string) error {
	// Get available bonuses
	bonuses, err := s.GetAggregatedBonuses(ctx, petType, userID)
	if err != nil {
		return err
	}

	// Filter out lock bonuses, as they cannot be assigned
	var assignableBonuses []*bonusmodel.Bonus
	for _, b := range bonuses {
		if b.Category != "lock" {
			assignableBonuses = append(assignableBonuses, b)
		}
	}

	// Collect IDs
	bonusIDs := make([]string, len(assignableBonuses))
	for i, b := range assignableBonuses {
		bonusIDs[i] = b.ID
	}

	// Assign them
	err = s.repo.AssignBonuses(ctx, bonusIDs, userID, donorResponseID)
	if err != nil {
		return err
	}

	// Set last donation if bonuses were assigned
	if len(bonusIDs) > 0 {
		return s.repo.SetLastDonation(ctx, userID, time.Now())
	}

	return nil
}

// UnassignBonuses unassigns bonuses from a user for a specific pet type.
func (s *BonusService) UnassignBonuses(ctx context.Context, userID string, petType common.PetType) error {
	err := s.repo.UnassignBonuses(ctx, userID, petType)
	if err != nil {
		return err
	}

	return nil
}

// ConfirmBonuses confirms bonuses for a user by setting stage to unused and adds priority search.
func (s *BonusService) ConfirmBonuses(ctx context.Context, userID string, petType common.PetType) error {
	err := s.repo.ConfirmBonuses(ctx, userID, petType)
	if err != nil {
		return err
	}
	return s.repo.AddPrioritySearch(ctx, userID)
}

// MarkBonusesAsUsed marks reserved bonuses for a user as used.
func (s *BonusService) MarkBonusesAsUsed(ctx context.Context, userID string, petType common.PetType) error {
	return s.repo.MarkBonusesAsUsed(ctx, userID, petType)
}
