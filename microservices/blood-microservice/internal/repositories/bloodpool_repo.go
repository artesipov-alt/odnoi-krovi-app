package repositories

import (
	"context"
	"time"

	bloodpoolv1 "github.com/artesipov-alt/odnoi-krovi-app/microservices/blood-microservice/gen/api/bloodpool/v1"
)

// PetRepository определяет интерфейс для работы с данными питомцев в Redis
type PetRepository interface {
	// AddPet добавляет или обновляет информацию о питомце в Redis с TTL
	AddPet(ctx context.Context, pet *bloodpoolv1.PetRow, ttl time.Duration) error

	// GetPetByID возвращает питомца по его идентификатору
	GetPetByID(ctx context.Context, petID string) (*bloodpoolv1.PetRow, error)

	// GetPetsByCriteria возвращает список питомцев по заданным критериям
	GetPetsByCriteria(ctx context.Context, criteria *bloodpoolv1.GetPetRows) ([]*bloodpoolv1.PetRow, error)

	// UpdatePetStatus обновляет статус питомца с возможностью обновления TTL
	UpdatePetStatus(ctx context.Context, petID string, status string, ttl time.Duration) error

	// DeletePet удаляет питомца из хранилища
	DeletePet(ctx context.Context, petID string) error

	// GetPetsByRegion возвращает питомцев в указанном регионе
	GetPetsByRegion(ctx context.Context, region int32) ([]*bloodpoolv1.PetRow, error)

	// GetPetsByBloodGroup возвращает питомцев с указанной группой крови
	GetPetsByBloodGroup(ctx context.Context, bloodGroup string) ([]*bloodpoolv1.PetRow, error)

	// GetAllPets возвращает всех питомцев (с пагинацией)
	GetAllPets(ctx context.Context, limit, offset int64) ([]*bloodpoolv1.PetRow, error)

	// Exists проверяет существование питомца
	Exists(ctx context.Context, petID string) (bool, error)

	// Count возвращает общее количество питомцев в хранилище
	Count(ctx context.Context) (int64, error)

	// GetTTL возвращает оставшееся время жизни записи
	GetTTL(ctx context.Context, petID string) (time.Duration, error)
}
