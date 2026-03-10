package auth

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/partner"
)

type ApiKeysValidator struct {
	partnerRepo partner.Repository
}

func NewApiKeysValidator(partnerRepo partner.Repository) *ApiKeysValidator {
	return &ApiKeysValidator{partnerRepo: partnerRepo}
}

// ValidateBySecret валидирует партнерский API ключ и возвращает информацию о партнере.
func (a ApiKeysValidator) ValidateBySecret(ctx context.Context, id string, apikey string) (providerID string, providerName string, role string) {
	partner, err := a.partnerRepo.GetByAPIKey(ctx, apikey)
	if err != nil {
		return "", "", ""
	}
	return partner.ID, partner.ProviderName, partner.Role
}
