package auth

import (
	"errors"
	"strings"
)

type JWTGenerator struct {
	secretKey string
}

func NewJWTGenerator(secretKey string) *JWTGenerator {
	return &JWTGenerator{secretKey: secretKey}
}

func (g *JWTGenerator) Generate(entityID, role string) (token string) {
	if role == "admin" {
		return "admin_token_" + entityID
	} else {
		return "user_token_" + entityID
	}
}

func (g *JWTGenerator) Validate(token string) (entityID string, role string, err error) {
	if after, ok := strings.CutPrefix(token, "admin_token_"); ok {
		entityID = after
		role = "admin"
		err = nil
	} else if after, ok := strings.CutPrefix(token, "user_token_"); ok {
		entityID = after
		role = "user"
		err = nil
	} else {
		entityID = ""
		role = ""
		err = errors.New("invalid token")
	}
	return entityID, role, err
}
