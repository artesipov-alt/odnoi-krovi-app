package pg

import (
	"context"
	"fmt"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bonus"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
	entbonus "github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/bonus"
)

// EntBonusRepository implements bonus.Repository using ENT
type EntBonusRepository struct {
	client *ent.Client
}

// NewEntBonusRepository creates a new ENT bonus repository
func NewEntBonusRepository(client *ent.Client) *EntBonusRepository {
	return &EntBonusRepository{
		client: client,
	}
}

// CreateBatch creates multiple bonuses in a single transaction.
func (r *EntBonusRepository) CreateBatch(ctx context.Context, bonuses []*bonus.Bonus) error {
	if len(bonuses) == 0 {
		return nil
	}

	creates := make([]*ent.BonusCreate, len(bonuses))
	for i, b := range bonuses {
		create := r.client.Bonus.Create().
			SetPartnerName(b.PartnerName).
			SetDescription(b.Description).
			SetTarget(entbonus.Target(b.Target)).
			SetRecipient(entbonus.Recipient(b.Recipient)).
			SetCategory(entbonus.Category(b.Category)).
			SetPromoCode(b.PromoCode).
			SetExpiresAt(b.ExpiresAt).
			SetPlatformName(b.PlatformName).
			SetIsActive(b.IsActive)

		if b.UserID != nil {
			create.SetUserID(*b.UserID)
		}

		if b.PlatformURL != nil {
			create.SetPlatformURL(*b.PlatformURL)
		}

		creates[i] = create
	}

	_, err := r.client.Bonus.CreateBulk(creates...).Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to create bonuses batch: %w", err)
	}

	return nil
}

// ExistsByPromoCodes returns a list of promo codes that already exist in the database.
func (r *EntBonusRepository) ExistsByPromoCodes(ctx context.Context, codes []string) ([]string, error) {
	if len(codes) == 0 {
		return nil, nil
	}

	existingBonuses, err := r.client.Bonus.Query().
		Where(entbonus.PromoCodeIn(codes...)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query existing promo codes: %w", err)
	}

	result := make([]string, len(existingBonuses))
	for i, b := range existingBonuses {
		result[i] = b.PromoCode
	}

	return result, nil
}
