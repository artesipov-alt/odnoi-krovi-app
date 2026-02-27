package reference

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/ent/bloodgroup"
	"github.com/artesipov-alt/odnoi-krovi-app/infa/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/user"
)

type ReferenceService struct {
	bloodInfoRepo domain.BloodInfoRepository
	breedsRepo    pet.BreedRepository
	locationRepo  user.LocationRepository
}

// NewReferenceService создает новый экземпляр ReferenceService
func NewReferenceService(bloodInfoRepo domain.BloodInfoRepository, breedsRepo pet.BreedRepository, locationRepo user.LocationRepository) *ReferenceService {
	return &ReferenceService{
		bloodInfoRepo: bloodInfoRepo,
		breedsRepo:    breedsRepo,
		locationRepo:  locationRepo,
	}
}

// GetAllBreeds возвращает все породы
func (s *ReferenceService) GetAllBreeds(ctx context.Context) ([]*ent.Breed, error) {
	return s.breedsRepo.GetAll(ctx)
}

// GetAllLocations возвращает все локации
func (s *ReferenceService) GetAllLocations(ctx context.Context) ([]*ent.Location, error) {
	return s.locationRepo.GetAll(ctx)
}

// GetAllBloodComponents возвращает все компоненты крови
func (s *ReferenceService) GetAllBloodComponents(ctx context.Context) ([]*ent.BloodComponent, error) {
	return s.bloodInfoRepo.AllComponents(ctx)
}

// GetBloodGroupsByPetType возвращает группы крови по типу животного
func (s *ReferenceService) GetBloodGroupsByPetType(ctx context.Context, petType bloodgroup.PetType) ([]*ent.BloodGroup, error) {
	return s.bloodInfoRepo.BloodGroupsByPetType(ctx, petType)
}
