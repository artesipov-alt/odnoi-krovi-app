package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	bloodsearchv1 "github.com/artesipov-alt/odnoi-krovi-app/microservices/blood-microservice/gen/api/bloodsearch/v1"
	"github.com/redis/go-redis/v9"
)

const (
	// DefaultTTLForRecipient - TTL по умолчанию для реципиентов (30 минут)
	DefaultTTLForRecipient = 30 * time.Minute
	// PetKeyPrefix - префикс для ключей питомцев
	PetKeyPrefix = "pet:"
	// IndexPetTypePrefix - префикс для индекса по типу питомца
	IndexPetTypePrefix = "index:pet_type:"
	// IndexBloodGroupPrefix - префикс для индекса по группе крови
	IndexBloodGroupPrefix = "index:blood_group:"
	// IndexRegionPrefix - префикс для индекса по региону
	IndexRegionPrefix = "index:region:"
	// BloodPoolZSet - отсортированный набор с timestamp'ом истечения для пула поиска крови
	// score = unix timestamp истечения записи (time.Now().Add(ttl).Unix())
	BloodPoolZSet = "blood_pool"
)

// RedisPetRepository реализует PetRepository с использованием Redis
type RedisPetRepository struct {
	client *redis.Client
}

// NewRedisPetRepository создает новый экземпляр RedisPetRepository
func NewRedisPetRepository(client *redis.Client) *RedisPetRepository {
	return &RedisPetRepository{client: client}
}

// AddPet добавляет или обновляет информацию о питомце в Redis с TTL
func (r *RedisPetRepository) AddPet(ctx context.Context, pet *bloodsearchv1.PetRow, ttl time.Duration) error {
	// Используем транзакцию для атомарности
	pipe := r.client.TxPipeline()

	// Сериализуем данные питомца в JSON
	data, err := json.Marshal(pet)
	if err != nil {
		return fmt.Errorf("failed to marshal pet data: %w", err)
	}

	// Сохраняем в Redis с ключом pet:{pet_id}
	key := PetKeyPrefix + pet.PetId
	pipe.Set(ctx, key, data, ttl)

	// Обновляем индексы с тем же TTL
	if err := r.updateIndexesWithTTL(ctx, pipe, pet, ttl); err != nil {
		return fmt.Errorf("failed to update indexes: %w", err)
	}

	// Добавляем в ZSET индекс пула доноров/реципиентов с score = unix timestamp истечения (если ttl > 0).
	// Это позволит быстро получать активных доноров/реципиентов и чистить устаревшие записи.
	if ttl > 0 {
		expireAt := time.Now().Add(ttl).Unix()
		pipe.ZAdd(ctx, BloodPoolZSet, redis.Z{Score: float64(expireAt), Member: key})
	}

	// Выполняем транзакцию
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to save pet to Redis: %w", err)
	}

	return nil
}

// GetPetByID возвращает питомца по его идентификатору
func (r *RedisPetRepository) GetPetByID(ctx context.Context, petID string) (*bloodsearchv1.PetRow, error) {
	key := PetKeyPrefix + petID
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Питомец не найден
		}
		return nil, fmt.Errorf("failed to get pet from Redis: %w", err)
	}

	var pet bloodsearchv1.PetRow
	if err := json.Unmarshal(data, &pet); err != nil {
		return nil, fmt.Errorf("failed to unmarshal pet data: %w", err)
	}

	return &pet, nil
}

