package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTGenerator генерирует и валидирует JWT токены.
type JWTGenerator struct {
	secretKey []byte
	issuer    string
	tokenTTL  time.Duration
}

// NewJWTGenerator создает новый генератор JWT токенов.
// secretKey - секретный ключ для подписи (рекомендуется 256+ бит для HS256)
// issuer - издатель токена (например, "odnoi-krovi-app")
// tokenTTL - время жизни токена (рекомендуется 24 часа)
func NewJWTGenerator(secretKey string, issuer string, tokenTTL time.Duration) *JWTGenerator {
	return &JWTGenerator{
		secretKey: []byte(secretKey),
		issuer:    issuer,
		tokenTTL:  tokenTTL,
	}
}

// Claims представляет собой кастомные claims для JWT токена.
type Claims struct {
	jwt.RegisteredClaims
	Role string
}

// Generate генерирует JWT токен для заданной сущности и роли.
// Возвращает подписанный JWT токен, время истечения или пустые значения в случае ошибки.
func (g *JWTGenerator) Generate(entityID, role string, now time.Time) (string, time.Time) {
	expiresAt := now.Add(g.tokenTTL)

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   entityID,
			Issuer:    g.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
		Role: role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(g.secretKey)
	if err != nil {
		// В продакшене здесь должен быть логирован error
		// log.Printf("failed to sign JWT token: %v", err)
		return "", time.Time{}
	}

	return signedToken, expiresAt
}

// Validate валидирует JWT токен и возвращает entityID, роль и ошибку.
// Проверяет:
//   - подпись токена
//   - срок действия (exp)
//   - издателя (iss)
func (g *JWTGenerator) Validate(tokenString string) (entityID string, role string, err error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Проверяем, что используется ожидаемый алгоритм подписи
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return g.secretKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", "", errors.New("token has expired")
		}
		return "", "", err
	}

	if !token.Valid {
		return "", "", errors.New("invalid token")
	}

	// Дополнительная проверка издателя (опционально, но рекомендуется)
	if claims.Issuer != g.issuer {
		return "", "", errors.New("invalid token issuer")
	}

	return claims.Subject, claims.Role, nil
}

// ValidateToken — альтернативный метод, возвращающий только ошибку валидации.
// Может быть полезен, если не нужны данные из токена.
func (g *JWTGenerator) ValidateToken(tokenString string) error {
	_, _, err := g.Validate(tokenString)
	return err
}
