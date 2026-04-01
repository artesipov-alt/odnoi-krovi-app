package query

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	authmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/auth/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/ports"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user"
	userevent "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user/events"
	usermodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/user/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/transport/http/middleware"
)

type GetContactHandler struct {
	userRepo  user.Repository
	publisher ports.EventPublisher
}

func NewGetContactHandler(userepo user.Repository, publisher ports.EventPublisher) *GetContactHandler {
	return &GetContactHandler{
		userRepo:  userepo,
		publisher: publisher,
	}
}

func (h *GetContactHandler) Handle(ctx context.Context, id string, provider string) (*usermodel.User, error) {
	userID := middleware.GetUserID(ctx)
	if userID == "" {
		return nil, apperrors.Unauthorized("user ID is missing in context")
	}

	userInitiator, err := h.userRepo.GetByID(ctx, userID, user.UserPreloadOptions{WithIdentities: true})
	if err != nil {
		return nil, err
	}

	var sendToID string
	for _, identity := range userInitiator.Identities {
		if identity.ProviderName == authmodel.ProviderName(provider) {
			sendToID = identity.ProviderUserID
			break
		}
	}

	userData, err := h.userRepo.GetByID(ctx, id, user.UserPreloadOptions{
		WithIdentities: true,
	})
	if err != nil {
		return nil, err // Доменная ошибка (например, ErrUserNotFound)
	}

	event := userevent.GenerateContact(userData)
	event.SendTo = sendToID
	event.NotifyProvider = authmodel.ProviderName(provider)

	if err := h.publisher.PublishUserContact(ctx, event); err != nil {
		return nil, err
	}

	return userData, nil
}