// GetPetsByCriteria возвращает список питомцев по заданным критериям
func (r *RedisPetRepository) GetPetsByCriteria(ctx context.Context, criteria *bloodsearchv1.GetPetRows) ([]*bloodsearchv1.PetRow, error) {
	// Собираем ключи для поиска
	var keys []string

	// Поиск по типу питомца
	if criteria.PetType != "" {
		typeKey := IndexPetTypePrefix + criteria.PetType
		typeKeys, err := r.client.SMembers(ctx, typeKey).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to get pets by type: %w", err)
		}
		if len(keys) == 0 {
			keys = typeKeys
		} else {
			keys = intersect(keys, typeKeys)
		}
	}

	// Поиск по группе крови
	if criteria.BloodGroup != "" {
		bloodKey := IndexBloodGroupPrefix + criteria.BloodGroup
		bloodKeys, err := r.client.SMembers(ctx, bloodKey).Result()
		if err != nil {
			return nil, fmt.Errorf("failed to get pets by blood group: %w", err)
		}
		if len(keys) == 0 {
			keys = bloodKeys
		} else {
			keys = intersect(keys, bloodKeys)
		}
	}

	// Поиск по регионам
	if len(criteria.Regions) > 0 {
		var regionKeys []string
		for _, region := range criteria.Regions {
			regionKey := fmt.Sprintf("%s%d", IndexRegionPrefix, region)
			rKeys, err := r.client.SMembers(ctx, regionKey).Result()
			if err != nil {
				return nil, fmt.Errorf("failed to get pets by region: %w", err)
			}
			regionKeys = append(regionKeys, rKeys...)
		}
		if len(keys) == 0 {
			keys = regionKeys
		} else {
			keys = intersect(keys, regionKeys)
		}
	}

	// Если нет критериев, получаем все ключи с использованием SCAN (безопаснее чем KEYS)
	if len(keys) == 0 {
		allKeys, err := r.scanKeys(ctx, PetKeyPrefix+"*")
		if err != nil {
			return nil, fmt.Errorf("failed to get all pet keys: %w", err)
		}
		keys = allKeys
	}

	// Получаем данные питомцев, фильтруя несуществующие (с истекшим TTL)
	var pets []*bloodsearchv1.PetRow
	for _, key := range keys {
		data, err := r.client.Get(ctx, key).Bytes()
		if err != nil {
			if err == redis.Nil {
				// Ключ не существует (возможно истек TTL)
				// Удаляем ключ из индексов при следующем обновлении
				continue
			}
			continue // Пропускаем если другая ошибка
		}

		var pet bloodsearchv1.PetRow
		if err := json.Unmarshal(data, &pet); err != nil {
			continue // Пропускаем если не удалось распарсить
		}

		pets = append(pets, &pet)
	}

	return pets, nil
}

// UpdatePetStatus обновляет статус питомца с возможностью обновления TTL
func (r *RedisPetRepository) UpdatePetStatus(ctx context.Context, petID string, status string, ttl time.Duration) error {
	pet, err := r.GetPetByID(ctx, petID)
	if err != nil {
		return err
	}
	if pet == nil {
		return fmt.Errorf("pet not found: %s", petID)
	}

	pet.Status = status
	return r.AddPet(ctx, pet, ttl)
}

// DeletePet удаляет питомца из хранилища
func (r *RedisPetRepository) DeletePet(ctx context.Context, petID string) error {
	// Получаем питомца для удаления из индексов
	pet, err := r.GetPetByID(ctx, petID)
	if err != nil {
		return err
	}
	if pet == nil {
		return nil // Питомец уже не существует
	}

	// Используем транзакцию для атомарности
	pipe := r.client.TxPipeline()

	// Удаляем из индексов
	if err := r.removeFromIndexes(ctx, pipe, pet); err != nil {
		return fmt.Errorf("failed to remove from indexes: %w", err)
	}

	// Удаляем основную запись
	key := PetKeyPrefix + petID
	pipe.Del(ctx, key)

	// Удаляем из ZSET пула (на случай, если запись еще присутствует в индексе)
	pipe.ZRem(ctx, BloodPoolZSet, key)

	// Выполняем транзакцию
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete pet from Redis: %w", err)
	}

	return nil
}

// GetPetsByRegion возвращает питомцев в указанном регионе
func (r *RedisPetRepository) GetPetsByRegion(ctx context.Context, region int32) ([]*bloodsearchv1.PetRow, error) {
	regionKey := fmt.Sprintf("%s%d", IndexRegionPrefix, region)
	keys, err := r.client.SMembers(ctx, regionKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get pets by region: %w", err)
	}

	// Получаем данные питомцев, фильтруя несуществующие (с истекшим TTL)
	var pets []*bloodsearchv1.PetRow
	for _, key := range keys {
		data, err := r.client.Get(ctx, key).Bytes()
		if err != nil {
			if err == redis.Nil {
				// Ключ не существует (возможно истек TTL)
				continue
			}
			continue
		}

		var pet bloodsearchv1.PetRow
		if err := json.Unmarshal(data, &pet); err != nil {
			continue
		}

		pets = append(pets, &pet)
	}

	return pets, nil
}

