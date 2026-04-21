package pg

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	bonusmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bonus/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/common"
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
func (r *EntBonusRepository) CreateBatch(ctx context.Context, bonuses []*bonusmodel.Bonus) error {
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

// GetAvailableBonuses retrieves available (unassigned) bonuses filtered by petType and isActive.
func (r *EntBonusRepository) GetAvailableBonuses(ctx context.Context, petType common.PetType, isActive bool) ([]*bonusmodel.Bonus, error) {
	query := r.client.Bonus.Query().
		Where(entbonus.IsActive(isActive)).
		Where(entbonus.UserIDIsNil()).
		Where(entbonus.TargetIn(entbonus.Target(petType), entbonus.TargetAll)).
		Order(entbonus.ByExpiresAt(sql.OrderDesc()))

	bonuses, err := query.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to query bonuses: %w", err)
	}

	result := make([]*bonusmodel.Bonus, len(bonuses))
	for i, b := range bonuses {
		result[i] = &bonusmodel.Bonus{
			ID:          b.ID,
			UserID:      &b.UserID,
			PartnerName: b.PartnerName,
			Description: b.Description,
			Target:      b.Target.String(),
			Recipient:   b.Recipient.String(),
			Category:    b.Category.String(),
			Subcategory: &b.Subcategory,
			// Возвращаем без промокода.
			// PromoCode:    b.PromoCode,
			ExpiresAt:    b.ExpiresAt,
			PlatformName: b.PlatformName,
			PlatformURL:  &b.PlatformURL,
			IsActive:     b.IsActive,
		}
	}

	return result, nil
}
