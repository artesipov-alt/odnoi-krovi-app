package main

import (
	"log"

	"github.com/artesipov-alt/odnoi-krovi-app/poolservice/gen/api/greet/v1/greetv1connect"
	v1 "github.com/artesipov-alt/odnoi-krovi-app/poolservice/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
)

func main() {
	app := fiber.New()

	// Создаем gRPC сервер
	greeter := &v1.GreetServer{}
	path, handler := greetv1connect.NewGreetServiceHandler(greeter)

	// Подключаем gRPC handler к Fiber
	app.Post(path+"*", adaptor.HTTPHandler(handler))

	log.Println("Сервер запущен на :8080")
	log.Fatal(app.Listen(":8080"))
}
