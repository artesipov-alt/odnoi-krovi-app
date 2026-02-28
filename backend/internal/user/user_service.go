package user

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/user/model"
)

type UserService struct {
	userRepo Repository
}

// NewUserService создает новый экземпляр UserService
func NewUserService(userRepo Repository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// RegisterUser регистрирует нового пользователя в системе
func (s *UserService) RegisterUser(ctx context.Context, user *model.User) (*model.User, error) {
	newUser, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to create user")
	}
	return newUser, nil
}

// DeleteUser удаляет пользователя по ID (soft delete)
func (s *UserService) DeleteUser(ctx context.Context, userID string) error {

	// Удаляем пользователя
	if err := s.userRepo.Delete(ctx, userID); err != nil {
		return apperrors.Internal(err, "failed to delete user")
	}

	return nil
}

func (s *UserService) Update(ctx context.Context, id string, input *model.User) error {

	err := s.userRepo.Update(ctx, id, input)
	if err != nil {
		// Используем ent.IsNotFound - это идиоматический способ для Ent
		if ent.IsNotFound(err) {
			return apperrors.ErrUserNotFound
		}
		return apperrors.Internal(err, "failed to update user")
	}

	return nil
}

// ResetUser сбрасывает пользователя к начальным настройкам
func (s *UserService) ResetUser(ctx context.Context, userID string) error {
	if err := s.userRepo.ResetUser(ctx, userID); err != nil {
		return apperrors.Internal(err, "failed to reset user")
	}
	return nil
}

// RestoreUser восстанавливает удаленного пользователя
func (s *UserService) RestoreUser(ctx context.Context, userID string) error {
	if err := s.userRepo.RestoreUser(ctx, userID); err != nil {
		return apperrors.Internal(err, "failed to restore user")
	}
	return nil
}

// GetDeletedUsers получает всех удаленных пользователей
func (s *UserService) GetDeletedUsers(ctx context.Context) ([]*model.User, error) {
	users, err := s.userRepo.GetDeletedUsers(ctx)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get deleted users")
	}
	return users, nil
}
