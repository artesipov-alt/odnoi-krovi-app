package auth

import (
	"log/slog"
	"strconv"
	"strings"
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

func (v *AppValidator) ValidateBySecret(id int64, secret string) (providerID int64) {
	slog.Info("ValidateBySecret", "id", id, "secret", secret, "default-secret", v.secret)
	if secret != v.secret {
		return 0
	}
	return id
}
