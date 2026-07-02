// internal/application/bloodsearch/cmd/notification_respond.go

package cmd

import (
	"context"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/ports"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/presistance/cache"
)

type NotificationAction string

const (
	NotificationActionYes NotificationAction = "yes"
	NotificationActionNo  NotificationAction = "no"
)

type NotificationRespondInput struct {
	BloodRequestID string
	Action         NotificationAction
}

type NotificationRespondHandler struct {
	closeHandler *CloseRequestHandler
	cache        cache.NotificationCache
}

func NewNotificationRespondHandler(
	closeHandler *CloseRequestHandler,
	cache cache.NotificationCache,
) *NotificationRespondHandler {
	return &NotificationRespondHandler{
		closeHandler: closeHandler,
		cache:        cache,
	}
}

func (h *NotificationRespondHandler) Handle(ctx context.Context, input NotificationRespondInput) error {
	switch input.Action {
	case NotificationActionYes:
		if err := h.cache.MarkYesPressed(ctx, input.BloodRequestID); err != nil {
			return apperrors.Internal(err, "failed to mark yes pressed")
		}

		// Сброс таймера: пользователь сказал "Да", обновляем sent-маркер еще на 24ч
		if err := h.cache.MarkSent(ctx, ports.NotifRecipientEmptyShowcase, input.BloodRequestID, 24*time.Hour); err != nil {
			return apperrors.Internal(err, "failed to reset sent timer")
		}

		return nil

	case NotificationActionNo:
		if err := h.closeHandler.Handle(ctx, input.BloodRequestID); err != nil {
			return err
		}
		return nil

	default:
		return apperrors.BadRequest("неизвестное действие")
	}
}
