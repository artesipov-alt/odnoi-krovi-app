# Архитектура проекта "Одной крови"

## Общее описание

"Одной крови" — это платформа для поиска доноров крови для домашних животных (собак и кошек). Позволяет:
- Регистрировать питомцев с медицинскими данными
- Создавать заявки на поиск крови
- Откликаться на заявки донорами
- Управлять статусами донорства/реципиентства

## Технологический стек

### Backend
- **Язык**: Go 1.21+
- **Фреймворк HTTP**: Huma v2 (OpenAPI-first, code generation)
- **ORM/Database**: ENT (Entity Framework от Facebook)
- **База данных**: PostgreSQL 15+
- **Object Storage**: S3-совместимое хранилище (MinIO/AWS)
- **Миграции**: ENT Migrate (автоматические)
- **Логирование**: slog (structured logging)

### Структура проекта (Clean Architecture)

```
backend/
├── cmd/api/                    # Entry point (main.go)
├── internal/
│   ├── domain/                 # Domain Layer (чистый, без внешних зависимостей)
│   │   ├── pet/               # Агрегат Pet
│   │   ├── bloodsearch/       # Агрегат BloodRequest
│   │   ├── user/              # Агрегат User
│   │   ├── filestorage/       # Агрегат File
│   │   └── reference/         # Справочники (породы, локации, группы крови)
│   ├── application/           # Application Layer (CQRS)
│   │   ├── pet/cmd/          # Команды (Create, Update, Delete)
│   │   ├── pet/query/        # Запросы (GetByID, GetByUser)
│   │   ├── bloodsearch/cmd/  # Команды
│   │   ├── bloodsearch/query/# Запросы
│   │   ├── file/cmd/
│   │   └── user/cmd/, user/query/
│   ├── mapper/                # Mapping Layer (DTO <-> Domain conversions)
│   │   ├── pet_mapper.go     # Pet DTO <-> Domain mappings
│   │   ├── user_mapper.go    # User DTO <-> Domain mappings
│   │   └── blood_request_mapper.go # BloodRequest mappings
│   ├── infra/                 # Infrastructure Layer
│   │   ├── presistance/pg/   # PostgreSQL реализации (ENT)
│   │   ├── presistance/s3/   # S3 реализация
│   │   └── ent/              # Сгенерированные ENT сущности
│   └── transport/http/        # Transport Layer (Huma handlers)
├── pkg/                       # Shared packages (config, logger, enums)
└── docs/                      # OpenAPI specs
```

## Ключевые принципы

### 1. Clean Architecture / Layered Architecture

```
┌─────────────────────────────────────┐
│  Transport (HTTP/Huma)              │
├─────────────────────────────────────┤
│  Mapper (DTO <-> Domain)            │
│  - Преобразование между DTO и       │
│    доменными моделями               │
├─────────────────────────────────────┤
│  Application (CQRS Handlers)        │
│  - cmd: бизнес-логика, транзакции   │
│  - query: чтение с проекциями       │
├─────────────────────────────────────┤
│  Domain (Entities, Repositories)    │
│  - model/: доменные модели          │
│  - *_repo.go: интерфейсы репозиториев│
├─────────────────────────────────────┤
│  Infrastructure (ENT, S3, PostgreSQL)│
└─────────────────────────────────────┘
```

**Правило зависимостей**: Внешние слои зависят от внутренних, но не наоборот.

### 2. CQRS (Command Query Responsibility Segregation)

- **Commands** (`cmd/`): Изменяют состояние, возвращают результат или ошибку
  - Пример: `CreatePetHandler`, `UpdateRequestHandler`, `ApplyForRequestHandler`
  
- **Queries** (`query/`): Только чтение, не изменяют состояние
  - Пример: `GetByIDHandler`, `ListRequestsHandler`

### 3. Repository Pattern

Каждый агрегат имеет интерфейс репозитория в `domain/`, реализацию в `infra/presistance/pg/`.

**Пример разделения:**
```go
// Domain layer
type PetReadRepository interface {
    GetByID(ctx, id string, opts PetPreloadOptions) (*model.Pet, error)
    Exists(ctx, id string) (bool, error)
}

type PetWriteRepository interface {
    Create(ctx, pet *model.Pet) (*model.Pet, error)
    Update(ctx, id string, pet *model.Pet) (*model.Pet, error)
    Delete(ctx, id string) error
}
```

### 4. Domain Models (чистые структуры)

