package seeds

import (
	"context"
	"log/slog"

	"github.com/artesipov-alt/odnoi-krovi-app/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/bloodcomponent"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/bloodgroup"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/breed"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/location"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/utils/seeds/data"
)

// SeedBloodGroups заполняет таблицу групп крови начальными данными через ENT
func SeedBloodGroups(ctx context.Context, client *ent.Client, log *slog.Logger) error {
	groups := []struct {
		PetType     bloodgroup.PetType
		BloodGroup  string
		Description string
	}{
		// Группы крови для собак
		{
			PetType:     bloodgroup.PetTypeDog,
			BloodGroup:  "DEA 1+",
			Description: "Универсальный донор для собак с положительным DEA 1+",
		},
		{
			PetType:     bloodgroup.PetTypeDog,
			BloodGroup:  "DEA 1-",
			Description: "Универсальный донор для всех собак",
		},
		// Группы крови для кошек
		{
			PetType:     bloodgroup.PetTypeCat,
			BloodGroup:  "A",
			Description: "Самая распространенная группа крови у кошек",
		},
		{
			PetType:     bloodgroup.PetTypeCat,
			BloodGroup:  "B",
			Description: "Часто встречается у определенных пород (британская, рекс)",
		},
		{
			PetType:     bloodgroup.PetTypeCat,
			BloodGroup:  "AB",
			Description: "Очень редкая группа крови",
		},
	}

	for _, g := range groups {
		// Проверяем, существует ли уже такая группа крови
		exists, err := client.BloodGroup.Query().
			Where(
				bloodgroup.PetTypeEQ(g.PetType),
				bloodgroup.BloodGroupEQ(g.BloodGroup),
			).
			Exist(ctx)

		if err != nil {
			log.ErrorContext(ctx, "Ошибка при проверке существования группы крови", "error", err)
			return err
		}

		if !exists {
			// Если не существует, создаем новую запись
			err := client.BloodGroup.Create().
				SetPetType(g.PetType).
				SetBloodGroup(g.BloodGroup).
				SetDescription(g.Description).
				Exec(ctx)

			if err != nil {
				log.ErrorContext(ctx, "Ошибка при создании группы крови",
					"pet_type", string(g.PetType),
					"blood_group", g.BloodGroup,
					"error", err,
				)
				return err
			}
			log.InfoContext(ctx, "Группа крови добавлена",
				"pet_type", string(g.PetType),
				"blood_group", g.BloodGroup,
			)
		} else {
			log.DebugContext(ctx, "Группа крови уже существует",
				"pet_type", string(g.PetType),
				"blood_group", g.BloodGroup,
			)
		}
	}

	log.InfoContext(ctx, "Заполнение таблицы групп крови завершено")
	return nil
}

// SeedBloodComponents заполняет таблицу компонентов крови начальными данными через ENT
func SeedBloodComponents(ctx context.Context, client *ent.Client, log *slog.Logger) error {
	components := []string{
		"Цельная кровь",
		"Эритроцитарная масса",
		"Свежезамороженная плазма",
		"Замороженная плазма",
		"Тромбоконцентрат",
		"Обогащенная тромбоцитами плазма",
		"Криопреципитат",
		"Криосупернатант",
	}

	for _, name := range components {
		// Проверяем, существует ли уже такой компонент крови
		exists, err := client.BloodComponent.Query().
			Where(bloodcomponent.NameEQ(name)).
			Exist(ctx)

		if err != nil {
			log.ErrorContext(ctx, "Ошибка при проверке существования компонента крови", "error", err)
			return err
		}

		if !exists {
			// Если не существует, создаем новую запись
			err := client.BloodComponent.Create().
				SetName(name).
				Exec(ctx)

			if err != nil {
				log.ErrorContext(ctx, "Ошибка при создании компонента крови",
					"name", name,
					"error", err,
				)
				return err
			}
			log.InfoContext(ctx, "Компонент крови добавлен",
				"name", name,
			)
		} else {
			log.DebugContext(ctx, "Компонент крови уже существует",
				"name", name,
			)
		}
	}

	log.InfoContext(ctx, "Заполнение таблицы компонентов крови завершено")
	return nil
}

// SeedLocations заполняет таблицу локаций начальными данными через ENT
func SeedLocations(ctx context.Context, client *ent.Client, log *slog.Logger) error {
	locations := []string{
		"Москва",
		"Московская область",
	}

	for _, name := range locations {
		exists, err := client.Location.Query().
			Where(location.NameEQ(name)).
			Exist(ctx)

		if err != nil {
			log.ErrorContext(ctx, "Ошибка при проверке существования локации", "error", err)
			return err
		}

		if !exists {
			err := client.Location.Create().
				SetName(name).
				Exec(ctx)

			if err != nil {
				log.ErrorContext(ctx, "Ошибка при создании локации",
					"name", name,
					"error", err,
				)
				return err
			}
			log.InfoContext(ctx, "Локация добавлена", "name", name)
		}
	}

	log.InfoContext(ctx, "Заполнение таблицы локаций завершено")
	return nil
}

// SeedBreeds заполняет таблицу пород начальными данными через ENT
func SeedBreeds(ctx context.Context, client *ent.Client, log *slog.Logger) error {
	for _, b := range data.Breeds {
		exists, err := client.Breed.Query().
			Where(breed.IDEQ(b.ID)).
			Exist(ctx)

		if err != nil {
			log.ErrorContext(ctx, "Ошибка при проверке существования породы", "error", err)
			return err
		}

		if !exists {
			err := client.Breed.Create().
				SetID(b.ID).
				SetName(b.Name).
				SetType(b.Type).
				Exec(ctx)

			if err != nil {
				log.ErrorContext(ctx, "Ошибка при создании породы",
					"id", b.ID,
					"name", b.Name,
					"error", err,
				)
				return err
			}
			log.InfoContext(ctx, "Порода добавлена", "name", b.Name)
		}
	}

	log.InfoContext(ctx, "Заполнение таблицы пород завершено")
	return nil
}
