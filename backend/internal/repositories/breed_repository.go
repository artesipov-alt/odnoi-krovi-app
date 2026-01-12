package repositories

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/breed"
)

// BreedRepository определяет интерфейс для операций с данными пород
type BreedRepository interface {
	// GetAll возвращает все породы из базы данных
	GetAll(ctx context.Context) ([]*ent.Breed, error)

	// GetByID получает породу по её ID
	GetByID(ctx context.Context, id int) (*ent.Breed, error)

	// GetByPetType получает породы по типу животного
	GetByPetType(ctx context.Context, petType breed.Type) ([]*ent.Breed, error)

	// Create создает новую породу в базе данных
	Create(ctx context.Context, b *ent.Breed) (*ent.Breed, error)

	// Update обновляет существующую породу в базе данных
	Update(ctx context.Context, b *ent.Breed) (*ent.Breed, error)

	// Delete удаляет породу по её ID
	Delete(ctx context.Context, id int) error

	// ExistsByName проверяет, существует ли порода с заданным названием
	ExistsByName(ctx context.Context, name string) (bool, error)
}
