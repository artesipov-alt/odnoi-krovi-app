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
func (s *BonusService) GetAggregatedBonuses(ctx context.Context, petType common.PetType, lastDonation *time.Time) ([]*bonusmodel.Bonus, error) {
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

	// Check if user has donated within the last 2 months, if so, add a lock bonus
	if lastDonation != nil && time.Since(*lastDonation) < 2*30*24*time.Hour {
		lockBonus := &bonusmodel.Bonus{
			PartnerName: "Портал",
			Description: "Пользователь уже получал свои бонусы в течение двух месяцев.",
			Category:    "lock",
			Target:      "all",
			Recipient:   "all",
			Stage:       "unused",
		}
		result = append(result, lockBonus)
	}

	return result, nil
}

// AssignBonuses assigns available bonuses for a pet type to a user.
func (s *BonusService) AssignBonuses(ctx context.Context, userID string, petType common.PetType, lastDonation *time.Time) error {
	// Get available bonuses
	bonuses, err := s.GetAggregatedBonuses(ctx, petType, lastDonation)
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
	return s.repo.AssignBonuses(ctx, bonusIDs, userID)
}

// UnassignBonuses unassigns bonuses from a user for a specific pet type.
func (s *BonusService) UnassignBonuses(ctx context.Context, userID string, petType common.PetType) error {
	return s.repo.UnassignBonuses(ctx, userID, petType)
}

// ConfirmBonuses confirms bonuses for a user by setting stage to unused and sets the last donation date.
func (s *BonusService) ConfirmBonuses(ctx context.Context, userID string, petType common.PetType, donationDate time.Time) error {
	return s.repo.ConfirmBonuses(ctx, userID, petType, donationDate)
}

// MarkBonusesAsUsed marks reserved bonuses for a user as used.
func (s *BonusService) MarkBonusesAsUsed(ctx context.Context, userID string, petType common.PetType) error {
	return s.repo.MarkBonusesAsUsed(ctx, userID, petType)
}
