package v1

import (
	"context"
	"time"

	bloodpoolv1 "github.com/artesipov-alt/odnoi-krovi-app/microservices/blood-microservice/gen/api/bloodpool/v1"
	"github.com/artesipov-alt/odnoi-krovi-app/microservices/blood-microservice/internal/repositories"
)

const (
	// RecipientTTL - время жизни реципиента в пуле (30 минут)
	RecipientTTL = 30 * time.Minute
	// RecipientStatusActive - статус активного реципиента
	RecipientStatusActive = "active"
	// RecipientStatusMatched - статус реципиента, на которого нашелся донор
	RecipientStatusMatched = "matched"
	// RecipientStatusExpired - статус реципиента с истекшим TTL
	RecipientStatusExpired = "expired"
)

// RecipientService сервис для работы с реципиентами
type RecipientService struct {
	repo repositories.PetRepository
}

// NewRecipientService создает новый экземпляр RecipientService
func NewRecipientService(repo repositories.PetRepository) *RecipientService {
	return &RecipientService{repo: repo}
}

// AddRecipient добавляет реципиента в пул с TTL 30 минут
func (s *RecipientService) AddRecipient(ctx context.Context, recipient *bloodpoolv1.PetRow) error {
	// Устанавливаем статус активного реципиента
	recipient.Status = RecipientStatusActive

	// Добавляем реципиента с TTL 30 минут
	return s.repo.AddPet(ctx, recipient, RecipientTTL)
}

// GetRecipientByID возвращает реципиента по ID
func (s *RecipientService) GetRecipientByID(ctx context.Context, recipientID string) (*bloodpoolv1.PetRow, error) {
	return s.repo.GetPetByID(ctx, recipientID)
}

// FindMatchingRecipients ищет реципиентов по критериям
func (s *RecipientService) FindMatchingRecipients(ctx context.Context, criteria *bloodpoolv1.GetPetRows) ([]*bloodpoolv1.PetRow, error) {
	// Получаем всех реципиентов по критериям
	recipients, err := s.repo.GetPetsByCriteria(ctx, criteria)
	if err != nil {
		return nil, err
	}

	// Фильтруем только активных реципиентов
	var activeRecipients []*bloodpoolv1.PetRow
	for _, recipient := range recipients {
		if recipient.Status == RecipientStatusActive {
			activeRecipients = append(activeRecipients, recipient)
		}
	}

	return activeRecipients, nil
}

// MarkAsMatched помечает реципиента как "найден донор"
func (s *RecipientService) MarkAsMatched(ctx context.Context, recipientID string) error {
	// При нахождении донора удаляем реципиента из пула
	// (или можно пометить как matched и оставить для истории)
	return s.repo.DeletePet(ctx, recipientID)
}

// GetRemainingTime возвращает оставшееся время жизни реципиента
func (s *RecipientService) GetRemainingTime(ctx context.Context, recipientID string) (time.Duration, error) {
	return s.repo.GetTTL(ctx, recipientID)
}

// ExtendTTL продлевает время жизни реципиента
func (s *RecipientService) ExtendTTL(ctx context.Context, recipientID string) error {
	// Получаем текущего реципиента
	recipient, err := s.repo.GetPetByID(ctx, recipientID)
	if err != nil {
		return err
	}

	if recipient == nil {
		return nil // Реципиент не найден (возможно истек TTL)
	}

	// Обновляем с новым TTL (еще 30 минут)
	return s.repo.UpdatePetStatus(ctx, recipientID, recipient.Status, RecipientTTL)
}

// GetActiveRecipientsCount возвращает количество активных реципиентов
func (s *RecipientService) GetActiveRecipientsCount(ctx context.Context) (int64, error) {
	// Получаем всех реципиентов
	allRecipients, err := s.repo.GetAllPets(ctx, 0, 0)
	if err != nil {
		return 0, err
	}

	// Считаем только активных
	var count int64
	for _, recipient := range allRecipients {
		if recipient.Status == RecipientStatusActive {
			count++
		}
	}

	return count, nil
}

// CleanupExpiredRecipients удаляет реципиентов с истекшим TTL
// Внимание: Redis автоматически удаляет ключи с истекшим TTL,
// но этот метод можно использовать для дополнительной очистки индексов
func (s *RecipientService) CleanupExpiredRecipients(ctx context.Context) (int, error) {
	// Получаем всех реципиентов
	allRecipients, err := s.repo.GetAllPets(ctx, 0, 0)
	if err != nil {
		return 0, err
	}

	cleanedCount := 0

	// Проверяем каждого реципиента
	for _, recipient := range allRecipients {
		// Проверяем TTL
		ttl, err := s.repo.GetTTL(ctx, recipient.PetId)
		if err != nil {
			continue
		}

		// Если TTL истек или равен 0
		if ttl <= 0 {
			// Удаляем реципиента
			if err := s.repo.DeletePet(ctx, recipient.PetId); err == nil {
				cleanedCount++
			}
		}
	}

	return cleanedCount, nil
}

// GetRecipientsByRegion возвращает активных реципиентов в регионе
func (s *RecipientService) GetRecipientsByRegion(ctx context.Context, region int32) ([]*bloodpoolv1.PetRow, error) {
	recipients, err := s.repo.GetPetsByRegion(ctx, region)
	if err != nil {
		return nil, err
	}

	// Фильтруем только активных
	var activeRecipients []*bloodpoolv1.PetRow
	for _, recipient := range recipients {
		if recipient.Status == RecipientStatusActive {
			activeRecipients = append(activeRecipients, recipient)
		}
	}

	return activeRecipients, nil
}

// GetRecipientsByBloodGroup возвращает активных реципиентов с группой крови
func (s *RecipientService) GetRecipientsByBloodGroup(ctx context.Context, bloodGroup string) ([]*bloodpoolv1.PetRow, error) {
	recipients, err := s.repo.GetPetsByBloodGroup(ctx, bloodGroup)
	if err != nil {
		return nil, err
	}

	// Фильтруем только активных
	var activeRecipients []*bloodpoolv1.PetRow
	for _, recipient := range recipients {
		if recipient.Status == RecipientStatusActive {
			activeRecipients = append(activeRecipients, recipient)
		}
	}

	return activeRecipients, nil
}
