package config

import "github.com/gofiber/fiber/v2/middleware/cors"

// CORSOptions возвращает конфигурацию CORS оптимизированную для Telegram MiniApp
func CORSOptions() cors.Config {
	return cors.Config{
		// Разрешаем все источники для Telegram MiniApp, так как он может быть встроен в различные домены
		AllowOrigins: "*",

		// Разрешаем методы, обычно используемые в REST API
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",

		// Разрешаем заголовки, необходимые для Telegram MiniApp и API запросов
		AllowHeaders: "Origin,Content-Type,Accept,Authorization,X-Requested-With,tg-init-data",

		// Отображаем заголовки, которые могут понадобиться фронтенду
		ExposeHeaders: "Content-Length,Content-Range",

		// Разрешаем учетные данные, если нужна аутентификация
		//AllowCredentials: true,

		// Кэшируем предварительные запросы на 12 часов
		MaxAge: 43200,
	}
}
