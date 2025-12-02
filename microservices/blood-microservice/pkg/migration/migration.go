package migration

import (
	"github.com/artesipov-alt/odnoi-krovi-app/microservices/blood-microservice/internal/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// AutoMigrate выполняет автоматическую миграцию всех моделей для blood microservice
func AutoMigrate(db *gorm.DB, logger *zap.Logger) {
	modelsToMigrate := []any{
		&models.PetRow{},
	}

	// Автоматическая миграция всех моделей
	if err := db.AutoMigrate(modelsToMigrate...); err != nil {
		logger.Fatal("Ошибка автоматической миграции blood microservice", zap.Error(err))
	}

	// Создание индексов для оптимизации запросов
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_pet_rows_pet_id ON pet_rows(pet_id)",
		"CREATE INDEX IF NOT EXISTS idx_pet_rows_blood_group ON pet_rows(blood_group)",
		"CREATE INDEX IF NOT EXISTS idx_pet_rows_pet_type ON pet_rows(pet_type)",
		"CREATE INDEX IF NOT EXISTS idx_pet_rows_status ON pet_rows(status)",
		"CREATE INDEX IF NOT EXISTS idx_pet_rows_regions_gin ON pet_rows USING gin(regions)",
		"CREATE INDEX IF NOT EXISTS idx_pet_rows_created_at ON pet_rows(created_at)",
	}

	for _, sql := range indexes {
		if err := db.Exec(sql).Error; err != nil {
			logger.Error("Ошибка создания индекса", zap.Error(err), zap.String("sql", sql))
		}
	}

	logger.Info("Автоматическая миграция blood microservice выполнена успешно")
}
