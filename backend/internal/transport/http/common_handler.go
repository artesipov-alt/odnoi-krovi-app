package http

import (
	"context"
	"net/http"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/application/analytics/query"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/transport/http/dto/common"

	"github.com/danielgtaylor/huma/v2"
)

type CommonHandler struct {
	portalStartsHandler *query.PortalStatsHandler
}

// NewDonorHandler creates a new handler for donor-related operations.
func NewCommonHandler(
	portalStatsHandler *query.PortalStatsHandler,
) *CommonHandler {
	return &CommonHandler{
		portalStartsHandler: portalStatsHandler,
	}
}

// Register регистрирует маршруты заявок на поиск крови в Huma API
func (h *CommonHandler) Register(api huma.API) {
	// Получить статистику портала
	huma.Register(api, huma.Operation{
		OperationID: "get-portal-stats",
		Method:      http.MethodGet,
		Path:        "/v1/portal/stats",
		Summary:     "Получить статистику портала",
		Description: "Возвращает статистику портала",
		Tags:        []string{"admin-v1"},
	}, h.GetPortalStats)
}

func (h *CommonHandler) GetPortalStats(ctx context.Context, req *struct{}) (*common.BodyOutput[common.PortalStats], error) {
	stats, err := h.portalStartsHandler.Handle(ctx)
	if err != nil {
		return nil, err
	}
	return &common.BodyOutput[common.PortalStats]{
		Body: common.PortalStats{
			TotalUsers:             stats.TotalUsers,
			UsersWithPhone:         stats.UsersWithPhone,
			PhoneConversionPercent: stats.PhoneConversionPercent,
			VerifiedUsers:          stats.VerifiedUsers,
			UnverifiedUsers:        stats.UnverifiedUsers,
			TotalPets:              stats.TotalPets,
			ActiveBloodRequests:    stats.ActiveBloodRequests,
			TotalDonations:         stats.TotalDonations,
			CompletedDonations:     stats.CompletedDonations,
		},
	}, nil
}
