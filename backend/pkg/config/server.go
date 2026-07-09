package config

import (
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
)

type MyServer struct {
	*http.Server
}

// NewServer создает новый экземпляр сервера с заданным mux
func NewServer(port int, mux http.Handler) *MyServer {
	return &MyServer{
		&http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		},
	}
}

// Use добавляет middleware в цепочку обработки.
// Принимает стандартные функции-обертки func(http.Handler) http.Handler.
// Middleware применяются так, что первое переданное в списке становится самым внешним слоем.
// Это позволяет соблюдать логический порядок: Recovery -> RequestID -> Logging -> Mux.
func (s *MyServer) Use(middlewares ...func(http.Handler) http.Handler) {
	for i := len(middlewares) - 1; i >= 0; i-- {
		s.Handler = middlewares[i](s.Handler)
	}
}

// NewHumaConfig создает и возвращает конфигурацию Huma API на основе README.md
func NewHumaConfig(miniappDomain string) huma.Config {
	config := huma.DefaultConfig("Одной Крови API", "3.20.4")

	config.Info = &huma.Info{
		Title:       "Одной Крови API",
		Version:     "3.20.4", // Или динамически брать из переменной окружения/сборки
		Description: "### Описание платформы\n**Одной Крови** — это Telegram Mini App, который помогает находить донорскую кровь для животных и позволяет владельцам питомцев становиться донорами вместе со своими любимцами.\n\n* **Поиск доноров**: Быстрый поиск доноров крови для животных в экстренных ситуациях.\n* **Регистрация доноров**: Возможность регистрации питомцев как потенциальных доноров.\n* **Геолокация**: Определение ближайших доноров через Telegram Web App.\n* **Уведомления**: Система оповещений через Telegram Bot API.\n* **Интеграция с Telegram**: Удобное общение между пользователями через Telegram.",
		Contact: &huma.Contact{
			Name:  "Команда Одной Крови",
			Email: "support@odnoi-krovi.app", // Пример
			URL:   "https://1krovi.app",      // Пример
		},
		License: &huma.License{
			Name: "Apache 2.0 License",
			URL:  "https://www.apache.org/licenses/LICENSE-2.0",
		},
	}

	config.ExternalDocs = &huma.ExternalDocs{
		Description: "Журнал изменений",
		URL:         "https://github.com/artesipov-alt/odnoi-krovi-app/blob/dev/backend/CHANGELOG.md",
	}

	// Добавляем серверы, включая локальный и MiniApp домен
	config.Servers = []*huma.Server{
		{
			URL:         "https://dev.1krovi.app/api",
			Description: "Development API для Telegram Mini App",
		},
		{
			URL:         "https://1krovi.app/api",
			Description: "Production API для Telegram Mini App",
		},
		{
			URL:         "http://localhost:3001/api",
			Description: "Локальная разработка API"},
	}

	// Явно указываем пути для OpenAPI спецификации и UI документации
	config.OpenAPIPath = "/openapi.json"
	config.DocsPath = "" // Отключаем встроенный UI Huma, чтобы использовать свой SwaggerUIHandler

	return config
}
