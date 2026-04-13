package seeds

import (
	"context"
	"log/slog"
	"strconv"
	"strings"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/bloodcomponent"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/breed"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/location"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/schema"
)

// SeedBloodComponents заполняет таблицу компонентов крови начальными данными через ENT
func SeedBloodComponents(ctx context.Context, client *ent.Client) error {
	// Получаем все существующие ID для вычисления следующего номера
	existingIDs, err := client.BloodComponent.Query().Select(bloodcomponent.FieldID).All(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Ошибка при получении существующих ID компонентов крови", "error", err)
		return err
	}
	maxNum := 0
	for _, bc := range existingIDs {
		parts := strings.Split(bc.ID, "-")
		if len(parts) == 2 {
			if num, err := strconv.Atoi(parts[1]); err == nil && num > maxNum {
				maxNum = num
			}
		}
	}
	counter := maxNum + 1

	for _, name := range AllBloodComponents {
		// Проверяем, существует ли уже такой компонент крови
		exists, err := client.BloodComponent.Query().
			Where(bloodcomponent.NameEQ(name)).
			Exist(ctx)

		if err != nil {
			slog.ErrorContext(ctx, "Ошибка при проверке существования компонента крови", "error", err)
			return err
		}

		if !exists {
			// Если не существует, создаем новую запись
			err := client.BloodComponent.Create().
				SetID(schema.BloodComponentPrefix + "-" + strconv.Itoa(counter)).
				SetName(name).
				Exec(ctx)
			counter++

			if err != nil {
				slog.ErrorContext(ctx, "Ошибка при создании компонента крови",
					"name", name,
					"error", err,
				)
				return err
			}
			slog.InfoContext(ctx, "Компонент крови добавлен",
				"name", name,
			)
		} else {
			slog.DebugContext(ctx, "Компонент крови уже существует",
				"name", name,
			)
		}
	}

	slog.InfoContext(ctx, "Заполнение таблицы компонентов крови завершено")
	return nil
}

// SeedLocations заполняет таблицу локаций начальными данными через ENT
func SeedLocations(ctx context.Context, client *ent.Client) error {
	for _, l := range AllLocations {
		exists, err := client.Location.Query().
			Where(location.NameEQ(l.Name)).
			Exist(ctx)

		if err != nil {
			slog.ErrorContext(ctx, "Ошибка при проверке существования локации", "error", err)
			return err
		}

		if !exists {
			err := client.Location.Create().
				SetID(l.ID).
				SetName(l.Name).
				Exec(ctx)

			if err != nil {
				slog.ErrorContext(ctx, "Ошибка при создании локации",
					"name", l.Name,
					"error", err,
				)
				return err
			}
			slog.InfoContext(ctx, "Локация добавлена", "name", l.Name)
		}
	}

	slog.InfoContext(ctx, "Заполнение таблицы локаций завершено")
	return nil
}

// SeedBreeds заполняет таблицу пород начальными данными через ENT
func SeedBreeds(ctx context.Context, client *ent.Client) error {
	for _, b := range Breeds {
		exists, err := client.Breed.Query().
			Where(breed.IDEQ(b.ID)).
			Exist(ctx)

		if err != nil {
			slog.ErrorContext(ctx, "Ошибка при проверке существования породы", "error", err)
			return err
		}

		if !exists {
			err := client.Breed.Create().
				SetID(b.ID).
				SetName(b.Name).
				SetType(b.Type).
				Exec(ctx)

			if err != nil {
				slog.ErrorContext(ctx, "Ошибка при создании породы",
					"id", b.ID,
					"name", b.Name,
					"error", err,
				)
				return err
			}
			slog.InfoContext(ctx, "Порода добавлена", "name", b.Name)
		}
	}

	slog.InfoContext(ctx, "Заполнение таблицы пород завершено")
	return nil
}