Все доменные модели находятся в `domain/*/model/` и не зависят от инфраструктуры:
- `model.Pet` — не знает про ENT
- `model.BloodRequest` — чистая структура
- Конвертация ENT ↔ Domain происходит в репозиториях

## Текущее состояние архитектуры

### ✅ Что реализовано хорошо

1. **Чистое разделение слоёв** — Domain не зависит от Infra
2. **CQRS** — Разделены команды и запросы
3. **Интерфейсы репозиториев** — Легко мокать для тестов
4. **Dependency Injection** — Явные конструкторы для всех handlers
5. **Обработка ошибок** — Централизованный пакет `apperrors`

### ⚠️ Известные проблемы и ограничения

#### 1. God Object — Pet модель
```go
// Сейчас Pet содержит 20+ полей:
type Pet struct {
    ID, Name, Type, WeightKg, Gender, BirthDate...
    Health *PetHealth           // Вложенная структура
    Treatments *PetTreatment    // Вложенная структура  
    Analyses []*PetAnalysis     // Слайз
    StopFactors []string
    WarnFactors []string
    Bonuses []string
    // ... и т.д.
}
```

**Проблема**: 
- Сложно тестировать
- Нарушает Single Responsibility
- Неясные границы агрегата

#### 2. Отсутствие Value Objects
Используются примитивы без валидации:
```go
WeightKg   float64  // Может быть -10 или 1000
ChipNumber string   // Нет формата
PhotoURLs  []string // Нет ограничений на количество
```

#### 3. Смешение логики валидации
Методы `GetStopFactors()`, `GetWarnFactors()` встроены в модель Pet, но это бизнес-правила, которые должны быть в сервисе/спецификации.

#### 4. Отсутствие Domain Events
Нет механизма публикации событий при изменениях:
- Питомец стал донором
- Создана заявка на кровь
- Донор откликнулся

#### 5. Transaction Management
Транзакции распределены по handlers. Нет единого Unit of Work.

#### 6. N+1 Query проблема
В `GetByUserHandler` для каждого питомца отдельный запрос на blood request.

## Рекомендации по улучшению

### 🔴 Высокий приоритет (сделать в ближайшие 2 недели)

#### 1. Разделить Pet на агрегаты

**Текущее состояние**: Один большой агрегат Pet

**Целевая архитектура**:
```
Pet (корень агрегата)
├── PetID (Value Object)
├── PetProfile (Value Object: Type, Breed, Weight, BirthDate)
├── PetStatus (Value Object: donor/recipient/none)
└── OwnerID (ссылка на User)

PetHealth (отдельный агрегат)
├── PetID (ссылка на Pet)
├── HealthStatus
├── LastDonation
└── Medications

PetTreatment (отдельный агрегат)
├── PetID
├── VaccinationDates
└── DewormingDates
```

**Шаги реализации**:
1. Создать отдельные таблицы (уже есть через ENT)
2. Создать отдельные репозитории: `PetHealthRepository`, `PetTreatmentRepository`
3. Разделить handlers: `UpdatePetHealthHandler`, `UpdatePetTreatmentHandler`
4. Pet агрегат должен содержать только `PetProfile` и `PetStatus`

#### 2. Внедрить Value Objects

**Примеры VO для создания**:

```go
// Weight — вес питомца с валидацией
type Weight struct {
    kilograms float64
}

func NewWeight(kg float64) (Weight, error) {
    if kg <= 0 || kg > 100 {
        return Weight{}, errors.New("weight must be between 0 and 100 kg")
    }
    return Weight{kilograms: kg}, nil
}

// ChipNumber — номер чипа с форматом
type ChipNumber struct {
    value string
}

func NewChipNumber(number string) (ChipNumber, error) {
    if len(number) != 15 {
        return ChipNumber{}, errors.New("chip number must be 15 digits")
    }
    if !isAllDigits(number) {
        return ChipNumber{}, errors.New("chip number must contain only digits")
    }
    return ChipNumber{value: number}, nil
}

// PhotoGallery — коллекция фото с ограничениями
type PhotoGallery struct {
    urls []string
}

func NewPhotoGallery(urls []string) (PhotoGallery, error) {
    if len(urls) > 10 {
        return PhotoGallery{}, errors.New("max 10 photos allowed")
    }
    return PhotoGallery{urls: urls}, nil
}
```

**Где использовать**:
- Заменить `WeightKg float64` → `Weight Weight`
- Заменить `ChipNumber string` → `ChipNumber ChipNumber`
- Заменить `PhotoURLs []string` → `PhotoGallery PhotoGallery`

