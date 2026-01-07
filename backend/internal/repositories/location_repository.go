package repositories

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/ent"
)

// LocationRepository определяет интерфейс для операций с данными локаций
type LocationRepository interface {

	// GetByID получает локацию по её ID
	GetByID(ctx context.Context, id int) (*ent.Location, error)

	// GetAll получает все локации из базы данных
	GetAll(ctx context.Context) ([]*ent.Location, error)
}
