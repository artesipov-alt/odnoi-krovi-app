package migration

import (
	"github.com/artesipov-alt/odnoi-krovi-app/microservices/blood-microservice/internal/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// AutoMigrate выполняет автоматическую миграцию всех моделей для blood microservice
func AutoMigrate(db *gorm.DB, logger *zap.Logger) {
	modelsToMigrate := []any{
		&models.BloodRequest{},
	}

	// Автоматическая миграция всех моделей
	if err := db.AutoMigrate(modelsToMigrate...); err != nil {
		logger.Fatal("Ошибка автоматической миграции blood microservice", zap.Error(err))
	}

	// Создание индексов для оптимизации запросов
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_blood_requests_pet_id ON blood_requests(pet_id)",
		"CREATE INDEX IF NOT EXISTS idx_blood_requests_blood_group ON blood_requests(blood_group)",
		"CREATE INDEX IF NOT EXISTS idx_blood_requests_pet_type ON blood_requests(pet_type)",
		"CREATE INDEX IF NOT EXISTS idx_blood_requests_status ON blood_requests(status)",
		"CREATE INDEX IF NOT EXISTS idx_blood_requests_regions_gin ON blood_requests USING gin(regions)",
		"CREATE INDEX IF NOT EXISTS idx_blood_requests_created_at ON blood_requests(created_at)",
	}

	for _, sql := range indexes {
		if err := db.Exec(sql).Error; err != nil {
			logger.Error("Ошибка создания индекса", zap.Error(err), zap.String("sql", sql))
		}
	}

	logger.Info("Автоматическая миграция blood microservice выполнена успешно")
}
