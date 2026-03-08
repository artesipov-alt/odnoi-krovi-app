package auth

import (
	"log/slog"
	"strconv"
	"strings"

	authmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/auth/model"
)

type AppValidator struct {
	botToken string
	secret   string
}

func NewAppValidator(botToken string, secret string) *AppValidator {
	return &AppValidator{
		botToken: botToken,
		secret:   secret,
	}
}

func (v *AppValidator) ValidateHash(initData string) (providerID int64) {
	splited := strings.Split(initData, "-")
	if len(splited) < 2 {
		return 0
	}
	if splited[1] != v.botToken {
		return 0
	}
	providerID, err := strconv.ParseInt(splited[0], 10, 64)
	if err != nil {
		return 0
	}
	return providerID
}

var apiKeys = map[string]string{
	"max_bot":      "API-MAX",
	"telegram_bot": "API-TELEGRAM",
}

func (v *AppValidator) ValidateBySecret(id int64, secret string) (providerID int64, providerName authmodel.ProviderName) {
	slog.Info("ValidateBySecret", "id", id, "secret", secret, "default-secret", v.secret)
	switch secret {
	case apiKeys["max_bot"]:
		return id, authmodel.ProviderMax
	case apiKeys["telegram_bot"]:
		return id, authmodel.ProviderTelegram
	}
	return 0, ""
}
