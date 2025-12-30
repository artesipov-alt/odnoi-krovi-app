package cache

import (
	"context"
	"fmt"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/cache"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/models"
	repositories "github.com/artesipov-alt/odnoi-krovi-app/internal/repositories"
)

// CachedBloodInfoRepository реализует кеширующий репозиторий для работы с группами крови
type CachedBloodInfoRepository struct {
	repo  repositories.BloodInfoRepository
	cache cache.ICache
}

// NewCachedBloodInfoRepository создает новый экземпляр кеширующего репозитория
func NewCachedBloodInfoRepository(repo repositories.BloodInfoRepository, cache cache.ICache) *CachedBloodInfoRepository {
	return &CachedBloodInfoRepository{
		repo:  repo,
		cache: cache,
	}
}

// GetAllComponents возвращает все компоненты крови с кешированием
func (r *CachedBloodInfoRepository) AllComponents(ctx context.Context) ([]models.BloodComponent, error) {
	cacheKey := fmt.Sprintf(cache.BloodTypesListKey)

	// Пытаемся получить из кэша
	var components []models.BloodComponent
	if err := r.cache.GetJSON(ctx, cacheKey, &components); err == nil {

		return components, nil
	}

	// Получаем из БД

	components, err := r.repo.AllComponents(ctx)
	if err != nil {
		return nil, err
	}

	// Сохраняем в кэш

	r.cache.SetJSON(ctx, cacheKey, components, cache.LongTTL)

	return components, nil
}

// GetComponentByID возвращает компонент крови по ID с кешированием
func (r *CachedBloodInfoRepository) ComponentByID(ctx context.Context, id int) (*models.BloodComponent, error) {
	cacheKey := fmt.Sprintf(cache.BloodComponentByIDKey, id)

	// Пытаемся получить из кэша
	var component models.BloodComponent
	if err := r.cache.GetJSON(ctx, cacheKey, &component); err == nil {

		return &component, nil
	}

	// Получаем из БД

	componentPtr, err := r.repo.ComponentByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Сохраняем в кэш

	r.cache.SetJSON(ctx, cacheKey, componentPtr, cache.LongTTL)

	return componentPtr, nil
}

// GetBloodGroupsByPetType возвращает группы крови по типу животного с кешированием
func (r *CachedBloodInfoRepository) BloodGroupsByPetType(ctx context.Context, petType models.PetType) ([]*models.BloodGroup, error) {
	cacheKey := fmt.Sprintf(cache.BloodGroupsByPetTypeKey, petType)

	// Пытаемся получить из кэша
	var bloodGroups []*models.BloodGroup
	if err := r.cache.GetJSON(ctx, cacheKey, &bloodGroups); err == nil {

		return bloodGroups, nil
	}

	// Получаем из БД

	bloodGroups, err := r.repo.BloodGroupsByPetType(ctx, petType)
	if err != nil {
		return nil, err
	}

	// Сохраняем в кэш

	r.cache.SetJSON(ctx, cacheKey, bloodGroups, cache.LongTTL)

	return bloodGroups, nil
}
