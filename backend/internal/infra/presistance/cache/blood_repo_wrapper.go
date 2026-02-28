package cache

import (
	"context"
	"fmt"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/bloodgroup"
)

// CachedBloodInfoRepository реализует кеширующий репозиторий для работы с группами крови
type CachedBloodInfoRepository struct {
	repo  domain.BloodInfoRepository
	cache ICache
}

// NewCachedBloodInfoRepository создает новый экземпляр кеширующего репозитория
func NewCachedBloodInfoRepository(repo domain.BloodInfoRepository, cache ICache) *CachedBloodInfoRepository {
	return &CachedBloodInfoRepository{
		repo:  repo,
		cache: cache,
	}
}

// AllComponents возвращает все компоненты крови с кешированием
func (r *CachedBloodInfoRepository) AllComponents(ctx context.Context) ([]*ent.BloodComponent, error) {
	cacheKey := BloodTypesListKey

	// Пытаемся получить из кэша
	var components []*ent.BloodComponent
	if err := r.cache.GetJSON(ctx, cacheKey, &components); err == nil {
		return components, nil
	}

	// Получаем из БД
	components, err := r.repo.AllComponents(ctx)
	if err != nil {
		return nil, err
	}

	// Сохраняем в кэш
	r.cache.SetJSON(ctx, cacheKey, components, LongTTL)

	return components, nil
}

// ComponentByID возвращает компонент крови по ID с кешированием
func (r *CachedBloodInfoRepository) ComponentByID(ctx context.Context, id string) (*ent.BloodComponent, error) {
	cacheKey := fmt.Sprintf(BloodComponentByIDKey, id)

	// Пытаемся получить из кэша
	var component ent.BloodComponent
	if err := r.cache.GetJSON(ctx, cacheKey, &component); err == nil {
		return &component, nil
	}

	// Получаем из БД
	componentPtr, err := r.repo.ComponentByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Сохраняем в кэш
	r.cache.SetJSON(ctx, cacheKey, componentPtr, LongTTL)

	return componentPtr, nil
}

// BloodGroupsByPetType возвращает группы крови по типу животного с кешированием
func (r *CachedBloodInfoRepository) BloodGroupsByPetType(ctx context.Context, petType bloodgroup.PetType) ([]*ent.BloodGroup, error) {
	cacheKey := fmt.Sprintf(BloodGroupsByPetTypeKey, petType)

	// Пытаемся получить из кэша
	var bloodGroups []*ent.BloodGroup
	if err := r.cache.GetJSON(ctx, cacheKey, &bloodGroups); err == nil {
		return bloodGroups, nil
	}

	// Получаем из БД
	bloodGroups, err := r.repo.BloodGroupsByPetType(ctx, petType)
	if err != nil {
		return nil, err
	}

	// Сохраняем в кэш
	r.cache.SetJSON(ctx, cacheKey, bloodGroups, LongTTL)

	return bloodGroups, nil
}
