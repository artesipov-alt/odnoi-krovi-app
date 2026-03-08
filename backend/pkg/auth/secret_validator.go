package auth

import (
	"context"
	"strings"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/partner"
)

type AppValidator struct {
	botToken    string
	partnerRepo partner.Repository
}

func NewAppValidator(botToken string, partnerRepo partner.Repository) *AppValidator {
	return &AppValidator{
		botToken:    botToken,
		partnerRepo: partnerRepo,
	}
}

func (v *AppValidator) ValidateHash(initData string) (providerID string, role string) {
	splited := strings.Split(initData, "-")
	if len(splited) < 2 {
		return "", ""
	}
	if splited[1] != v.botToken {
		return "", ""
	}
	providerID = splited[0]
	role = "USER"
	return providerID, role
}

func (v *AppValidator) ValidateBySecret(ctx context.Context, id string, apikey string) (providerID string, providerName string, role string) {
	partner, err := v.partnerRepo.GetByAPIKey(ctx, apikey)
	if err != nil {
		return "", "", ""
	}
	providerID = id
	if strings.HasSuffix(partner.ID, "_bot") {
		providerName = partner.ID
	} else {
		providerName = "service"
	}
	return providerID, providerName, partner.Role
}
