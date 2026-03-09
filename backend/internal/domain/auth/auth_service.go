package user

import (
	"context"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/pkg/auth"
)

type AppValidator interface {
	ValidateMock(initData string) (providerID string, role string)
	ValidateWebAppInitData(ctx context.Context, initData string) (*auth.WebAppInitData, error)
}

type ApiKeysValidator interface {
	ValidateBySecret(ctx context.Context, id string, apikey string) (providerID string, providerName string, role string)
}

type TokenGenerator interface {
	Generate(entityID, role string, now time.Time) (token string, expiresAt time.Time)
	Validate(token string) (entityID string, role string, err error)
}
