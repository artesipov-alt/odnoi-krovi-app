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
		&models.VetClinic{},
		&models.Breed{},
		&models.BloodSearch{},
		&models.BloodStock{},
		&models.BloodComponent{},
		&models.BloodGroup{},
		&models.Location{},
		&models.Donation{},
	}

	// Автоматическая миграция всех моделей
	if err := db.AutoMigrate(modelsToMigrate...); err != nil {
		logger.Fatal("Ошибка автоматической миграции", zap.Error(err))
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
