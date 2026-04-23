package pg

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	bonusmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bonus/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/common"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
	entbonus "github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/bonus"
	entuser "github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/user"
)

// EntBonusRepository реализует bonus.Repository с использованием ENT
type EntBonusRepository struct {
	db *ent.Client
}

// client возвращает ent.Client из контекста, если активна транзакция,
// иначе возвращает клиента по умолчанию
func (r *EntBonusRepository) client(ctx context.Context) *ent.Client {
	if tx := ent.TxFromContext(ctx); tx != nil {
		return tx.Client()
	}
	return r.db
}

// NewEntBonusRepository создает новый репозиторий бонусов ENT
func NewEntBonusRepository(client *ent.Client) *EntBonusRepository {
	return &EntBonusRepository{
		db: client,
	}
}

// CreateBatch создает несколько бонусов в одной транзакции.
func (r *EntBonusRepository) CreateBatch(ctx context.Context, bonuses []*bonusmodel.Bonus) error {
	if len(bonuses) == 0 {
		return nil
	}

	creates := make([]*ent.BonusCreate, len(bonuses))
	for i, b := range bonuses {
		create := r.client(ctx).Bonus.Create().
			SetPartnerName(b.PartnerName).
			SetDescription(b.Description).
			SetTarget(entbonus.Target(b.Target)).
			SetRecipient(entbonus.Recipient(b.Recipient)).
			SetCategory(entbonus.Category(b.Category)).
			SetPromoCode(b.PromoCode).
			SetExpiresAt(b.ExpiresAt).
			SetPlatformName(b.PlatformName)

		if b.UserID != nil {
			create.SetUserID(*b.UserID)
		}

		if b.PlatformURL != nil {
			create.SetPlatformURL(*b.PlatformURL)
		}

		creates[i] = create
	}

	_, err := r.client(ctx).Bonus.CreateBulk(creates...).Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to create bonuses batch: %w", err)
	}

	return nil
}

// ExistsByPromoCodes возвращает список промокодов, которые уже существуют в базе данных.
func (r *EntBonusRepository) ExistsByPromoCodes(ctx context.Context, codes []string) ([]string, error) {
	if len(codes) == 0 {
		return nil, nil
	}

	existingBonuses, err := r.client(ctx).Bonus.Query().
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

// AssignBonuses присваивает бонусы пользователю, обновляя их UserID и DonorResponseID, и устанавливая статус на reserved.
func (r *EntBonusRepository) AssignBonuses(ctx context.Context, bonusIDs []string, userID string, donorResponseID string) error {
	if len(bonusIDs) == 0 {
		return nil
	}

	_, err := r.client(ctx).Bonus.Update().
		Where(entbonus.IDIn(bonusIDs...)).
		SetUserID(userID).
		SetDonorResponseID(donorResponseID).
		SetStage(entbonus.StageReserved).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to assign bonuses: %w", err)
	}

	return nil
}

// UnassignBonuses отменяет присвоение зарезервированных бонусов от пользователя для определенного типа питомца, устанавливая UserID и DonorResponseID в nil и статус на unused.
func (r *EntBonusRepository) UnassignReservedBonuses(ctx context.Context, userID string, petType common.PetType) error {
	_, err := r.client(ctx).Bonus.Update().
		Where(entbonus.UserID(userID)).
		Where(entbonus.TargetIn(entbonus.Target(petType), entbonus.TargetAll)).
		Where(entbonus.StageEQ(entbonus.StageReserved)).
		ClearUserID().
		ClearDonorResponseID().
		SetStage(entbonus.StageUnused).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to unassign bonuses: %w", err)
	}

	return nil
}

// ConfirmBonuses подтверждает бонусы для пользователя, устанавливая статус на unused без очистки UserID.
func (r *EntBonusRepository) ConfirmBonuses(ctx context.Context, userID string, petType common.PetType) error {
	_, err := r.client(ctx).Bonus.Update().
		Where(entbonus.UserID(userID)).
		Where(entbonus.TargetIn(entbonus.Target(petType), entbonus.TargetAll)).
		Where(entbonus.StageEQ(entbonus.StageReserved)).
		SetStage(entbonus.StageUnused).
		SetAssignedAt(time.Now()).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to unreserve bonuses: %w", err)
	}

	return nil
}

// MarkBonusesAsUsed помечает зарезервированные бонусы для пользователя как использованные, устанавливая статус на used.
func (r *EntBonusRepository) MarkBonusesAsUsed(ctx context.Context, userID string, petType common.PetType) error {
	_, err := r.client(ctx).Bonus.Update().
		Where(entbonus.UserID(userID)).
		Where(entbonus.TargetIn(entbonus.Target(petType), entbonus.TargetAll)).
		Where(entbonus.StageEQ(entbonus.StageReserved)).
		SetStage(entbonus.StageUsed).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("failed to mark bonuses as used: %w", err)
	}

	return nil
}

