package middleware

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"slices"
	"strings"

	"github.com/artesipov-alt/odnoi-krovi-app/pkg/auth"
)

// AuthContextKey — приватный тип для ключей контекста авторизации.
type AuthContextKey string

// Константы для ключей контекста авторизации.
const (
	// UserIDKey — ключ для хранения ID пользователя в контексте.
	UserIDKey AuthContextKey = "user_id"
	// UserRoleKey — ключ для хранения роли пользователя в контексте.
	UserRoleKey AuthContextKey = "user_role"
)

// AuthMiddleware создает middleware для аутентификации JWT токенов.
// JWTGenerator используется для валидации токена и извлечения данных пользователя.
// excludedPaths — список путей, которые не требуют аутентификации (например, "/auth/signin").
// Если токен отсутствует или невалиден, возвращается 401 ошибка.
func AuthMiddleware(jwtGenerator *auth.JWTGenerator, excludedPaths ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Проверяем, нужно ли пропустить аутентификацию для этого пути
			path := r.URL.Path
			for _, excluded := range excludedPaths {
				if strings.HasPrefix(path, excluded) {
					next.ServeHTTP(w, r)
					return
				}
			}

			// Извлекаем токен из заголовка Authorization
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				slog.WarnContext(r.Context(), "missing authorization header",
					"path", r.URL.Path,
				)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// Проверяем формат "Bearer <token>"
			const bearerPrefix = "Bearer "
			if !strings.HasPrefix(authHeader, bearerPrefix) {
				slog.WarnContext(r.Context(), "invalid authorization header format",
					"path", r.URL.Path,
				)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, bearerPrefix)
			if tokenString == "" {
				slog.WarnContext(r.Context(), "empty token in authorization header",
					"path", r.URL.Path,
				)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// Валидируем токен и извлекаем данные
			userID, role, err := jwtGenerator.Validate(tokenString)
			if err != nil {
				slog.WarnContext(r.Context(), "invalid JWT token",
					"path", r.URL.Path,
					"error", err.Error(),
				)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// Проверяем, что userID не пустой
			if userID == "" {
				slog.WarnContext(r.Context(), "empty user ID in token",
					"path", r.URL.Path,
				)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// Создаем новый контекст с данными пользователя
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			ctx = context.WithValue(ctx, UserRoleKey, role)

			slog.DebugContext(ctx, "user authenticated",
				"user_id", userID,
				"role", role,
			)

			// Передаем управление следующему обработчику с обновленным контекстом
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID возвращает ID пользователя из контекста запроса.
// Если данные отсутствуют, возвращает пустую строку.
func GetUserID(ctx context.Context) string {
	if userID, ok := ctx.Value(UserIDKey).(string); ok {
		return userID
	}
	return ""
}

// GetUserRole возвращает роль пользователя из контекста запроса.
// Если данные отсутствуют, возвращает пустую строку.
func GetUserRole(ctx context.Context) string {
	if role, ok := ctx.Value(UserRoleKey).(string); ok {
		return role
	}
	return ""
}

// GetUserIDWithError возвращает ID пользователя из контекста или ошибку, если данные отсутствуют.
// Это удобно для использования в обработчиках, где требуется гарантированное наличие userID.
func GetUserIDWithError(ctx context.Context) (string, error) {
	userID := GetUserID(ctx)
	if userID == "" {
		return "", errors.New("user ID not found in context: authorization required")
	}
	return userID, nil
}

// RequireRole создает middleware, который проверяет наличие у пользователя указанной роли.
// Если роль не совпадает, возвращается 403 Forbidden.
// Комбинируйте с AuthMiddleware: AuthMiddleware(generator, RequireRole("admin"))
func RequireRole(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole := GetUserRole(r.Context())

			if userRole != requiredRole {
				slog.WarnContext(r.Context(), "insufficient permissions",
					"path", r.URL.Path,
					"user_role", userRole,
					"required_role", requiredRole,
				)
				w.WriteHeader(http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireAnyRole создает middleware, который проверяет наличие у пользователя одной из указанных ролей.
// Если роль не совпадает ни с одной из требуемых, возвращается 403 Forbidden.
func RequireAnyRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole := GetUserRole(r.Context())

			if slices.Contains(roles, userRole) {
				next.ServeHTTP(w, r)
				return
			}

			slog.WarnContext(r.Context(), "insufficient permissions",
				"path", r.URL.Path,
				"user_role", userRole,
				"required_roles", roles,
			)
			w.WriteHeader(http.StatusForbidden)
		})
	}
}
