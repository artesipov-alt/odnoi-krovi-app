package job

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/application/bloodsearch/cmd"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/ports"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/presistance/cache"
)

type NotificationJob struct {
	db           *sql.DB
	closeHandler *cmd.CloseRequestHandler
	publisher    ports.EventPublisher
	cache        cache.NotificationCache
}

func NewNotificationJob(db *sql.DB, closeHandler *cmd.CloseRequestHandler, publisher ports.EventPublisher, cache cache.NotificationCache) *NotificationJob {
	return &NotificationJob{
		db:           db,
		closeHandler: closeHandler,
		publisher:    publisher,
		cache:        cache,
	}
}

func (n *NotificationJob) Run(ctx context.Context) {
	n.checkDonorWaiting(ctx)
	n.checkRecipientInactive6h(ctx)
	n.checkRecipientInactive12h(ctx)
	n.checkRecipientEmptyShowcase24h(ctx)
	n.checkRecipientEmptyShowcase48h(ctx)
	n.checkDonorNotAccepted(ctx)
}

func (n *NotificationJob) checkDonorNotAccepted(ctx context.Context) {
	rows, err := n.db.QueryContext(ctx, queryDonorNotAccepted)
	if err != nil {
		slog.Error("checkDonorNotAccepted: query failed", "err", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			responseID       string
			recipientPetName string
			recipientBG      string
			volume           float64
			telegramID       sql.NullString
			maxID            sql.NullString
		)
		if err := rows.Scan(&responseID, &recipientPetName, &recipientBG, &volume, &telegramID, &maxID); err != nil {
			slog.Error("checkDonorNotAccepted: scan failed", "err", err)
			continue
		}

		sent, err := n.cache.WasSent(ctx, ports.NotifDonorNotAccepted, responseID)
		if err != nil {
			slog.Error("checkDonorNotAccepted: cache check failed", "responseID", responseID, "err", err)
			continue
		}
		if sent {
			continue
		}

		err = n.publisher.PublishNotification(ctx, ports.Notification{
			Type: ports.NotifDonorNotAccepted,
			Targets: ports.NotifTargets{
				TelegramID: telegramID.String,
				MaxID:      maxID.String,
			},
			Payload: map[string]any{
				"recipientPetName":    recipientPetName,
				"recipientBloodGroup": recipientBG,
				"volume":              volume,
			},
			CreatedAt: time.Now(),
		})
		if err != nil {
			slog.Error("checkDonorNotAccepted: publish failed", "responseID", responseID, "err", err)
			continue
		}

		if err := n.cache.MarkSent(ctx, ports.NotifDonorNotAccepted, responseID, time.Hour); err != nil {
			slog.Error("checkDonorNotAccepted: mark sent failed", "responseID", responseID, "err", err)
		}
	}
}

func (n *NotificationJob) checkDonorWaiting(ctx context.Context) {
	rows, err := n.db.QueryContext(ctx, queryDonorWaiting)
	if err != nil {
		slog.Error("checkDonorWaiting: query failed", "err", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			responseID   string
			donorPetName string
			donorBG      string
			telegramID   sql.NullString
			maxID        sql.NullString
		)
		if err := rows.Scan(&responseID, &donorPetName, &donorBG, &telegramID, &maxID); err != nil {
			slog.Error("checkDonorWaiting: scan failed", "err", err)
			continue
		}

		sent, err := n.cache.WasSent(ctx, ports.NotifRecipientDonorWaiting, responseID)
		if err != nil {
			slog.Error("checkDonorWaiting: cache check failed", "responseID", responseID, "err", err)
			continue
		}
		if sent {
			continue
		}

		err = n.publisher.PublishNotification(ctx, ports.Notification{
			Type: ports.NotifRecipientDonorWaiting,
			Targets: ports.NotifTargets{
				TelegramID: telegramID.String,
				MaxID:      maxID.String,
			},
			Payload: map[string]any{
				"donorPetName":    donorPetName,
				"donorBloodGroup": donorBG,
			},
			CreatedAt: time.Now(),
		})
		if err != nil {
			slog.Error("checkDonorWaiting: publish failed", "responseID", responseID, "err", err)
			continue
		}

		if err := n.cache.MarkSent(ctx, ports.NotifRecipientDonorWaiting, responseID, 30*time.Minute); err != nil {
			slog.Error("checkDonorWaiting: mark sent failed", "responseID", responseID, "err", err)
		}
	}
}

