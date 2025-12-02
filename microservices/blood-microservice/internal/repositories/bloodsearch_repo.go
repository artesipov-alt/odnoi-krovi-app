package repositories

import (
	"context"

	bloodsearchv1 "github.com/artesipov-alt/odnoi-krovi-app/microservices/blood-microservice/gen/api/bloodsearch/v1"
)

// BloodRequestRepository определяет интерфейс для работы с данными заявок на поиск крови питомцев
type BloodSearchRepository interface {
	// AddRequest добавляет или обновляет информацию о заявке
	AddRequest(ctx context.Context, request *bloodsearchv1.PetRow) error

	// GetRequestByID возвращает заявку по её идентификатору
	GetRequestByID(ctx context.Context, requestID string) (*bloodsearchv1.PetRow, error)

	// GetRequestsByCriteria возвращает список заявок по заданным критериям
	GetRequestsByCriteria(ctx context.Context, criteria *bloodsearchv1.GetPetRows) ([]*bloodsearchv1.PetRow, error)

	// UpdateRequestStatus обновляет статус заявки
	UpdateRequestStatus(ctx context.Context, requestID string, status string) error

	// DeleteRequest удаляет заявку из хранилища
	DeleteRequest(ctx context.Context, requestID string) error

	// GetRequestsByRegion возвращает заявки в указанном регионе
	GetRequestsByRegion(ctx context.Context, region int32) ([]*bloodsearchv1.PetRow, error)

	// GetRequestsByBloodGroup возвращает заявки с указанной группой крови
	GetRequestsByBloodGroup(ctx context.Context, bloodGroup string) ([]*bloodsearchv1.PetRow, error)

	// GetAllRequests возвращает все заявки (с пагинацией)
	GetAllRequests(ctx context.Context, limit, offset int64) ([]*bloodsearchv1.PetRow, error)

	// Exists проверяет существование заявки
	Exists(ctx context.Context, requestID string) (bool, error)

	// Count возвращает общее количество заявок в хранилище
	Count(ctx context.Context) (int64, error)
}
