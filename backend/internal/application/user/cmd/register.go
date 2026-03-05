package cmd

import (
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
)

type RegisterHandler struct {
	userRepo user.Repository
}

func NewRegisterHandler(userRepo user.Repository) *RegisterHandler {
	return &RegisterHandler{
		userRepo: userRepo,
	}
}

// func (h *RegisterHandler) Handle(ctx context.Context, user *usermodel.User) (*usermodel.User, error) {
// 	newUser, err := h.userRepo.Create(ctx, user)
// 	if err != nil {
// 		return nil, apperrors.Internal(err, "failed to create user")
// 	}
// 	return newUser, nil
// }
