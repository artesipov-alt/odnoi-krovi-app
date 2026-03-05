package cmd

import (
	"context"
	"errors"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
	usermodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/presistance"
)

type UpdateHandler struct {
	userRepo  user.Repository
	txManager *presistance.TxManager
}

func NewUpdateHandler(userRepo user.Repository, txManager *presistance.TxManager) *UpdateHandler {
	return &UpdateHandler{
		userRepo:  userRepo,
		txManager: txManager,
	}
}

func (h *UpdateHandler) Handle(ctx context.Context, id string, input *usermodel.User) (*time.Time, error) {
	now := time.Now()

	err := h.txManager.WithTx(ctx, func(txCtx context.Context) error {
		// Если обновляется телефон, проверяем конфликт
		if input.Phone != "" {
			existingID, err := h.userRepo.GetByPhone(txCtx, input.Phone)
			if err != nil && !errors.Is(err, apperrors.ErrUserNotFound) {
				return apperrors.Internal(err, "failed to check phone uniqueness")
			}

			// Если телефон занят другим пользователем - объединяем аккаунты
			if existingID != "" && existingID != id {
				// 1. Переносим UserIdentity к существующему пользователю
				if err := h.userRepo.TransferUserIdentity(txCtx, id, existingID); err != nil {
					return apperrors.Internal(err, "failed to transfer user identity")
				}

				// 2. Переносим UTM-историю к существующему пользователю (сохраняем аналитику)
				if err := h.userRepo.TransferUTMHistory(txCtx, id, existingID); err != nil {
					return apperrors.Internal(err, "failed to transfer UTM history")
				}

				// 3. Обновляем существующего пользователя данными из input
				if err := h.userRepo.UpdateUserFields(txCtx, existingID, input); err != nil {
					return apperrors.Internal(err, "failed to update existing user")
				}

				// 4. Удаляем настройки донора у текущего пользователя (перед удалением)
				if err := h.userRepo.DeleteDonorPreferenceByUserID(txCtx, id); err != nil {
					return apperrors.Internal(err, "failed to delete donor preference")
				}

				// 5. Удаляем текущего пользователя (полное удаление)
				if err := h.userRepo.DeleteUserHard(txCtx, id); err != nil {
					return apperrors.Internal(err, "failed to delete current user")
				}

				// Аккаунты объединены
				return nil
			}
		}

		// Обычное обновление (без конфликта телефона)
		if err := h.userRepo.UpdateUserFields(txCtx, id, input); err != nil {
			if ent.IsNotFound(err) {
				return apperrors.ErrUserNotFound
			}
			return apperrors.Internal(err, "failed to update user")
		}

		// Обновляем/создаем DonorPreference если передан
		if input.DonorPreference != nil {
			if err := h.userRepo.UpsertDonorPreference(txCtx, id, input.DonorPreference); err != nil {
				return apperrors.Internal(err, "failed to upsert donor preference")
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &now, nil
}
