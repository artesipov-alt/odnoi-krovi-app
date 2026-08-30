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
	n.checkAccepted12h(ctx)
	n.checkAccepted24h(ctx)
	n.checkVerifiedNoPets(ctx)
	n.checkNotVerified(ctx)
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

		if err := n.cache.MarkSent(ctx, ports.NotifDonorNotAccepted, responseID, 72*time.Hour); err != nil {
			slog.Error("checkDonorNotAccepted: mark sent failed", "responseID", responseID, "err", err)
		}

	}
	if err := rows.Err(); err != nil {
		slog.Error("checkDonorNotAccepted: rows iteration error", "err", err)
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
	if err := rows.Err(); err != nil {
		slog.Error("checkDonorWaiting: rows iteration error", "err", err)
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
	if err := rows.Err(); err != nil {
		slog.Error("checkRecipientInactive6h: rows iteration error", "err", err)
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
	if err := rows.Err(); err != nil {
		slog.Error("checkRecipientInactive12h: rows iteration error", "err", err)
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
			petName    string
			bloodGroup string
			volume     sql.NullFloat64
			telegramID sql.NullString
			maxID      sql.NullString
		)
		if err := rows.Scan(&requestID, &petName, &bloodGroup, &volume, &telegramID, &maxID); err != nil {
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
				"requestId":  requestID,
				"petName":    petName,
				"bloodGroup": bloodGroup,
				"volume":     volume.Float64,
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
	if err := rows.Err(); err != nil {
		slog.Error("checkRecipientEmptyShowcase24h: rows iteration error", "err", err)
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
	if err := rows.Err(); err != nil {
		slog.Error("checkRecipientEmptyShowcase48h: rows iteration error", "err", err)
	}
}

// checkAccepted12h — заявки в active/reserved_full с accepted-откликом старше 12ч.
// Отправляет два уведомления: реципиенту и принятому донору.
func (n *NotificationJob) checkAccepted12h(ctx context.Context) {
	rows, err := n.db.QueryContext(ctx, queryAccepted12h)
	if err != nil {
		slog.Error("checkAccepted12h: query failed", "err", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			responseID          string
			requestID           string
			recipientPetName    string
			recipientBloodGroup string
			donorPetName        string
			donorBloodGroup     string
			amount              float64
			recipientTelegramID sql.NullString
			recipientMaxID      sql.NullString
			donorTelegramID     sql.NullString
			donorMaxID          sql.NullString
		)
		if err := rows.Scan(
			&responseID, &requestID,
			&recipientPetName, &recipientBloodGroup,
			&donorPetName, &donorBloodGroup, &amount,
			&recipientTelegramID, &recipientMaxID,
			&donorTelegramID, &donorMaxID,
		); err != nil {
			slog.Error("checkAccepted12h: scan failed", "err", err)
			continue
		}

		// Уведомление реципиенту
		recipientSent, err := n.cache.WasSent(ctx, ports.NotifRecipientAcceptedReminder12h, responseID)
		if err != nil {
			slog.Error("checkAccepted12h: recipient cache check failed", "responseID", responseID, "err", err)
			continue
		}
		if !recipientSent {
			err = n.publisher.PublishNotification(ctx, ports.Notification{
				Type: ports.NotifRecipientAcceptedReminder12h,
				Targets: ports.NotifTargets{
					TelegramID: recipientTelegramID.String,
					MaxID:      recipientMaxID.String,
				},
				Payload: map[string]any{
					"requestId":       requestID,
					"donorPetName":    donorPetName,
					"donorBloodGroup": donorBloodGroup,
				},
				CreatedAt: time.Now(),
			})
			if err != nil {
				slog.Error("checkAccepted12h: recipient publish failed", "responseID", responseID, "err", err)
				continue
			}
			if err := n.cache.MarkSent(ctx, ports.NotifRecipientAcceptedReminder12h, responseID, 12*time.Hour); err != nil {
				slog.Error("checkAccepted12h: recipient mark sent failed", "responseID", responseID, "err", err)
			}
		}

		// Уведомление донору
		donorSent, err := n.cache.WasSent(ctx, ports.NotifDonorAcceptedReminder12h, responseID)
		if err != nil {
			slog.Error("checkAccepted12h: donor cache check failed", "responseID", responseID, "err", err)
			continue
		}
		if !donorSent {
			err = n.publisher.PublishNotification(ctx, ports.Notification{
				Type: ports.NotifDonorAcceptedReminder12h,
				Targets: ports.NotifTargets{
					TelegramID: donorTelegramID.String,
					MaxID:      donorMaxID.String,
				},
				Payload: map[string]any{
					"responseId":          responseID,
					"recipientPetName":    recipientPetName,
					"recipientBloodGroup": recipientBloodGroup,
					"volume":              amount,
				},
				CreatedAt: time.Now(),
			})
			if err != nil {
				slog.Error("checkAccepted12h: donor publish failed", "responseID", responseID, "err", err)
				continue
			}
			if err := n.cache.MarkSent(ctx, ports.NotifDonorAcceptedReminder12h, responseID, 12*time.Hour); err != nil {
				slog.Error("checkAccepted12h: donor mark sent failed", "responseID", responseID, "err", err)
			}
		}
	}
	if err := rows.Err(); err != nil {
		slog.Error("checkAccepted12h: rows iteration error", "err", err)
	}
}

// checkAccepted24h — заявки в active/reserved_full с accepted-откликом старше 24ч.
// Отправляет уведомление реципиенту и запускает сценарий «72 часа» (фактически —
// AutoConfirmJob сам найдёт отклик через 96ч от accepted, т.е. 72ч от этого напоминания).
func (n *NotificationJob) checkAccepted24h(ctx context.Context) {
	rows, err := n.db.QueryContext(ctx, queryAccepted24h)
	if err != nil {
		slog.Error("checkAccepted24h: query failed", "err", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			responseID          string
			requestID           string
			recipientPetName    string
			recipientBloodGroup string
			recipientTelegramID sql.NullString
			recipientMaxID      sql.NullString
		)
		if err := rows.Scan(
			&responseID, &requestID,
			&recipientPetName, &recipientBloodGroup,
			&recipientTelegramID, &recipientMaxID,
		); err != nil {
			slog.Error("checkAccepted24h: scan failed", "err", err)
			continue
		}

		sent, err := n.cache.WasSent(ctx, ports.NotifRecipientAcceptedReminder24h, responseID)
		if err != nil {
			slog.Error("checkAccepted24h: cache check failed", "responseID", responseID, "err", err)
			continue
		}
		if sent {
			continue
		}

		err = n.publisher.PublishNotification(ctx, ports.Notification{
			Type: ports.NotifRecipientAcceptedReminder24h,
			Targets: ports.NotifTargets{
				TelegramID: recipientTelegramID.String,
				MaxID:      recipientMaxID.String,
			},
			Payload: map[string]any{
				"requestId":           requestID,
				"recipientPetName":    recipientPetName,
				"recipientBloodGroup": recipientBloodGroup,
			},
			CreatedAt: time.Now(),
		})
		if err != nil {
			slog.Error("checkAccepted24h: publish failed", "responseID", responseID, "err", err)
			continue
		}

		if err := n.cache.MarkSent(ctx, ports.NotifRecipientAcceptedReminder24h, responseID, 24*time.Hour); err != nil {
			slog.Error("checkAccepted24h: mark sent failed", "responseID", responseID, "err", err)
		}
	}
	if err := rows.Err(); err != nil {
		slog.Error("checkAccepted24h: rows iteration error", "err", err)
	}
}

// checkVerifiedNoPets — верифицированные пользователи без питомцев.
// Первое уведомление через 24ч после верификации, повтор каждые 48ч (дедупликация —
// кеш с TTL 48ч), пока пользователь не добавит питомца. Лимита повторов нет.
func (n *NotificationJob) checkVerifiedNoPets(ctx context.Context) {
	rows, err := n.db.QueryContext(ctx, queryVerifiedNoPets)
	if err != nil {
		slog.Error("checkVerifiedNoPets: query failed", "err", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			userID     string
			telegramID sql.NullString
			maxID      sql.NullString
		)
		if err := rows.Scan(&userID, &telegramID, &maxID); err != nil {
			slog.Error("checkVerifiedNoPets: scan failed", "err", err)
			continue
		}

		sent, err := n.cache.WasSent(ctx, ports.NotifUserVerifiedNoPets, userID)
		if err != nil {
			slog.Error("checkVerifiedNoPets: cache check failed", "userID", userID, "err", err)
			continue
		}
		if sent {
			continue
		}

		err = n.publisher.PublishNotification(ctx, ports.Notification{
			Type: ports.NotifUserVerifiedNoPets,
			Targets: ports.NotifTargets{
				TelegramID: telegramID.String,
				MaxID:      maxID.String,
			},
			Payload:   map[string]any{},
			CreatedAt: time.Now(),
		})
		if err != nil {
			slog.Error("checkVerifiedNoPets: publish failed", "userID", userID, "err", err)
			continue
		}

		if err := n.cache.MarkSent(ctx, ports.NotifUserVerifiedNoPets, userID, 48*time.Hour); err != nil {
			slog.Error("checkVerifiedNoPets: mark sent failed", "userID", userID, "err", err)
		}
	}
	if err := rows.Err(); err != nil {
		slog.Error("checkVerifiedNoPets: rows iteration error", "err", err)
	}
}

// checkNotVerified — пользователи, не подтвердившие телефон.
// Первое уведомление через 24ч после регистрации (created_at — первый /start в боте),
// повтор каждые 48ч (дедупликация — кеш с TTL 48ч), пока не подтвердят телефон.
// Лимита повторов нет.
func (n *NotificationJob) checkNotVerified(ctx context.Context) {
	rows, err := n.db.QueryContext(ctx, queryNotVerified)
	if err != nil {
		slog.Error("checkNotVerified: query failed", "err", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			userID     string
			telegramID sql.NullString
			maxID      sql.NullString
		)
		if err := rows.Scan(&userID, &telegramID, &maxID); err != nil {
			slog.Error("checkNotVerified: scan failed", "err", err)
			continue
		}

		sent, err := n.cache.WasSent(ctx, ports.NotifUserNotVerified, userID)
		if err != nil {
			slog.Error("checkNotVerified: cache check failed", "userID", userID, "err", err)
			continue
		}
		if sent {
			continue
		}

		err = n.publisher.PublishNotification(ctx, ports.Notification{
			Type: ports.NotifUserNotVerified,
			Targets: ports.NotifTargets{
				TelegramID: telegramID.String,
				MaxID:      maxID.String,
			},
			Payload:   map[string]any{},
			CreatedAt: time.Now(),
		})
		if err != nil {
			slog.Error("checkNotVerified: publish failed", "userID", userID, "err", err)
			continue
		}

		if err := n.cache.MarkSent(ctx, ports.NotifUserNotVerified, userID, 48*time.Hour); err != nil {
			slog.Error("checkNotVerified: mark sent failed", "userID", userID, "err", err)
		}
	}
	if err := rows.Err(); err != nil {
		slog.Error("checkNotVerified: rows iteration error", "err", err)
	}
}
