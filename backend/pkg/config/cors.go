package config

import (
	"log/slog"
	"net/http"

	"github.com/rs/cors"
)

// SetupCORS настраивает CORS middleware в зависимости от окружения.
func SetupCORS(env, miniappDomain string) *cors.Cors {
	var allowedOrigins []string

	if env == "development" {
		allowedOrigins = []string{"*"} // Разрешаем все источники для разработки
		slog.Warn("CORS is configured to allow all origins in development mode.")
	} else {
		// В продакшене или других окружениях разрешаем только указанный домен
		if miniappDomain != "" {
			// Добавляем схему, если она отсутствует, так как домены обычно указываются с ней
			if !hasScheme(miniappDomain) {
				miniappDomain = "https://" + miniappDomain
			}
			allowedOrigins = []string{miniappDomain}
			slog.Info("CORS is configured to allow specific origin.", "origin", miniappDomain)
		} else {
			// Если домен не указан, можно выбрать политику по умолчанию, например, запретить все
			// или разрешить только текущий домен сервера.
			// Для безопасности лучше явно указать разрешенные домены.
			allowedOrigins = []string{} // По умолчанию запрещаем, если домен не указан
			slog.Error("MINIAPP_DOMAIN is not set for non-development environment. CORS will block requests.")
		}
	}

	return cors.New(cors.Options{
		AllowedOrigins: allowedOrigins,
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
		},
		AllowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"X-CSRF-Token",
			"X-Request-ID",
		},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // 5 минут кеширования для preflight запросов
	})
}

// hasScheme проверяет, содержит ли URL схему (http/https).
func hasScheme(url string) bool {
	return len(url) >= 7 && (url[:7] == "http://" || url[:8] == "https://")
}
