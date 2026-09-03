// Command migrate применяет one-off SQL-миграции из migrations/one-off.
//
// Ведёт учёт применённых миграций в таблице applied_migrations: применяются
// только файлы, отсутствующие в ней. Каждый файл выполняется целиком одним
// Exec — файл сам управляет транзакцией (BEGIN/COMMIT, см. migrations/README.md);
// если BEGIN/COMMIT внутри файла нет, PostgreSQL выполняет несколько стейтментов
// одного Exec как неявную транзакцию (всё или ничего).
//
// Использование (из каталога backend/):
//
//	go run ./cmd/migrate [env]            применить pending-миграции
//	go run ./cmd/migrate -baseline [env]  пометить все как применённые, не выполняя
//
// env: local (по умолчанию), dev, prod — выбирает целевую БД через config.NewEntConfig.
// -baseline нужен для свежей БД: one-off миграции там не нужны, их надо только
// отметить (обычно делается автоматически задачей db:migrate:bootstrap).
package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/artesipov-alt/odnoi-krovi-app/pkg/config"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const migrationsDir = "migrations/one-off"

func main() {
	baseline := flag.Bool("baseline", false, "пометить все миграции применёнными без выполнения (для свежей БД)")
	flag.Parse()

	env := "local"
	if args := flag.Args(); len(args) > 0 {
		env = args[0]
	}

	// .env лежит в корне проекта (как в cmd/dburl). godotenv.Load молча
	// пропускает отсутствие файла — в Docker переменные приходят через окружение.
	_ = godotenv.Load("../.env")

	db, err := sql.Open("postgres", config.NewEntConfig(env).GetDSN())
	if err != nil {
		log.Fatalf("не удалось открыть подключение к БД (env=%s): %v", env, err)
	}
	defer db.Close()

	if err := ensureTrackingTable(db); err != nil {
		log.Fatalf("не удалось создать таблицу applied_migrations: %v", err)
	}

	applied, err := loadApplied(db)
	if err != nil {
		log.Fatalf("не удалось прочитать applied_migrations: %v", err)
	}

	files, err := listMigrations()
	if err != nil {
		log.Fatal(err)
	}

	var pending []string
	for _, f := range files {
		if !applied[f] {
			pending = append(pending, f)
		}
	}

	if len(pending) == 0 {
		fmt.Println("one-off миграций к применению нет (все применены)")
		return
	}

	if *baseline {
		for _, f := range pending {
			if err := markApplied(db, f); err != nil {
				log.Fatalf("%s: не удалось записать в applied_migrations: %v", f, err)
			}
			fmt.Printf("baseline: %s\n", f)
		}
		fmt.Printf("Помечено применёнными (без выполнения): %d\n", len(pending))
		return
	}

	for _, f := range pending {
		if err := applyOne(db, f); err != nil {
			log.Fatalf("%s: %v\nОставшиеся one-off миграции НЕ применены — устраните ошибку и повторите запуск.", f, err)
		}
	}
	fmt.Printf("Применено миграций: %d\n", len(pending))
}

// ensureTrackingTable создаёт таблицу учёта применённых one-off миграций.
func ensureTrackingTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS applied_migrations (
			name       text PRIMARY KEY,
			applied_at timestamptz NOT NULL DEFAULT now()
		)
	`)
	return err
}

// loadApplied возвращает множество имён уже применённых миграций.
func loadApplied(db *sql.DB) (map[string]bool, error) {
	rows, err := db.Query(`SELECT name FROM applied_migrations`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		applied[name] = true
	}
	return applied, rows.Err()
}

// listMigrations возвращает отсортированный список SQL-файлов в migrationsDir.
func listMigrations() ([]string, error) {
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать каталог %s: %w", migrationsDir, err)
	}

	var files []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".sql") {
			continue
		}
		files = append(files, e.Name())
	}
	sort.Strings(files)

	if len(files) == 0 {
		return nil, errors.New("в каталоге " + migrationsDir + " нет SQL-файлов")
	}
	return files, nil
}

// applyOne выполняет файл миграции целиком и отмечает его применённым.
// Отметка отдельным Exec после успеха: при ошибке в файле он останется
// pending и будет применён при следующем запуске.
func applyOne(db *sql.DB, name string) error {
	content, err := os.ReadFile(filepath.Join(migrationsDir, name))
	if err != nil {
		return err
	}

	fmt.Printf("Применяю %s ... ", name)
	if _, err := db.Exec(string(content)); err != nil {
		fmt.Println("ОШИБКА")
		return err
	}
	if err := markApplied(db, name); err != nil {
		fmt.Println("ОШИБКА")
		return err
	}
	fmt.Println("OK")
	return nil
}

func markApplied(db *sql.DB, name string) error {
	_, err := db.Exec(`INSERT INTO applied_migrations (name) VALUES ($1) ON CONFLICT (name) DO NOTHING`, name)
	return err
}
