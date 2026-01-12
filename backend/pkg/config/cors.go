package config

import (
	"log/slog"

	"github.com/rs/cors"
)

func SetupCORS(env, miniappDomain string) *cors.Cors {
	var allowedOrigins []string

	if env == "development" {
		// ВАЖНО: При AllowCredentials: true нельзя использовать "*"
		// Лучше указать конкретные локальные порты твоих фронтенд-инструментов (Vite/Next)
		allowedOrigins = []string{"http://localhost:5173", "http://localhost:3000", "http://127.0.0.1:5173"}
		slog.Warn("CORS development mode: allowed local origins")
	} else {
		if miniappDomain != "" {
			if !hasScheme(miniappDomain) {
				miniappDomain = "https://" + miniappDomain
			}
			allowedOrigins = []string{miniappDomain}
			slog.Info("CORS production mode", "origin", miniappDomain)
		} else {
			allowedOrigins = []string{}
			slog.Error("MINIAPP_DOMAIN is empty! API might be inaccessible.")
		}
	}

	return cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,                 // Позволяет передавать куки/авторизацию
		Debug:            env == "development", // Включает подробные логи CORS в консоль
	})
}

// hasScheme проверяет, содержит ли URL схему (http/https).
func hasScheme(url string) bool {
	return len(url) >= 7 && (url[:7] == "http://" || url[:8] == "https://")
}
