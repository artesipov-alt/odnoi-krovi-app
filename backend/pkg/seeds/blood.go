package seeds

import (
	"github.com/artesipov-alt/odnoi-krovi-app/internal/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// SeedBloodGroups заполняет таблицу групп крови начальными данными
func SeedBloodGroups(db *gorm.DB, log *zap.Logger) error {
	bloodGroups := []models.BloodGroup{
		// Группы крови для собак
		{
			PetType:     "dog",
			BloodGroup:  "DEA 1.1+",
			Description: "Универсальный донор для собак с положительным DEA 1.1",
		},
		{
			PetType:     "dog",
			BloodGroup:  "DEA 1.1-",
			Description: "Универсальный донор для всех собак",
		},
		{
			PetType:     "dog",
			BloodGroup:  "DEA 1.2+",
			Description: "Положительная группа DEA 1.2",
		},
		{
			PetType:     "dog",
			BloodGroup:  "DEA 1.2-",
			Description: "Отрицательная группа DEA 1.2",
		},
		{
			PetType:     "dog",
			BloodGroup:  "DEA 3+",
			Description: "Положительная группа DEA 3",
		},
		{
			PetType:     "dog",
			BloodGroup:  "DEA 3-",
			Description: "Отрицательная группа DEA 3",
		},
		{
			PetType:     "dog",
			BloodGroup:  "DEA 4+",
			Description: "Положительная группа DEA 4",
		},
		{
			PetType:     "dog",
			BloodGroup:  "DEA 4-",
			Description: "Отрицательная группа DEA 4",
		},
		{
			PetType:     "dog",
			BloodGroup:  "DEA 5+",
			Description: "Положительная группа DEA 5",
		},
		{
			PetType:     "dog",
			BloodGroup:  "DEA 5-",
			Description: "Отрицательная группа DEA 5",
		},
		{
			PetType:     "dog",
			BloodGroup:  "DEA 7+",
			Description: "Положительная группа DEA 7",
		},
		{
			PetType:     "dog",
			BloodGroup:  "DEA 7-",
			Description: "Отрицательная группа DEA 7",
		},

		// Группы крови для кошек
		{
			PetType:     "cat",
			BloodGroup:  "A",
			Description: "Группа крови A - самая распространенная у кошек (около 95%)",
		},
		{
			PetType:     "cat",
			BloodGroup:  "B",
			Description: "Группа крови B - встречается реже (около 5%)",
		},
		{
			PetType:     "cat",
			BloodGroup:  "AB",
			Description: "Группа крови AB - очень редкая (менее 1%)",
		},
	}

	for _, bloodGroup := range bloodGroups {
		// Проверяем, существует ли уже такая группа крови
		var existing models.BloodGroup
		result := db.Where("pet_type = ? AND blood_group = ?", bloodGroup.PetType, bloodGroup.BloodGroup).First(&existing)

		if result.Error == gorm.ErrRecordNotFound {
			// Если не существует, создаем новую запись
			if err := db.Create(&bloodGroup).Error; err != nil {
				log.Error("Ошибка при создании группы крови",
					zap.String("pet_type", bloodGroup.PetType),
					zap.String("blood_group", bloodGroup.BloodGroup),
					zap.Error(err),
				)
				return err
			}
			log.Info("Группа крови добавлена",
				zap.String("pet_type", bloodGroup.PetType),
				zap.String("blood_group", bloodGroup.BloodGroup),
			)
		} else if result.Error != nil {
			log.Error("Ошибка при проверке существования группы крови", zap.Error(result.Error))
			return result.Error
		} else {
			log.Debug("Группа крови уже существует",
				zap.String("pet_type", bloodGroup.PetType),
				zap.String("blood_group", bloodGroup.BloodGroup),
			)
		}
	}

	log.Info("Заполнение таблицы групп крови завершено")
	return nil
}

// SeedBloodComponents заполняет таблицу компонентов крови начальными данными
func SeedBloodComponents(db *gorm.DB, log *zap.Logger) error {
	bloodComponents := []models.BloodComponent{
		{
			Name: "Цельная кровь",
		},
		{
			Name: "Эритроцитарная масса",
		},
		{
			Name: "Плазма",
		},
		{
			Name: "Замороженная плазма",
		},
		{
			Name: "Тромбоконцентрат",
		},
		{
			Name: "Обогащенная тромбоцитами плазма",
		},
		{
			Name: "Криопреципитат",
		},
		{
			Name: "Криосупернатант",
		},
	}

	for _, component := range bloodComponents {
		// Проверяем, существует ли уже такой компонент крови
		var existing models.BloodComponent
		result := db.Where("name = ?", component.Name).First(&existing)

		if result.Error == gorm.ErrRecordNotFound {
			// Если не существует, создаем новую запись
			if err := db.Create(&component).Error; err != nil {
				log.Error("Ошибка при создании компонента крови",
					zap.String("name", component.Name),
					zap.Error(err),
				)
				return err
			}
			log.Info("Компонент крови добавлен",
				zap.String("name", component.Name),
			)
		} else if result.Error != nil {
			log.Error("Ошибка при проверке существования компонента крови", zap.Error(result.Error))
			return result.Error
		} else {
			log.Debug("Компонент крови уже существует",
				zap.String("name", component.Name),
			)
		}
	}

	log.Info("Заполнение таблицы компонентов крови завершено")
	return nil
}
