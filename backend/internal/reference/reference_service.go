package reference

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/bloodgroup"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/breed"
)

type ReferenceService struct {
	bloodInfoRepo domain.BloodInfoRepository
	breedsRepo    domain.BreedRepository
	locationRepo  domain.LocationRepository
}

// NewReferenceService создает новый экземпляр ReferenceService
func NewReferenceService(bloodInfoRepo domain.BloodInfoRepository, breedsRepo domain.BreedRepository, locationRepo domain.LocationRepository) *ReferenceService {
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

// GetBreedsByPetType возвращает породы по типу животного
func (s *ReferenceService) GetBreedsByPetType(ctx context.Context, petType breed.Type) ([]*ent.Breed, error) {
	return s.breedsRepo.GetByPetType(ctx, petType)
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

// AllComponents возвращает все компоненты крови
func (s *ReferenceService) AllComponents(ctx context.Context) ([]*ent.BloodComponent, error) {
	return s.bloodInfoRepo.AllComponents(ctx)
}

// GetByPetType возвращает породы по типу животного
func (s *ReferenceService) GetByPetType(ctx context.Context, petType breed.Type) ([]*ent.Breed, error) {
	return s.breedsRepo.GetByPetType(ctx, petType)
}
