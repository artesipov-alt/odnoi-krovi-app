package partner

import (
	"context"
	"time"

	authmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/auth/model"
	partnermodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/partner/model"
)

type Repository interface {
	// Write methods

	// создает нового партнера
	CreatePartner(ctx context.Context, name, apiKey, role, status, description string) (*partnermodel.Partner, error)

	// обновляет статус партнера
	UpdateStatus(ctx context.Context, id string, status string) error

	// обновляет время последнего использования
	UpdateLastUsedAt(ctx context.Context, id string, lastUsedAt time.Time) error

	// удаляет партнера по ID
	Delete(ctx context.Context, id string) error

	// Read methods

	// возвращает партнера по API-ключу
	GetByAPIKey(ctx context.Context, apiKey string) (*partnermodel.Partner, error)

	// возвращает партнера по ID
	GetByID(ctx context.Context, id string) (*partnermodel.Partner, error)

	// проверяет, существует ли партнер с заданным API-ключом
	ExistsByAPIKey(ctx context.Context, apiKey string) (bool, error)

	// возвращает все identity партнера
	GetPartnerIdentities(ctx context.Context, partnerID string) ([]*authmodel.Identity, error)
}
