package bonus

import (
	"context"

	bonusmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bonus/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/common"
)

// Repository defines the interface for bonus data access.
type Repository interface {
	// Write methods

	// создает несколько бонусов в одной транзакции
	CreateBatch(ctx context.Context, bonuses []*bonusmodel.Bonus) error

	// открепляет зарезервированные бонусы от пользователя по конкретному отклику донора
	UnassignReservedBonuses(ctx context.Context, userID string, petType common.PetType, donorResponseID string) error

	// назначает бонусы пользователю (UserID, DonorResponseID, stage = reserved)
	AssignBonuses(ctx context.Context, bonusIDs []string, userID string, donorResponseID string) error

	// подтверждает бонусы пользователя (stage = unused, UserID не очищается)
	ConfirmBonuses(ctx context.Context, userID string, petType common.PetType) error

	// помечает бонусы как использованные (stage = used)
	MarkBonusesAsUsed(ctx context.Context, userID string, petType common.PetType) error

	// увеличивает счетчик приоритетного поиска на 1
	AddPrioritySearch(ctx context.Context, id string) error

	// уменьшает счетчик приоритетного поиска на 1
	SubtractPrioritySearch(ctx context.Context, id string) error

	// Read methods

	// возвращает список промокодов, которые уже существуют в БД
	ExistsByPromoCodes(ctx context.Context, codes []string) ([]string, error)

	// возвращает доступные (неназначенные) бонусы по фильтрам
	GetAvailableBonuses(ctx context.Context, petType common.PetType) ([]*bonusmodel.Bonus, error)

	// возвращает последний бонус пользователя по UpdatedAt
	GetLastBonus(ctx context.Context, userID string) (*bonusmodel.Bonus, error)

	// возвращает назначенные бонусы пользователя
	GetAssignedBonuses(ctx context.Context, userID string) ([]*bonusmodel.Bonus, error)

	// возвращает бонусы по ID отклика донора
	GetBonusesByDonorResponseID(ctx context.Context, donorResponseID string) ([]*bonusmodel.Bonus, error)
}