// GetAvailableBonuses получает доступные (неприсвоенные) бонусы, отфильтрованные по petType.
func (r *EntBonusRepository) GetAvailableBonuses(ctx context.Context, petType common.PetType) ([]*bonusmodel.Bonus, error) {
	query := r.client(ctx).Bonus.Query().
		Where(entbonus.StageEQ(entbonus.StageUnused)).
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
			ExpiresAt:    b.ExpiresAt,
			PlatformName: b.PlatformName,
			PlatformURL:  &b.PlatformURL,
			Stage:        string(b.Stage),
			CreatedAt:    b.CreatedAt,
			UpdatedAt:    b.UpdatedAt,
			AssignedAt:   b.AssignedAt,
			DeletedAt:    b.DeletedAt,
		}
	}

	return result, nil
}

// GetLastBonus получает самый свежий бонус для пользователя по AssignedAt
//
// Примечание: AssignedAt устанавливается при подтверждении бонуса (ConfirmBonuses),
// что отражает время окончательного получения бонуса пользователем.
// Это обеспечивает корректную логику блокировки: если последний подтвержденный бонус
// был менее 2 месяцев назад, пользователь блокируется для новых бонусов.
func (r *EntBonusRepository) GetLastBonus(ctx context.Context, userID string) (*bonusmodel.Bonus, error) {
	bonus, err := r.client(ctx).Bonus.Query().
		Where(entbonus.UserID(userID)).
		Where(entbonus.AssignedAtNotNil()).
		Order(entbonus.ByAssignedAt(sql.OrderDesc())).
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get last bonus: %w", err)
	}

	return &bonusmodel.Bonus{
		ID:           bonus.ID,
		UserID:       &bonus.UserID,
		PartnerName:  bonus.PartnerName,
		Description:  bonus.Description,
		Target:       bonus.Target.String(),
		Recipient:    bonus.Recipient.String(),
		Category:     bonus.Category.String(),
		Subcategory:  &bonus.Subcategory,
		PromoCode:    bonus.PromoCode,
		ExpiresAt:    bonus.ExpiresAt,
		PlatformName: bonus.PlatformName,
		PlatformURL:  &bonus.PlatformURL,
		Stage:        string(bonus.Stage),
		CreatedAt:    bonus.CreatedAt,
		UpdatedAt:    bonus.UpdatedAt,
		AssignedAt:   bonus.AssignedAt,
		DeletedAt:    bonus.DeletedAt,
	}, nil
}

// AddPrioritySearch увеличивает счетчик приоритетного поиска для пользователя на 1
func (r *EntBonusRepository) AddPrioritySearch(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("invalid user ID")
	}

	err := r.client(ctx).User.Update().
		Where(entuser.ID(id)).
		AddPrioritySearchCount(1).
		Exec(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return fmt.Errorf("user with id %s not found", id)
		}
		return fmt.Errorf("failed to add priority search: %w", err)
	}

	return nil
}

// SubtractPrioritySearch уменьшает счетчик приоритетного поиска для пользователя на 1
func (r *EntBonusRepository) SubtractPrioritySearch(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("invalid user ID")
	}

	err := r.client(ctx).User.Update().
		Where(entuser.ID(id)).
		AddPrioritySearchCount(-1).
		Exec(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return fmt.Errorf("user with id %s not found", id)
		}
		return fmt.Errorf("failed to subtract priority search: %w", err)
	}

	return nil
}

// GetBonusesByDonorResponseID получает бонусы, связанные с конкретным откликом донора.
func (r *EntBonusRepository) GetBonusesByDonorResponseID(ctx context.Context, donorResponseID string) ([]*bonusmodel.Bonus, error) {
	bonuses, err := r.client(ctx).Bonus.Query().
		Where(entbonus.DonorResponseID(donorResponseID)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get bonuses by donor response ID: %w", err)
	}

	result := make([]*bonusmodel.Bonus, len(bonuses))
	for i, b := range bonuses {
		result[i] = &bonusmodel.Bonus{
			ID:           b.ID,
			UserID:       &b.UserID,
			PartnerName:  b.PartnerName,
			Description:  b.Description,
			Target:       b.Target.String(),
			Recipient:    b.Recipient.String(),
			Category:     b.Category.String(),
			Subcategory:  &b.Subcategory,
			PromoCode:    b.PromoCode,
			ExpiresAt:    b.ExpiresAt,
			PlatformName: b.PlatformName,
			PlatformURL:  &b.PlatformURL,
			Stage:        string(b.Stage),
			CreatedAt:    b.CreatedAt,
			UpdatedAt:    b.UpdatedAt,
			AssignedAt:   b.AssignedAt,
			DeletedAt:    b.DeletedAt,
		}
	}

	return result, nil
}