func (n *NotificationJob) checkRecipientInactive6h(ctx context.Context) {
	rows, err := n.db.QueryContext(ctx, queryRecipientInactive6h)
	if err != nil {
		slog.Error("checkRecipientInactive6h: query failed", "err", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			requestID  string
			telegramID sql.NullString
			maxID      sql.NullString
		)
		if err := rows.Scan(&requestID, &telegramID, &maxID); err != nil {
			slog.Error("checkRecipientInactive6h: scan failed", "err", err)
			continue
		}

		sent, err := n.cache.WasSent(ctx, ports.NotifRecipientInactiveWarning, requestID)
		if err != nil {
			slog.Error("checkRecipientInactive6h: cache check failed", "requestID", requestID, "err", err)
			continue
		}
		if sent {
			continue
		}

		err = n.publisher.PublishNotification(ctx, ports.Notification{
			Type: ports.NotifRecipientInactiveWarning,
			Targets: ports.NotifTargets{
				TelegramID: telegramID.String,
				MaxID:      maxID.String,
			},
			Payload:   map[string]any{},
			CreatedAt: time.Now(),
		})
		if err != nil {
			slog.Error("checkRecipientInactive6h: publish failed", "requestID", requestID, "err", err)
			continue
		}

		if err := n.cache.MarkSent(ctx, ports.NotifRecipientInactiveWarning, requestID, 6*time.Hour); err != nil {
			slog.Error("checkRecipientInactive6h: mark sent failed", "requestID", requestID, "err", err)
		}
	}
}

func (n *NotificationJob) checkRecipientInactive12h(ctx context.Context) {
	rows, err := n.db.QueryContext(ctx, queryRecipientInactive12h)
	if err != nil {
		slog.Error("checkRecipientInactive12h: query failed", "err", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			requestID  string
			telegramID sql.NullString
			maxID      sql.NullString
		)
		if err := rows.Scan(&requestID, &telegramID, &maxID); err != nil {
			slog.Error("checkRecipientInactive12h: scan failed", "err", err)
			continue
		}

		sent, err := n.cache.WasSent(ctx, ports.NotifRecipientSearchClosed, requestID)
		if err != nil {
			slog.Error("checkRecipientInactive12h: cache check failed", "requestID", requestID, "err", err)
			continue
		}
		if sent {
			continue
		}

		// Сначала закрываем поиск
		if err := n.closeHandler.Handle(ctx, requestID); err != nil {
			slog.Error("checkRecipientInactive12h: close request failed", "requestID", requestID, "err", err)
			continue
		}

		// Потом уведомляем
		err = n.publisher.PublishNotification(ctx, ports.Notification{
			Type: ports.NotifRecipientSearchClosed,
			Targets: ports.NotifTargets{
				TelegramID: telegramID.String,
				MaxID:      maxID.String,
			},
			Payload:   map[string]any{},
			CreatedAt: time.Now(),
		})
		if err != nil {
			slog.Error("checkRecipientInactive12h: publish failed", "requestID", requestID, "err", err)
			continue
		}

		if err := n.cache.MarkSent(ctx, ports.NotifRecipientSearchClosed, requestID, 12*time.Hour); err != nil {
			slog.Error("checkRecipientInactive12h: mark sent failed", "requestID", requestID, "err", err)
		}
	}
}

