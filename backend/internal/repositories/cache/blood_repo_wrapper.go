package cache

import (
	"context"
	"fmt"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/cache"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/cache/interfaces"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/models"
	repositories "github.com/artesipov-alt/odnoi-krovi-app/internal/repositories/interfaces"
)

// CachedBloodRepository реализует кеширующий репозиторий для работы с группами крови
type CachedBloodRepository struct {
	repo  repositories.BloodRepository
	cache interfaces.Cache
}

// NewCachedBloodRepository создает новый экземпляр кеширующего репозитория
func NewCachedBloodRepository(repo repositories.BloodRepository, cache interfaces.Cache) *CachedBloodRepository {
	return &CachedBloodRepository{
		repo:  repo,
		cache: cache,
	}
}

// GetAllComponents возвращает все компоненты крови с кешированием
func (r *CachedBloodRepository) GetAllComponents(ctx context.Context) ([]models.BloodComponent, error) {
	cacheKey := fmt.Sprintf(cache.BloodTypesListKey)

	// Пытаемся получить из кэша
	var components []models.BloodComponent
	if err := r.cache.GetJSON(ctx, cacheKey, &components); err == nil {

		return components, nil
	}

	// Получаем из БД

	components, err := r.repo.GetAllComponents(ctx)
	if err != nil {
		return nil, err
	}

	// Сохраняем в кэш

	r.cache.SetJSON(ctx, cacheKey, components, cache.LongTTL)

	return components, nil
}

// GetComponentByID возвращает компонент крови по ID с кешированием
func (r *CachedBloodRepository) GetComponentByID(ctx context.Context, id int) (*models.BloodComponent, error) {
	cacheKey := fmt.Sprintf(cache.BloodComponentByIDKey, id)

	// Пытаемся получить из кэша
	var component models.BloodComponent
	if err := r.cache.GetJSON(ctx, cacheKey, &component); err == nil {

		return &component, nil
	}

	// Получаем из БД

	componentPtr, err := r.repo.GetComponentByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Сохраняем в кэш

	r.cache.SetJSON(ctx, cacheKey, componentPtr, cache.LongTTL)

	return componentPtr, nil
}

// GetBloodGroupsByPetType возвращает группы крови по типу животного с кешированием
func (r *CachedBloodRepository) GetBloodGroupsByPetType(ctx context.Context, petType models.PetType) ([]*models.BloodGroup, error) {
	cacheKey := fmt.Sprintf(cache.BloodGroupsByPetTypeKey, petType)

	// Пытаемся получить из кэша
	var bloodGroups []*models.BloodGroup
	if err := r.cache.GetJSON(ctx, cacheKey, &bloodGroups); err == nil {

		return bloodGroups, nil
	}

	// Получаем из БД

	bloodGroups, err := r.repo.GetBloodGroupsByPetType(ctx, petType)
	if err != nil {
		return nil, err
	}

	// Сохраняем в кэш

	r.cache.SetJSON(ctx, cacheKey, bloodGroups, cache.LongTTL)

	return bloodGroups, nil
}
