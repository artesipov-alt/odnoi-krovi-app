package auth

import (
	"context"
	"strings"

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
	if strings.HasSuffix(partner.ID, "_bot") {
		providerName = partner.ID
	} else {
		providerName = "service"
	}
	return partner.ID, providerName, partner.Role
}