func (n *NotificationJob) checkRecipientEmptyShowcase24h(ctx context.Context) {
	rows, err := n.db.QueryContext(ctx, queryRecipientEmptyShowcase24h)
	if err != nil {
		slog.Error("checkRecipientEmptyShowcase24h: query failed", "err", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			requestID  string
			telegramID sql.NullString
			maxID      sql.NullString
		)
		if err := rows.Scan(&requestID, &telegramID, &maxID); err != nil {
			slog.Error("checkRecipientEmptyShowcase24h: scan failed", "err", err)
			continue
		}

		sent, err := n.cache.WasSent(ctx, ports.NotifRecipientEmptyShowcase, requestID)
		if err != nil {
			slog.Error("checkRecipientEmptyShowcase24h: cache check failed", "requestID", requestID, "err", err)
			continue
		}
		if sent {
			continue
		}

		err = n.publisher.PublishNotification(ctx, ports.Notification{
			Type: ports.NotifRecipientEmptyShowcase,
			Targets: ports.NotifTargets{
				TelegramID: telegramID.String,
				MaxID:      maxID.String,
			},
			Payload: map[string]any{
				"requestId": requestID,
			},
			CreatedAt: time.Now(),
		})
		if err != nil {
			slog.Error("checkRecipientEmptyShowcase24h: publish failed", "requestID", requestID, "err", err)
			continue
		}

		if err := n.cache.MarkSent(ctx, ports.NotifRecipientEmptyShowcase, requestID, 24*time.Hour); err != nil {
			slog.Error("checkRecipientEmptyShowcase24h: mark sent failed", "requestID", requestID, "err", err)
		}
	}
}

func (n *NotificationJob) checkRecipientEmptyShowcase48h(ctx context.Context) {
	rows, err := n.db.QueryContext(ctx, queryRecipientEmptyShowcase48h)
	if err != nil {
		slog.Error("checkRecipientEmptyShowcase48h: query failed", "err", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			requestID  string
			telegramID sql.NullString
			maxID      sql.NullString
		)
		if err := rows.Scan(&requestID, &telegramID, &maxID); err != nil {
			slog.Error("checkRecipientEmptyShowcase48h: scan failed", "err", err)
			continue
		}

		// Если пользователь нажал "Да" — не закрываем
		yesPressed, err := n.cache.WasYesPressed(ctx, requestID)
		if err != nil {
			slog.Error("checkRecipientEmptyShowcase48h: yes pressed check failed", "requestID", requestID, "err", err)
			continue
		}
		if yesPressed {
			continue
		}

		sent, err := n.cache.WasSent(ctx, ports.NotifRecipientSearchClosedInactive, requestID)
		if err != nil {
			slog.Error("checkRecipientEmptyShowcase48h: cache check failed", "requestID", requestID, "err", err)
			continue
		}
		if sent {
			continue
		}

		// Сначала закрываем поиск
		if err := n.closeHandler.Handle(ctx, requestID); err != nil {
			slog.Error("checkRecipientEmptyShowcase48h: close request failed", "requestID", requestID, "err", err)
			continue
		}

		// Потом уведомляем
		err = n.publisher.PublishNotification(ctx, ports.Notification{
			Type: ports.NotifRecipientSearchClosedInactive,
			Targets: ports.NotifTargets{
				TelegramID: telegramID.String,
				MaxID:      maxID.String,
			},
			Payload:   map[string]any{},
			CreatedAt: time.Now(),
		})
		if err != nil {
			slog.Error("checkRecipientEmptyShowcase48h: publish failed", "requestID", requestID, "err", err)
			continue
		}

		if err := n.cache.MarkSent(ctx, ports.NotifRecipientSearchClosedInactive, requestID, 48*time.Hour); err != nil {
			slog.Error("checkRecipientEmptyShowcase48h: mark sent failed", "requestID", requestID, "err", err)
		}
	}
}