### 🟡 Средний приоритет (в течение месяца)

#### 3. Domain Events

**Необходимые события**:
```go
// events/pet_events.go
type PetBecameDonorEvent struct {
    PetID     string
    OwnerID   string
    Timestamp time.Time
}

type BloodRequestCreatedEvent struct {
    RequestID string
    PetID     string
    OwnerID   string
    BloodGroupNames []string
}

type DonorRespondedEvent struct {
    ResponseID string
    RequestID  string
    DonorID    string
    DonorOwnerID string
}
```

**Инфраструктура**:
```go
// infra/eventbus/event_bus.go
type EventBus interface {
    Publish(ctx context.Context, event DomainEvent) error
    Subscribe(eventType string, handler EventHandler) error
}

// Примеры обработчиков:
// - PetBecameDonorEvent → отправка уведомления владельцу
// - BloodRequestCreatedEvent → индексация для поиска
// - DonorRespondedEvent → уведомление владельцу заявки
```

**Примечание**: Event-логика — дополнительная, не обязательная для core-функциональности. Можно реализовать позже.

#### 4. Unit of Work

**Проблема**: Сейчас транзакции в handlers через `txManager.WithTx()`.

**Решение**:
```go
// application/uow/unit_of_work.go
type UnitOfWork interface {
    Execute(ctx context.Context, fn func(txCtx context.Context) error) error
}

// В handler:
func (h *CreateHandler) Handle(ctx context.Context, ...) error {
    return h.uow.Execute(ctx, func(txCtx context.Context) error {
        pet, err := h.petRepo.Create(txCtx, pet)
        if err != nil {
            return err
        }
        
        h.eventBus.Publish(txCtx, PetCreatedEvent{...})
        return nil
    })
}
```

#### 5. Specification Pattern для валидации доноров

**Текущее**: Логика в методах `Pet.GetStopFactors()`

**Целевое**:
```go
// domain/pet/specification/donor_specification.go
type DonorSpecification interface {
    IsSatisfiedBy(pet *model.Pet) (bool, model.FactorCode)
}

type AgeSpecification struct{}
func (s AgeSpecification) IsSatisfiedBy(pet *model.Pet) (bool, model.FactorCode) {
    age := calculateAge(pet.BirthDate)
    if age > 8 {
        return false, model.StopFactorTooOld
    }
    return true, ""
}

type VaccinationSpecification struct{}
type HealthSpecification struct{}

// Сервис валидации использует спецификации:
type DonorValidationService struct {
    specifications []DonorSpecification
}

func (s *DonorValidationService) Validate(pet *model.Pet) (stopFactors, warnFactors []model.FactorCode) {
    for _, spec := range s.specifications {
        if satisfied, factor := spec.IsSatisfiedBy(pet); !satisfied {
            stopFactors = append(stopFactors, factor)
        }
    }
    return
}
```

### 🟢 Низкий приоритет (по возможности)

#### 6. Read Models для сложных запросов

**Проблема**: `GetByUserHandler` возвращает полные агрегаты, но UI часто нужна агрегированная информация.

**Решение**:
```go
// application/pet/query/pet_list_item.go
type PetListItem struct {
    PetID           string
    PetName         string
    PetType         string
    Status          string
    ActiveRequests  int
    LastDonation    *time.Time
    PhotoURL        string // только главное фото
}

// Отдельный query handler:
type ListUserPetsHandler struct {
    readModel PetReadModel
}

func (h *ListUserPetsHandler) Handle(ctx context.Context, userID string) ([]PetListItem, error) {
    // Использует JOIN и агрегации на уровне SQL
    return h.readModel.ListByUserID(ctx, userID)
}
```

#### 7. API Versioning

Когда будут breaking changes:
```
/v1/pet/{id}          # текущая версия
/v2/pet/{id}          # новая версия
```

#### 8. Outbox Pattern для надёжной доставки событий

Если используем Domain Events + Message Queue, нужен Outbox Pattern для гарантированной доставки.

## Changelog

### [2024-XX-XX] Рефакторинг архитектуры: Fat Services → CQRS

#### Добавлено
- Разделение Application слоя на `cmd/` (команды) и `query/` (запросы)
- Отдельные handlers для каждой операции:
  - Pet: `CreateHandler`, `UpdateHandler`, `DeleteHandler`, `RevalidateDonorHandler`, `GetByIDHandler`, `GetByUserHandler`
  - BloodSearch: 8 handlers (create, update, delete, apply, getByID, getByPetID, list, updateStatus)
  - File: `GetPresignedURLsHandler`, `ConfirmUploadHandler`
