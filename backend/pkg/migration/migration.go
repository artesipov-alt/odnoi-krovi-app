package migration

import (
	"github.com/artesipov-alt/odnoi-krovi-app/internal/models"
	"github.com/artesipov-alt/odnoi-krovi-app/pkg/seeds"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// AutoMigrate выполняет автоматическую миграцию всех моделей
func AutoMigrate(db *gorm.DB, logger *zap.Logger) {
	modelsToMigrate := []any{
		&models.User{},
		&models.Pet{},
		&models.PetHealth{},
		&models.PetTreatment{},
		&models.PetAnalysis{},
		&models.Breed{},
		&models.BloodComponent{},
		&models.BloodGroup{},
		&models.Location{},
	}

	// Автоматическая миграция всех моделей
	if err := db.AutoMigrate(modelsToMigrate...); err != nil {
		logger.Fatal("Ошибка автоматической миграции", zap.Error(err))
	}

	// Создание последовательностей для префиксных ID
	sequences := []string{
		"CREATE SEQUENCE IF NOT EXISTS vet_clinic_id_seq START 1",
		"CREATE SEQUENCE IF NOT EXISTS user_id_seq START 1",
		"CREATE SEQUENCE IF NOT EXISTS pet_id_seq START 1",
	}

	for _, sql := range sequences {
		if err := db.Exec(sql).Error; err != nil {
			logger.Fatal("Ошибка создания последовательности", zap.Error(err), zap.String("sql", sql))
		}
	}

	logger.Info("Автоматическая миграция выполнена успешно")
}

// SeedDatabase заполняет базу данных начальными данными
func SeedDatabase(db *gorm.DB, logger *zap.Logger) {
	logger.Info("Начало заполнения базы данных начальными данными")

	// Заполнение групп крови
	if err := seeds.SeedBloodGroups(db, logger); err != nil {
		logger.Error("Ошибка при заполнении групп крови", zap.Error(err))
	}

	// Заполнение компонентов крови
	if err := seeds.SeedBloodComponents(db, logger); err != nil {
		logger.Error("Ошибка при заполнении компонентов крови", zap.Error(err))
	}

	logger.Info("Заполнение базы данных завершено")
}
