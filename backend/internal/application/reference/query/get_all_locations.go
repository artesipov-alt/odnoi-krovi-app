package query

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/reference"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/reference/model"
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
	return locations, nil
}
