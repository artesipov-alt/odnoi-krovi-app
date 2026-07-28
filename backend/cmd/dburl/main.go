// Command dburl печатает строку подключения к БД для указанного окружения.
//
// Используется в Taskfile для передачи DSN в psql при применении SQL-миграций:
//
// \tpsql "$(go run ./cmd/dburl {{.ENV}})" -f migrations/...
//
// Поддерживаемые значения ENV: local (по умолчанию), dev, prod.
// Для dev/prod переменные окружения (DB_HOST, DB_USER, ...) должны быть заданы.
package main

import (
	"fmt"
	"os"

	"github.com/artesipov-alt/odnoi-krovi-app/pkg/config"
	"github.com/joho/godotenv"
)

func main() {
	// Загружаем .env из корня проекта (как в cmd/api/main.go).
	// В Docker не нужен — переменные приходят через docker-compose environment.
	// godotenv.Load молча пропускает отсутствие файла.
	_ = godotenv.Load("../.env")

	env := "local"
	if len(os.Args) > 1 {
		env = os.Args[1]
	}
	cfg := config.NewEntConfig(env)
	fmt.Print(cfg.GetPsqlURL())
}
