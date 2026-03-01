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

## Структура проекта (Clean Architecture)

```
backend/
├── internal/
│   ├── domain/              # Чистый домен, без зависимостей
│   │   ├── pet/            # Агрегат Питомец
│   │   │   ├── model/      # Доменные модели (entities, value objects)
│   │   │   ├── pet_repo.go # Интерфейсы репозиториев
│   │   ├── user/           # Агрегат Пользователь
│   │   ├── bloodsearch/    # Агрегаты Заявка на кровь, Отклик донора
│   │   └── reference/      # Справочники (породы, группы крови)
│   ├── application/        # Use Cases (CQRS)
│   │   ├── pet/
│   │   │   ├── cmd/        # Команды (изменяют состояние)
│   │   │   │   ├── create.go
│   │   │   │   ├── update.go
│   │   │   │   └── ...
│   │   │   └── query/      # Запросы (только чтение)
│   │   │       ├── get_by_id.go
│   │   │       └── ...
│   ├── transport/          # Адаптеры внешнего мира
│   │   └── http/
│   │       ├── dto/        # DTO для Huma (Input/Output)
│   │       └── handlers/   # HTTP handlers
│   └── infra/              # Инфраструктура
│       └── presistance/
│           └── pg/         # Реализация репозиториев на PostgreSQL/Ent
└── docs/                   # OpenAPI спецификации
```

## Ключевые архитектурные принципы

### 1. Чистая архитектура (Clean Architecture)

**Зависимости направлены внутрь:**
```
transport → application → domain → (ничего)
```

- **Domain** — не зависит ни от чего, содержит бизнес-логику
- **Application** — зависит только от domain, оркестрирует use cases
- **Transport** — зависит от application, преобразует HTTP в команды/запросы
- **Infra** — реализует интерфейсы домена (репозитории)

### 2. CQRS (Command Query Responsibility Segregation)

Разделение операций:
- **Command** — изменяют состояние, возвращают минимум данных (ID, CreatedAt/UpdatedAt)
- **Query** — только читают, не изменяют состояние, возвращают полные данные

### 3. Агрегаты (Aggregates)

**Pet (Питомец)** — основной агрегат, содержит:
- Сущности (Entities): `PetHealth`, `PetTreatment`, `PetAnalysis`
- Value Objects: `FactorCode`, `PetStatus`, `PetType`, `Gender`
- Инварианты: стоп-факторы вычисляются в домене

**Почему Pet — один агрегат с Health/Treatments/Analyses?**

1. **Жизненный цикл**: Health, Treatments, Analyses не существуют без Pet
2. **Загрузка через Ent edges**: Используем `PetPreloadOptions` для eager loading
3. **Целостность**: При сохранении Pet в одной транзакции сохраняются все связанные сущности
4. **Инварианты**: Статус донора зависит от всех этих данных

```go
type Pet struct {
    ID      string
    Name    string
    Health  *PetHealth        // Entity
    Treatments *PetTreatment  // Entity  
    Analyses []*PetAnalysis   // Entities
    // ... другие поля
}
```

**User** — отдельный агрегат

**BloodRequest + DonorResponse** — отдельные агрегаты (связаны через ID)

### 4. Repository Pattern

Интерфейсы объявлены в домене, реализация в инфраструктуре:

```go
// domain/pet/pet_repo.go
type PetReadRepository interface {
    GetByID(ctx context.Context, id string, opts PetPreloadOptions) (*model.Pet, error)
    Exists(ctx context.Context, id string) (bool, error)
}

type PetWriteRepository interface {
    Create(ctx context.Context, pet *model.Pet) (*model.Pet, error)  // Возвращает полный агрегат!
    Update(ctx context.Context, id string, pet *model.Pet) (*model.Pet, error)
    Delete(ctx context.Context, id string) error
}
```

### 5. Конструкторы с валидацией

Каждый агрегат создаётся только через конструктор:

```go
func NewPet(
    name string,
    petType PetType,
    weightKg float64,
    gender Gender,
    ownerID string,
    // ... остальные поля
) (*Pet, error) {
    // Валидация здесь!
    if name == "" {
        return nil, errors.New("pet name is required")
    }
    // ...
}
```

### 6. Маппинг DTO ↔ Domain

