package query

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/user"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/user/model"
)

type GetDeletedUsersHandler struct {
	userRepo user.Repository
}

func NewGetDeletedUsersHandler(userRepo user.Repository) *GetDeletedUsersHandler {
	return &GetDeletedUsersHandler{
		userRepo: userRepo,
	}
}

func (h *GetDeletedUsersHandler) Handle(ctx context.Context) ([]*model.User, error) {
	users, err := h.userRepo.GetDeletedUsers(ctx)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get deleted users")
	}
	return users, nil
}
