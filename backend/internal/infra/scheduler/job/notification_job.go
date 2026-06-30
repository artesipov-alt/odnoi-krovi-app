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
				"recipient_pet_name":    recipientPetName,
				"recipient_blood_group": recipientBG,
				"volume":                volume,
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
				"donor_pet_name":    donorPetName,
				"donor_blood_group": donorBG,
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

}

func (n *NotificationJob) checkRecipientEmptyShowcase24h(ctx context.Context) {

}

func (n *NotificationJob) checkRecipientEmptyShowcase48h(ctx context.Context) {

}