// GetPetsByBloodGroup возвращает питомцев с указанной группой крови
func (r *RedisPetRepository) GetPetsByBloodGroup(ctx context.Context, bloodGroup string) ([]*bloodsearchv1.PetRow, error) {
	bloodKey := IndexBloodGroupPrefix + bloodGroup
	keys, err := r.client.SMembers(ctx, bloodKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get pets by blood group: %w", err)
	}

	// Получаем данные питомцев, фильтруя несуществующие (с истекшим TTL)
	var pets []*bloodsearchv1.PetRow
	for _, key := range keys {
		data, err := r.client.Get(ctx, key).Bytes()
		if err != nil {
			if err == redis.Nil {
				// Ключ не существует (возможно истек TTL)
				continue
			}
			continue
		}

		var pet bloodsearchv1.PetRow
		if err := json.Unmarshal(data, &pet); err != nil {
			continue
		}

		pets = append(pets, &pet)
	}

	return pets, nil
}

// GetAllPets возвращает всех питомцев (с пагинацией)
func (r *RedisPetRepository) GetAllPets(ctx context.Context, limit, offset int64) ([]*bloodsearchv1.PetRow, error) {
	// Получаем все ключи питомцев с использованием SCAN
	allKeys, err := r.scanKeys(ctx, PetKeyPrefix+"*")
	if err != nil {
		return nil, fmt.Errorf("failed to get all pet keys: %w", err)
	}

	// Применяем пагинацию
	start := offset
	end := offset + limit
	if start > int64(len(allKeys)) {
		return []*bloodsearchv1.PetRow{}, nil
	}
	if end > int64(len(allKeys)) {
		end = int64(len(allKeys))
	}
	keys := allKeys[start:end]

	// Получаем данные питомцев, фильтруя несуществующие (с истекшим TTL)
	var pets []*bloodsearchv1.PetRow
	for _, key := range keys {
		data, err := r.client.Get(ctx, key).Bytes()
		if err != nil {
			if err == redis.Nil {
				// Ключ не существует (возможно истек TTL)
				continue
			}
			continue
		}

		var pet bloodsearchv1.PetRow
		if err := json.Unmarshal(data, &pet); err != nil {
			continue
		}

		pets = append(pets, &pet)
	}

	return pets, nil
}

// Exists проверяет существование питомца
func (r *RedisPetRepository) Exists(ctx context.Context, petID string) (bool, error) {
	key := PetKeyPrefix + petID
	exists, err := r.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check pet existence: %w", err)
	}
	return exists == 1, nil
}

// Count возвращает общее количество питомцев в хранилище
func (r *RedisPetRepository) Count(ctx context.Context) (int64, error) {
	keys, err := r.scanKeys(ctx, PetKeyPrefix+"*")
	if err != nil {
		return 0, fmt.Errorf("failed to count pets: %w", err)
	}
	return int64(len(keys)), nil
}

// GetTTL возвращает оставшееся время жизни записи
func (r *RedisPetRepository) GetTTL(ctx context.Context, petID string) (time.Duration, error) {
	key := PetKeyPrefix + petID
	ttl, err := r.client.TTL(ctx, key).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to get TTL for pet: %w", err)
	}

	// Если TTL равен -1 (без TTL) или -2 (ключ не существует), возвращаем 0
	if ttl == -1 || ttl == -2 {
		return 0, nil
	}

	return ttl, nil
}

// GetActivePetsFromPool возвращает всех активных питомцев из ZSET-пула (те, у которых score > now)
// Возвращаемые записи дополнительно проверяются на существование и парсятся из JSON.
func (r *RedisPetRepository) GetActivePetsFromPool(ctx context.Context) ([]*bloodsearchv1.PetRow, error) {
	now := time.Now().Unix()
	min := fmt.Sprintf("%d", now+1) // strictly greater than now
	keys, err := r.client.ZRangeByScore(ctx, BloodPoolZSet, &redis.ZRangeBy{
		Min: min,
		Max: "+inf",
	}).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get active pets from pool: %w", err)
	}

	var pets []*bloodsearchv1.PetRow
	for _, key := range keys {
		data, err := r.client.Get(ctx, key).Bytes()
		if err != nil {
			if err == redis.Nil {
				// Ключ не существует — возможно TTL истек, будем чистить периодически
				continue
			}
			// если другая ошибка — пропускаем текущую запись
			continue
		}

		var pet bloodsearchv1.PetRow
		if err := json.Unmarshal(data, &pet); err != nil {
			continue
		}
		pets = append(pets, &pet)
	}

	return pets, nil
}

