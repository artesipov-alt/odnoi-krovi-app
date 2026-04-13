package seeds

import (
	"context"
	"log/slog"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/breed"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/location"
)

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
