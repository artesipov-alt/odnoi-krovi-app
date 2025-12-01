package v1

import (
	"context"
	"sort"
	"time"

	bloodpoolv1 "github.com/artesipov-alt/odnoi-krovi-app/microservices/blood-microservice/gen/api/bloodpool/v1"
	"github.com/artesipov-alt/odnoi-krovi-app/microservices/blood-microservice/internal/repositories"
)

type bloodPoolService struct {
	repo repositories.PetRepository
}

func NewBloodPoolService(repo repositories.PetRepository) *bloodPoolService {
	return &bloodPoolService{repo: repo}
}

// AddPet обрабатывает запрос добавления питомца в пул поиска крови
func (s *bloodPoolService) AddPet(
	ctx context.Context,
	req *bloodpoolv1.PetRow,
) (*bloodpoolv1.PetRowStatus, error) {
	// Для реципиентов устанавливаем TTL 30 минут
	// Для доноров можно установить другой TTL или 0 (без TTL)
	var ttl time.Duration

	// Определяем TTL в зависимости от типа питомца или статуса
	// В данном случае предполагаем, что все добавляемые питомцы - реципиенты с TTL 30 минут
	ttl = 30 * time.Minute

	// Устанавливаем статус "active" для реципиента
	req.Status = "active"

	// Сохраняем питомца в репозиторий с TTL
	err := s.repo.AddPet(ctx, req, ttl)
	if err != nil {
		return nil, err
	}

	return &bloodpoolv1.PetRowStatus{
		PetId:  req.PetId,
		Status: "added_with_ttl_30min",
	}, nil
}

// GetPets обрабатывает запрос получения питомцев по критериям
func (s *bloodPoolService) GetPets(
	ctx context.Context,
	req *bloodpoolv1.GetPetRows,
) (*bloodpoolv1.PetRows, error) {
	// Получаем питомцев по критериям
	pets, err := s.repo.GetPetsByCriteria(ctx, req)
	if err != nil {
		return nil, err
	}

	// Дедупликация по pet_id и фильтрация только активных
	seen := make(map[string]struct{}, len(pets))
	var activePets []*bloodpoolv1.PetRow
	for _, pet := range pets {
		if pet.Status != "active" {
			continue
		}
		if _, ok := seen[pet.PetId]; ok {
			continue
		}
		seen[pet.PetId] = struct{}{}
		activePets = append(activePets, pet)
	}

	// Предрасчет оставшегося времени жизни для сортировки
	timeLeft := make(map[string]time.Duration, len(activePets))
	for _, pet := range activePets {
		ttl, err := s.repo.GetTTL(ctx, pet.PetId)
		if err != nil || ttl <= 0 {
			// Если TTL недоступен/ошибка — ставим максимально возможное значение,
			// чтобы такие записи шли в конце при сортировке по времени
			ttl = time.Duration(1<<63 - 1)
		}
		timeLeft[pet.PetId] = ttl
	}

	// Сортировка: сперва по приоритету (desc), затем по оставшемуся времени (asc)
	sort.Slice(activePets, func(i, j int) bool {
		pi, pj := activePets[i], activePets[j]
		if pi.PriorityLevel != pj.PriorityLevel {
			return pi.PriorityLevel > pj.PriorityLevel
		}
		ti := timeLeft[pi.PetId]
		tj := timeLeft[pj.PetId]
		return ti < tj
	})

	return &bloodpoolv1.PetRows{
		Pets: activePets,
	}, nil
}
