package cmd

import (
	"context"

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

func (h *UpdateHandler) Handle(ctx context.Context, id string, input *usermodel.User) (*usermodel.User, error) {
	// 1. Загружаем существующий агрегат
	existingUser, err := h.userRepo.GetByID(ctx, id, user.UserPreloadOptions{WithDonorPreference: true})
	if err != nil {
		return nil, err
	}

	// 2. Применяем изменения через доменный метод агрегата (контролируемая мутация + валидация)
	if err := existingUser.UpdateFrom(input); err != nil {
		return nil, apperrors.Validation("Ошибка валидации", nil).WithInternal(err)
	}

	// 3. Сохраняем агрегат в транзакции
	err = h.txManager.WithTx(ctx, func(txCtx context.Context) error {
		if err := h.userRepo.UpdateUserFields(txCtx, id, existingUser); err != nil {
			if ent.IsConstraintError(err) {
				return apperrors.Conflict("Пользователь с такой почтой уже существует, при подтверждении номера аккаунты будут связаны")
			}
			return apperrors.Internal(err, "failed to update user")
		}

		if existingUser.DonorPreference != nil {
			if err := h.userRepo.UpsertDonorPreference(txCtx, id, existingUser.DonorPreference); err != nil {
				return apperrors.Internal(err, "failed to upsert donor preference")
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 4. Возвращаем обновлённого пользователя (без preload — только базовые поля)
	usr, err := h.userRepo.GetByID(ctx, id, user.UserPreloadOptions{})
	if err != nil {
		return nil, err
	}

	return usr, nil
}