// CleanExpiredPool удаляет из ZSET все записи с истекшим временем (score <= now).
// Опционально пытается удалить соответствующие ключи (если они все еще существуют).
func (r *RedisPetRepository) CleanExpiredPool(ctx context.Context) error {
	now := time.Now().Unix()
	max := fmt.Sprintf("%d", now)

	// Получаем устаревшие члены, чтобы попытаться удалить соответствующие ключи (если нужно)
	expiredMembers, err := r.client.ZRangeByScore(ctx, BloodPoolZSet, &redis.ZRangeBy{
		Min: "-inf",
		Max: max,
	}).Result()
	if err != nil {
		return fmt.Errorf("failed to fetch expired pool members: %w", err)
	}

	// Если ничего нет — всё ок, но всё равно очищаем диапазон (без ошибок)
	if len(expiredMembers) == 0 {
		_, _ = r.client.ZRemRangeByScore(ctx, BloodPoolZSet, "-inf", max).Result()
		return nil
	}

	pipe := r.client.TxPipeline()
	// Удаляем range из ZSET
	pipe.ZRemRangeByScore(ctx, BloodPoolZSet, "-inf", max)
	// Пытаемся удалить соответствующие ключи (если они ещё существуют)
	for _, member := range expiredMembers {
		pipe.Del(ctx, member)
	}
	_, err = pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to clean expired pool entries: %w", err)
	}
	return nil
}

// updateIndexesWithTTL обновляет индексы для питомца с TTL
func (r *RedisPetRepository) updateIndexesWithTTL(ctx context.Context, pipe redis.Pipeliner, pet *bloodsearchv1.PetRow, ttl time.Duration) error {
	key := PetKeyPrefix + pet.PetId

	// Индекс по типу питомца
	typeKey := IndexPetTypePrefix + pet.PetType
	pipe.SAdd(ctx, typeKey, key)
	if ttl > 0 {
		// Устанавливаем TTL для индекса
		pipe.Expire(ctx, typeKey, ttl)
	}

	// Индекс по группе крови
	bloodKey := IndexBloodGroupPrefix + pet.BloodGroup
	pipe.SAdd(ctx, bloodKey, key)
	if ttl > 0 {
		pipe.Expire(ctx, bloodKey, ttl)
	}

	// Индексы по регионам
	for _, region := range pet.Regions {
		regionKey := fmt.Sprintf("%s%d", IndexRegionPrefix, region)
		pipe.SAdd(ctx, regionKey, key)
		if ttl > 0 {
			pipe.Expire(ctx, regionKey, ttl)
		}
	}

	// Устанавливаем TTL для индексов (Redis автоматически удалит их при истечении)
	// Это важно, чтобы индексы не содержали ссылки на несуществующие ключи

	return nil
}

// removeFromIndexes удаляет питомца из индексов
func (r *RedisPetRepository) removeFromIndexes(ctx context.Context, pipe redis.Pipeliner, pet *bloodsearchv1.PetRow) error {
	key := PetKeyPrefix + pet.PetId

	// Удаляем из индекса типа питомца
	typeKey := IndexPetTypePrefix + pet.PetType
	pipe.SRem(ctx, typeKey, key)

	// Удаляем из индекса группы крови
	bloodKey := IndexBloodGroupPrefix + pet.BloodGroup
	pipe.SRem(ctx, bloodKey, key)

	// Удаляем из индексов регионов
	for _, region := range pet.Regions {
		regionKey := fmt.Sprintf("%s%d", IndexRegionPrefix, region)
		pipe.SRem(ctx, regionKey, key)
	}

	return nil
}

// scanKeys безопасно сканирует ключи по шаблону
func (r *RedisPetRepository) scanKeys(ctx context.Context, pattern string) ([]string, error) {
	var keys []string
	var cursor uint64 = 0
	const count = 100 // Количество ключей за один SCAN

	for {
		var scannedKeys []string
		var err error
		scannedKeys, cursor, err = r.client.Scan(ctx, cursor, pattern, count).Result()
		if err != nil {
			return nil, err
		}

		keys = append(keys, scannedKeys...)

		if cursor == 0 {
			break
		}
	}

	return keys, nil
}

// intersect возвращает пересечение двух срезов строк
func intersect(a, b []string) []string {
	m := make(map[string]bool)
	for _, item := range a {
		m[item] = true
	}

	var result []string
	for _, item := range b {
		if m[item] {
			result = append(result, item)
		}
	}
	return result
}
