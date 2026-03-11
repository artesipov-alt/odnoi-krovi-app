package user

import (
	"context"
	"time"

	authmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/auth/model"
	"github.com/artesipov-alt/odnoi-krovi-app/pkg/auth"
)

type MiniAppValidator interface {
	ValidateMock(initData string) (providerID string, role string)
	ValidateWebAppInitData(ctx context.Context, initData string, providerName authmodel.ProviderName) (*auth.WebAppInitData, error)
}

type ApiKeysValidator interface {
	ValidateBySecret(ctx context.Context, id, providerName, apikey string) (providerID string, role string)
}

type TokenGenerator interface {
	Generate(entityID, role string, now time.Time) (token string, expiresAt time.Time)
	Validate(token string) (entityID string, role string, err error)
}
