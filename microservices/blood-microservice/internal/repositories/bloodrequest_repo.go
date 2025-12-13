package repositories

import (
	"context"

	bloodrequestv1 "github.com/artesipov-alt/odnoi-krovi-app/microservices/blood-microservice/gen/api/bloodrequest/v1"
)

// BloodRequestRepository определяет интерфейс для работы с данными заявок на поиск крови питомцев
type BloodRequestRepository interface {
	// AddRequest добавляет или обновляет информацию о заявке
	AddRequest(ctx context.Context, request *bloodrequestv1.BloodRequest) error

	// GetRequestByID возвращает заявку по её идентификатору
	GetRequestByID(ctx context.Context, requestID string) (*bloodrequestv1.BloodRequest, error)

	// GetRequestsByCriteria возвращает список заявок по заданным критериям
	GetRequestsByCriteria(ctx context.Context, criteria *bloodrequestv1.GetBloodRequests) ([]*bloodrequestv1.BloodRequest, error)

	// UpdateRequestStatus обновляет статус заявки
	UpdateRequestStatus(ctx context.Context, requestID string, status string) error

	// DeleteRequest удаляет заявку из хранилища
	DeleteRequest(ctx context.Context, requestID string) error

	// GetRequestsByRegion возвращает заявки в указанном регионе
	GetRequestsByRegion(ctx context.Context, region int32) ([]*bloodrequestv1.BloodRequest, error)

	// GetRequestsByBloodGroup возвращает заявки с указанной группой крови
	GetRequestsByBloodGroup(ctx context.Context, bloodGroup string) ([]*bloodrequestv1.BloodRequest, error)

	// GetAllRequests возвращает все заявки (с пагинацией)
	GetAllRequests(ctx context.Context, limit, offset int64) ([]*bloodrequestv1.BloodRequest, error)

	// Exists проверяет существование заявки
	Exists(ctx context.Context, requestID string) (bool, error)

	// Count возвращает общее количество заявок в хранилище
	Count(ctx context.Context) (int64, error)
}
