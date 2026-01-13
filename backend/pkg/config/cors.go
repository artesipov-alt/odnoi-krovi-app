package config

import (
	"log/slog"
	"net/url"

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
			allowedOrigins = []string{miniappDomain, "http://localhost:5173"}
			slog.Info("CORS production mode", "origin", miniappDomain)
			slog.Warn("CORS temporarily allowing http://localhost:5173 for debugging! REMOVE THIS BEFORE PRODUCTION DEPLOYMENT.")
		} else {
			allowedOrigins = []string{"http://localhost:5173"} // Temporarily allow localhost
			slog.Error("MINIAPP_DOMAIN is empty! API might be inaccessible. Temporarily allowing http://localhost:5173.")
			slog.Warn("CORS temporarily allowing http://localhost:5173 for debugging due to empty MINIAPP_DOMAIN! REMOVE THIS BEFORE PRODUCTION DEPLOYMENT.")
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
func hasScheme(rawURL string) bool {
	u, err := url.Parse(rawURL)
	return err == nil && u.Scheme != ""
}
