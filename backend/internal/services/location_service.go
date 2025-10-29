package services

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/models"
	repositories "github.com/artesipov-alt/odnoi-krovi-app/internal/repositories/interfaces"
)

// LocationService определяет интерфейс для бизнес-логики локаций
type LocationService interface {
	// CreateLocation создает новую локацию в системе
	CreateLocation(ctx context.Context, locationData LocationCreate) (*models.Location, error)

	// GetLocationByID получает локацию по ID
	GetLocationByID(ctx context.Context, id int) (*models.Location, error)

	// GetAllLocations получает все локации
	GetAllLocations(ctx context.Context) ([]*models.Location, error)

	// UpdateLocation обновляет информацию о локации
	UpdateLocation(ctx context.Context, id int, updates LocationUpdate) error

	// DeleteLocation удаляет локацию по ID
	DeleteLocation(ctx context.Context, id int) error
}

// LocationCreate содержит данные для создания локации
type LocationCreate struct {
	Name string `json:"name" validate:"required,min=2,max=255"`
}

// LocationUpdate содержит поля, которые можно обновить для локации
type LocationUpdate struct {
	Name *string `json:"name,omitempty" validate:"omitempty,min=2,max=255"`
}

// LocationServiceImpl реализует LocationService
type LocationServiceImpl struct {
	locationRepo repositories.LocationRepository
}

// NewLocationService создает новый сервис локаций
func NewLocationService(locationRepo repositories.LocationRepository) *LocationServiceImpl {
	return &LocationServiceImpl{
		locationRepo: locationRepo,
	}
}

// CreateLocation создает новую локацию в системе
func (s *LocationServiceImpl) CreateLocation(ctx context.Context, locationData LocationCreate) (*models.Location, error) {
	// Проверяем, существует ли локация с таким названием
	exists, err := s.locationRepo.ExistsByName(ctx, locationData.Name)
	if err != nil {
		return nil, apperrors.Internal(err, "не удалось проверить существование локации")
	}

	if exists {
		return nil, apperrors.Conflict("локация с таким названием уже существует").WithDetails(map[string]interface{}{
			"name": locationData.Name,
		})
	}

	// Создаем новую локацию
	location := &models.Location{
		Name: locationData.Name,
	}

	if err := s.locationRepo.Create(ctx, location); err != nil {
		return nil, apperrors.Internal(err, "не удалось создать локацию")
	}

	return location, nil
}

// GetLocationByID получает локацию по ID
func (s *LocationServiceImpl) GetLocationByID(ctx context.Context, id int) (*models.Location, error) {
	location, err := s.locationRepo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.Internal(err, "не удалось получить локацию")
	}

	if location == nil {
		return nil, apperrors.NotFound("локация не найдена").WithDetails(map[string]interface{}{
			"location_id": id,
		})
	}

	return location, nil
}

// GetAllLocations получает все локации
func (s *LocationServiceImpl) GetAllLocations(ctx context.Context) ([]*models.Location, error) {
	locations, err := s.locationRepo.GetAll(ctx)
	if err != nil {
		return nil, apperrors.Internal(err, "не удалось получить список локаций")
	}

	return locations, nil
}

// UpdateLocation обновляет информацию о локации
func (s *LocationServiceImpl) UpdateLocation(ctx context.Context, id int, updates LocationUpdate) error {
	// Получаем существующую локацию
	location, err := s.locationRepo.GetByID(ctx, id)
	if err != nil {
		return apperrors.Internal(err, "не удалось получить локацию")
	}

	if location == nil {
		return apperrors.NotFound("локация не найдена").WithDetails(map[string]interface{}{
			"location_id": id,
		})
	}

	// Если обновляется название, проверяем уникальность
	if updates.Name != nil && *updates.Name != location.Name {
		exists, err := s.locationRepo.ExistsByName(ctx, *updates.Name)
		if err != nil {
			return apperrors.Internal(err, "не удалось проверить существование локации")
		}

		if exists {
			return apperrors.Conflict("локация с таким названием уже существует").WithDetails(map[string]interface{}{
				"name": *updates.Name,
			})
		}

		location.Name = *updates.Name
	}

	// Сохраняем обновленную локацию
	if err := s.locationRepo.Update(ctx, location); err != nil {
		return apperrors.Internal(err, "не удалось обновить локацию")
	}

	return nil
}

// DeleteLocation удаляет локацию по ID
func (s *LocationServiceImpl) DeleteLocation(ctx context.Context, id int) error {
	// Проверяем, существует ли локация
	location, err := s.locationRepo.GetByID(ctx, id)
	if err != nil {
		return apperrors.Internal(err, "не удалось получить локацию")
	}

	if location == nil {
		return apperrors.NotFound("локация не найдена").WithDetails(map[string]interface{}{
			"location_id": id,
		})
	}

	// Удаляем локацию
	if err := s.locationRepo.Delete(ctx, id); err != nil {
		return apperrors.Internal(err, "не удалось удалить локацию")
	}

	return nil
}
