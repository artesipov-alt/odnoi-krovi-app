package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/presistance"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/presistance/redis/otp"
)

type VerifyPhoneHandler struct {
	userRepo  user.Repository
	otpRepo   otp.OTPRepository
	txManager *presistance.TxManager
}

func NewVerifyPhoneHandler(userRepo user.Repository, otpRepo otp.OTPRepository, txManager *presistance.TxManager) *VerifyPhoneHandler {
	return &VerifyPhoneHandler{
		userRepo:  userRepo,
		otpRepo:   otpRepo,
		txManager: txManager,
	}
}

func (h *VerifyPhoneHandler) Handle(ctx context.Context, userID string, code string) error {
	otpData, err := h.otpRepo.Get(ctx, userID)
	if err != nil {
		return apperrors.BadRequest("неверный или просроченный код подтверждения")
	}

	if otpData.Code != code {
		return apperrors.BadRequest("неверный код подтверждения")
	}

	newPhone := otpData.NewPhone

	err = h.txManager.WithTx(ctx, func(txCtx context.Context) error {
		existingID, err := h.userRepo.GetByPhone(txCtx, newPhone)
		if err != nil && !errors.Is(err, apperrors.ErrUserNotFound) {
			return apperrors.Internal(err, "failed to check phone uniqueness")
		}

		// Если телефон занят другим пользователем — объединяем аккаунты
		if existingID != "" && existingID != userID {
			// Загружаем identity обоих пользователей для проверки пересечения провайдеров
			existingUser, err := h.userRepo.GetByID(txCtx, existingID, user.UserPreloadOptions{
				WithIdentities: true,
			})
			if err != nil {
				return apperrors.Internal(err, "failed to get existing user")
			}

			currentUser, err := h.userRepo.GetByID(txCtx, userID, user.UserPreloadOptions{
				WithIdentities: true,
			})
			if err != nil {
				return apperrors.Internal(err, "failed to get current user")
			}

			// Проверяем, нет ли пересечения по провайдерам (защита от чёрного донорства)
			if existingUser.HasProviderConflict(currentUser) {
				return apperrors.Conflict("нельзя объединить аккаунты: у вас уже есть аккаунт в этом сервисе")
			}

			// 1. Переносим UserIdentity к существующему пользователю
			if err := h.userRepo.TransferUserIdentity(txCtx, userID, existingID); err != nil {
				return apperrors.Internal(err, "failed to transfer user identity")
			}

			// 2. Переносим питомцев к существующему пользователю (со всеми заявками и откликами)
			if err := h.userRepo.TransferPets(txCtx, userID, existingID); err != nil {
				return apperrors.Internal(err, "failed to transfer pets")
			}

			// 3. Переносим UTM-историю к существующему пользователю (сохраняем аналитику)
			if err := h.userRepo.TransferUTMHistory(txCtx, userID, existingID); err != nil {
				return apperrors.Internal(err, "failed to transfer UTM history")
			}

			// 4. Обновляем телефон существующего пользователя
			if err := h.userRepo.UpdatePhone(txCtx, existingID, newPhone); err != nil {
				return apperrors.Internal(err, "failed to update phone on existing user")
			}

			// 4. Удаляем настройки донора у текущего пользователя (перед удалением)
			if err := h.userRepo.DeleteDonorPreferenceByUserID(txCtx, userID); err != nil {
				return apperrors.Internal(err, "failed to delete donor preference")
			}

			// 5. Удаляем текущего пользователя (полное удаление)
			if err := h.userRepo.DeleteUserHard(txCtx, userID); err != nil {
				return apperrors.Internal(err, "failed to delete current user")
			}

			return nil
		}

		// Телефон свободен — просто обновляем
		if err := h.userRepo.UpdatePhone(txCtx, userID, newPhone); err != nil {
			return apperrors.Internal(err, "failed to update phone")
		}

		return nil
	})

	if err != nil {
		return err
	}

	if err := h.otpRepo.Delete(ctx, userID); err != nil {
		return fmt.Errorf("verify phone: delete otp: %w", err)
	}

	return nil
}
