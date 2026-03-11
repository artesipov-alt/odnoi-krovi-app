package user

import (
	"context"
	"time"

	authmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/auth/model"
	"github.com/artesipov-alt/odnoi-krovi-app/pkg/auth"
)

type MiniAppValidator interface {
	ValidateWebAppInitData(ctx context.Context, initData string, providerName authmodel.ProviderName) (*auth.WebAppInitData, error)
}

type TokenGenerator interface {
	Generate(entityID, role string, now time.Time) (token string, expiresAt time.Time)
	Validate(token string) (entityID string, role string, err error)
}
