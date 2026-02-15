package pg

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/artesipov-alt/odnoi-krovi-app/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/breed"
)

// EntBreedRepository implements BreedRepository using ENT
type EntBreedRepository struct {
	client *ent.Client
}

// NewEntBreedRepository creates a new ENT breed repository
func NewEntBreedRepository(client *ent.Client) *EntBreedRepository {
	return &EntBreedRepository{
		client: client,
	}
}

// GetAll returns all breeds from the database
func (r *EntBreedRepository) GetAll(ctx context.Context) ([]*ent.Breed, error) {
	breeds, err := r.client.Breed.Query().
		Order(ent.Asc(breed.FieldName)).
		All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get all breeds: %w", err)
	}

	return breeds, nil
}

// GetByID retrieves a breed by its ID
func (r *EntBreedRepository) GetByID(ctx context.Context, id string) (*ent.Breed, error) {
	if id == "" {
		return nil, errors.New("invalid breed ID")
	}

	b, err := r.client.Breed.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("breed with id %s not found: %w", id, err)
		}
		return nil, fmt.Errorf("failed to get breed by id %s: %w", id, err)
	}

	return b, nil
}

// GetByPetType retrieves breeds by pet type
func (r *EntBreedRepository) GetByPetType(ctx context.Context, petType breed.Type) ([]*ent.Breed, error) {
	if petType == "" {
		return nil, errors.New("pet type cannot be empty")
	}

	breeds, err := r.client.Breed.Query().
		Where(breed.TypeEQ(petType)).
		All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get breeds for pet type %s: %w", petType, err)
	}

	// Сортировка: для собак и кошек "МЕТИС" первым, остальные по алфавиту
	if petType == breed.TypeDog || petType == breed.TypeCat {
		var metis *ent.Breed
		var others []*ent.Breed
		for _, b := range breeds {
			if b.Name == "МЕТИС" {
				metis = b
			} else {
				others = append(others, b)
			}
		}
		sort.Slice(others, func(i, j int) bool {
			return others[i].Name < others[j].Name
		})
		if metis != nil {
			result := []*ent.Breed{metis}
			result = append(result, others...)
			return result, nil
		}
		return others, nil
	}

	// Для других типов сортировка по алфавиту
	sort.Slice(breeds, func(i, j int) bool {
		return breeds[i].Name < breeds[j].Name
	})

	return breeds, nil
}

// Create creates a new breed in the database
func (r *EntBreedRepository) Create(ctx context.Context, b *ent.Breed) (*ent.Breed, error) {
	if b == nil {
		return nil, errors.New("breed cannot be nil")
	}

	newBreed, err := r.client.Breed.Create().
		SetName(b.Name).
		SetType(b.Type).
		Save(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to create breed: %w", err)
	}

	return newBreed, nil
}

// Update updates an existing breed in the database
func (r *EntBreedRepository) Update(ctx context.Context, b *ent.Breed) (*ent.Breed, error) {
	if b == nil {
		return nil, errors.New("breed cannot be nil")
	}

	if b.ID == "" {
		return nil, errors.New("invalid breed ID")
	}

	updatedBreed, err := r.client.Breed.UpdateOneID(b.ID).
		SetName(b.Name).
		SetType(b.Type).
		Save(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("breed with id %s not found", b.ID)
		}
		return nil, fmt.Errorf("failed to update breed: %w", err)
	}

	return updatedBreed, nil
}

// Delete deletes a breed by its ID
func (r *EntBreedRepository) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("invalid breed ID")
	}

	err := r.client.Breed.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return fmt.Errorf("breed with id %s not found", id)
		}
		return fmt.Errorf("failed to delete breed: %w", err)
	}

	return nil
}

// ExistsByName checks if a breed with the given name exists
func (r *EntBreedRepository) ExistsByName(ctx context.Context, name string) (bool, error) {
	if name == "" {
		return false, errors.New("breed name cannot be empty")
	}

	exists, err := r.client.Breed.Query().
		Where(breed.Name(name)).
		Exist(ctx)

	if err != nil {
		return false, fmt.Errorf("failed to check breed existence by name %s: %w", name, err)
	}

	return exists, nil
}
