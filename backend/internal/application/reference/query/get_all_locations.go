package query

import (
	"cmp"
	"context"
	"slices"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/reference"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/reference/model"
)

type GetAllLocationsHandler struct {
	locationRepo reference.LocationInfoRepository
}

func NewGetAllLocationsHandler(locationRepo reference.LocationInfoRepository) *GetAllLocationsHandler {
	return &GetAllLocationsHandler{
		locationRepo: locationRepo,
	}
}

func (h *GetAllLocationsHandler) Handle(ctx context.Context) ([]*model.Location, error) {
	locations, err := h.locationRepo.GetAll(ctx)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get all locations")
	}

	// Карта приоритетов для топ-регионов (ключи — строковые ID локаций)
	// Регионы с меньшим числом приоритета отображаются первыми в списке
	priority := map[string]int{
		"77": 1, // Москва
		"50": 2, // Московская область
		"78": 3, // Санкт-Петербург
		"47": 4, // Ленинградская область
	}

	// Сортировка: топ-регионы сначала (по приоритету), остальные — по алфавиту
	slices.SortStableFunc(locations, func(a, b *model.Location) int {
		pa, pb := priority[a.ID], priority[b.ID]

		// Оба — топ-регионы: сохраняем порядок приоритетов
		if pa != 0 && pb != 0 {
			return cmp.Compare(pa, pb)
		}

		// Только один — топ-регион: он идёт первым
		if pa != 0 {
			return -1
		}
		if pb != 0 {
			return 1
		}

		// Оба обычные: сортируем по алфавиту
		return cmp.Compare(a.Name, b.Name)
	})

	return locations, nil
}