Маппинг происходит только в:
- **transport/http/dto/** — DTO структуры с Huma-тегами
- **mapper/** — конвертация DTO ↔ Domain

Запрещено использовать доменные модели напрямую в transport!

## ✅ Что реализовано

### 1. Разделение стоп-факторов

**Статические факторы** (хранятся в БД):
- `StopFactorNoPhoto`, `StopFactorHasDiseases`, `StopFactorTransfused`
- `StopFactorNoRabiesVaccination`, `StopFactorNoInfectionVaccination`
- Метод: `GetStaticStopFactors()`

**Динамические факторы** (вычисляются на лету):
- `StopFactorTooOld`, `StopFactorTooYoung` — возраст
- `StopFactorVaccinationExpired`, `StopFactorVaccinationTooRecent`
- `StopFactorDonationTooRecent` — срок после донации
- Метод: `GetDynamicStopFactors(now)`

**Проверка при наличии активной заявки:**
- `ActualStopFactors(now, hasActiveRequest)` — объединяет все факторы

### 2. Статус питомца вычисляется, не хранится

```go
func (p *Pet) CalculateStatus(now time.Time, hasActiveRequest bool, hasResponses bool) PetStatus {
    if hasActiveRequest && hasResponses {
        return PetStatusBloodFound
    }
    if hasActiveRequest {
        return PetStatusRecipient
    }
    if len(p.ActualStopFactors(now, false)) > 0 {
        return PetStatusNone
    }
    return PetStatusDonor
}
```

### 3. Явный маппинг без рефлексии

- ✅ Убран `copier.Copy` полностью
- ✅ Убраны `SetInput()` — только явные `Set()` методы
- ✅ Явная функция `petToDomain()` для конвертации Ent → Domain

### 4. DTO Input/Output структуры

Каждая операция имеет свои типы:
- `CreatePetInput` / `CreatePetOutput` (возвращает ID + CreatedAt)
- `UpdatePetInput` / `UpdatePetOutput` (возвращает ID + UpdatedAt)
- `GetPetByIDInput` / `GetPetByIDOutput` (возвращает полные данные)

### 5. CQRS на уровне Application

```
application/pet/
├── cmd/           # Команды
│   ├── create.go
│   ├── update.go
│   ├── delete.go
│   └── revalidate_donor.go
└── query/         # Запросы
    ├── get_by_id.go
    └── get_by_user.go
```

### 6. Справочники — только Query, без Command

```
application/reference/cmd/ — пустая директория ✅
application/reference/query/ — только чтение
```

## ⚠️ Архитектурный технический долг

### 🔴 Высокий приоритет (сделать сейчас)

1. **Доделать DTO рефакторинг**
   - User DTO — структурировать на Input/Output
   - BloodSearch DTO — структурировать на Input/Output
   - Убрать дублирование общих типов

2. **Исправить ошибки компиляции**
   - Текущие проблемы с типами после рефакторинга

### 🟡 Средний приоритет (пропускаем сейчас, нужно будет сделать)

3. **Command Result структуры** ⏸️
   - Сейчас: `CreateHandler.Handle()` возвращает `*model.Pet`
   - Нужно: Возвращать `CreatePetResult{ID, CreatedAt}` — минимум данных

4. **ReadModel/View для Query** ⏸️
   - Сейчас: Query возвращает доменную модель `*model.Pet`
   - Нужно: Создать `PetView` с только нужными полями для чтения

5. **Value Objects** ⏸️
   - Сейчас: `ChipNumber`, `Email`, `Phone` — простые строки
   - Нужно: Отдельные типы с валидацией в конструкторе

### 🟢 Низкий приоритет (по возможности)

6. **Domain Events**
   - `PetBecameDonorEvent`, `BloodRequestCreatedEvent`

7. **Unit of Work**
   - Для транзакций spanning multiple aggregates

8. **Read Models для сложных запросов**
   - Проекции для списков с фильтрацией

## Changelog

### [2024-XX-XX] Рефакторинг: Fat Services → CQRS + Clean Architecture

#### ✅ Добавлено
- Разделение на слои: domain → application → transport
- CQRS: Command и Query handlers
- Конструкторы с валидацией: `NewPet()`, `NewUser()`, `NewDonorResponse()`
- Разделение стоп-факторов на статические и динамические
- Метод `CalculateStatus()` для вычисления статуса питомца
- Явный маппинг `petToDomain()` без `copier.Copy`
- DTO Input/Output структуры для всех операций

#### ✅ Изменено
- Репозиторий возвращает полный агрегат (не частичный)
- `Create()` и `Update()` перечитывают созданный/обновлённый объект из БД
- Application handlers принимают доменную модель, не Ent-структуры
- Статус питомца вычисляется при каждом запросе, не хранится в БД

#### ✅ Удалено
- `copier.Copy` — полностью убрана зависимость
- `PetService` — логика перенесена в Command/Query handlers
- Анонимные структуры в хендлерах

## Контакты

Проект: "Одной крови"
Архитектура: Clean Architecture + CQRS