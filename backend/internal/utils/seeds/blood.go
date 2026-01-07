package seeds

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/bloodcomponent"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/bloodgroup"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/breed"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/location"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/utils/seeds/data"
	"go.uber.org/zap"
)

// SeedBloodGroups заполняет таблицу групп крови начальными данными через ENT
func SeedBloodGroups(ctx context.Context, client *ent.Client, log *zap.Logger) error {
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
			log.Error("Ошибка при проверке существования группы крови", zap.Error(err))
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
				log.Error("Ошибка при создании группы крови",
					zap.String("pet_type", string(g.PetType)),
					zap.String("blood_group", g.BloodGroup),
					zap.Error(err),
				)
				return err
			}
			log.Info("Группа крови добавлена",
				zap.String("pet_type", string(g.PetType)),
				zap.String("blood_group", g.BloodGroup),
			)
		} else {
			log.Debug("Группа крови уже существует",
				zap.String("pet_type", string(g.PetType)),
				zap.String("blood_group", g.BloodGroup),
			)
		}
	}

	log.Info("Заполнение таблицы групп крови завершено")
	return nil
}

// SeedBloodComponents заполняет таблицу компонентов крови начальными данными через ENT
func SeedBloodComponents(ctx context.Context, client *ent.Client, log *zap.Logger) error {
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
			log.Error("Ошибка при проверке существования компонента крови", zap.Error(err))
			return err
		}

		if !exists {
			// Если не существует, создаем новую запись
			err := client.BloodComponent.Create().
				SetName(name).
				Exec(ctx)

			if err != nil {
				log.Error("Ошибка при создании компонента крови",
					zap.String("name", name),
					zap.Error(err),
				)
				return err
			}
			log.Info("Компонент крови добавлен",
				zap.String("name", name),
			)
		} else {
			log.Debug("Компонент крови уже существует",
				zap.String("name", name),
			)
		}
	}

	log.Info("Заполнение таблицы компонентов крови завершено")
	return nil
}

// SeedLocations заполняет таблицу локаций начальными данными через ENT
func SeedLocations(ctx context.Context, client *ent.Client, log *zap.Logger) error {
	locations := []string{
		"Москва",
		"Санкт-Петербург",
		"Новосибирск",
		"Екатеринбург",
		"Казань",
		"Нижний Новгород",
		"Челябинск",
		"Самара",
		"Омск",
		"Ростов-на-Дону",
	}

	for _, name := range locations {
		exists, err := client.Location.Query().
			Where(location.NameEQ(name)).
			Exist(ctx)

		if err != nil {
			log.Error("Ошибка при проверке существования локации", zap.Error(err))
			return err
		}

		if !exists {
			err := client.Location.Create().
				SetName(name).
				Exec(ctx)

			if err != nil {
				log.Error("Ошибка при создании локации",
					zap.String("name", name),
					zap.Error(err),
				)
				return err
			}
			log.Info("Локация добавлена", zap.String("name", name))
		}
	}

	log.Info("Заполнение таблицы локаций завершено")
	return nil
}

// SeedBreeds заполняет таблицу пород начальными данными через ENT
func SeedBreeds(ctx context.Context, client *ent.Client, log *zap.Logger) error {
	for _, b := range data.Breeds {
		exists, err := client.Breed.Query().
			Where(breed.IDEQ(b.ID)).
			Exist(ctx)

		if err != nil {
			log.Error("Ошибка при проверке существования породы", zap.Error(err))
			return err
		}

		if !exists {
			err := client.Breed.Create().
				SetID(b.ID).
				SetName(b.Name).
				SetType(b.Type).
				Exec(ctx)

			if err != nil {
				log.Error("Ошибка при создании породы",
					zap.Int("id", b.ID),
					zap.String("name", b.Name),
					zap.Error(err),
				)
				return err
			}
			log.Info("Порода добавлена", zap.String("name", b.Name))
		}
	}

	log.Info("Заполнение таблицы пород завершено")
	return nil
}
