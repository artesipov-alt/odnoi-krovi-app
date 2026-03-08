package user

import "context"

type AppValidator interface {
	ValidateHash(initData string) (providerID string, role string)
	ValidateBySecret(ctx context.Context, id string, apikey string) (providerID string, providerName string, role string)
}

type TokenGenerator interface {
	Generate(entityID, role string) (token string)
	Validate(token string) (entityID string, role string, err error)
}