- Разделение `PetRepository` на специализированные интерфейсы:
  - `PetReadRepository` (GetByID, GetByUserID, Exists)
  - `PetWriteRepository` (Create, Update, Delete)
  - `PetStatsRepository` (CountSuitableDonors)
  - `PetPhotoRepository` (AddPhotoURLs)
- Доменные модели в `domain/bloodsearch/model/`:
  - `BloodRequest` с методами `IsActive()`, `Close()`, `ReserveVolume()`
  - `DonorResponse` со статусами
  - `BloodRequestFilter` для фильтрации

#### Изменено
- **BREAKING**: Интерфейсы репозиториев теперь работают с доменными моделями вместо ENT типов:
  - `Create(ctx, *ent.CreateBloodSearchRequestInput)` → `Create(ctx, *model.BloodRequest)`
  - `GetByID(ctx, id) (*ent.BloodSearchRequest, error)` → `GetByID(ctx, id) (*model.BloodRequest, error)`
- **BREAKING**: `BloodRequestRepository.Update` теперь принимает `*model.BloodRequest` вместо `*ent.UpdateBloodSearchRequestInput`
- Репозитории конвертируют ENT ↔ Domain через методы `toDomainModel()`
- HTTP handlers маппят DTO ↔ Domain models через `mapDTOToBloodRequest()`, `mapBloodRequestToDTO()`
- `UpdateRequestHandler` обновляет поля напрямую в доменной модели вместо использования `ent.UpdateBloodSearchRequestInput`
- `GetByIDHandler` и `GetByUserHandler` используют `bloodReq.ResponseIDs` вместо `bloodReq.Edges.Responses`

#### Исправлено
- Удалены дубликаты типов `DonorResponseStatus` и `DonorResponse` из `blood_request.go`
- Устранены циклические зависимости между domain слоями
- `LocationRepository` теперь возвращает `[]*model.Location` вместо `[]*ent.Location`
- Используется единый метод `storage.BuildPhotoURLs()` вместо дублирования логики в каждом handler

#### Удалено (Deprecated)
- Старые fat services закомментированы/не используются:
  - `domain/pet/pet_service.go`
  - `domain/bloodsearch/bloodrequest_service.go`
  - `domain/filestorage/file_service.go`

### [2024-XX-XX] Выделение Mapper Layer (DTO <-> Domain conversions)

#### Добавлено
- Новый пакет `internal/mapper/` для централизованного маппинга DTO <-> Domain:
  - `PetMapper`: `ToDTO()`, `ToDTOs()`, `ToDomainCreate()`, `ToDomainUpdate()`, `ToSimplifiedDTO()`
  - `UserMapper`: `ToDTO()`, `ToDTOs()` (использует PetMapper для вложенных питомцев)
  - `BloodRequestMapper`: `ToDTO()`, `ToDTOs()`, `ToDomain()`
- Мапперы инициализируются внутри HTTP handlers (self-contained)
- Удалены дублирующие функции маппинга из handlers:
  - `pet_handler_http.go`: удалены `ToDTO()`, `ToDTOs()`, `ToDomain()`, `calculateBirthDateFromAge()`, `nilable()`
  - `user_handler.go`: удалены `toDTO()`, `petToDTO()`
  - `blood_request_handler.go`: удалены `mapBloodRequestToDTO()`, `mapDTOToBloodRequest()`

#### Изменено
- HTTP handlers теперь используют мапперы через поле `*mapper.XXXMapper`
- Консистентный подход к маппингу во всех handlers

#### Технический долг (TODO)
- [ ] Удалить закомментированные старые сервисы полностью
- [ ] Разделить Pet агрегат на под-агрегаты (Health, Treatment, Analyses)
- [ ] Внедрить Value Objects (Weight, ChipNumber, PhotoGallery)
- [ ] Добавить Domain Events (PetBecameDonor, BloodRequestCreated)
- [ ] Реализовать Unit of Work для транзакций
- [ ] Вынести валидацию доноров в Specification Pattern
- [ ] Создать Read Models для оптимизации N+1 запросов

## Контакты и ресурсы

- **Технологии**: Go, ENT, Huma, PostgreSQL
- **Архитектура**: Clean Architecture, CQRS, Repository Pattern
- **База данных**: PostgreSQL 15+ с soft delete (deleted_at)
- **Хранение файлов**: S3-совместимое (MinIO)