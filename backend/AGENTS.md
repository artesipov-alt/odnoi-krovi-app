# Backend — Go-монолит

## Общий обзор

Go-монолит, реализованный в стилистике **DDD (Domain-Driven Design)** с элементами **CQRS** на уровне приложения. Веб-фреймворк — **Huma v2** (REST + OpenAPI 3.1). ORM — **Ent**. База — **PostgreSQL**. Кеш/события — **Redis**. Файлы — **S3**. Аутентификация — **JWT** + Telegram Mini App data validation.

## Слои и организация пакетов

```
backend/
├── cmd/                     # Точки входа
│   ├── api/main.go         # Основной API-сервер (DI-композиция)
│   └── dburl/main.go       # Утилита: печатает DSN для psql по ENV (для Taskfile)
├── internal/
│   ├── apperrors/           # Единая система ошибок (AppError) + интеграция с Huma
│   ├── domain/              # DOMAIN LAYER — бизнес-логика и модели
│   │   ├── common/          # Общие value objects (PetType, Компоненты крови, и т.д.)
│   │   ├── ports/           # Порт-интерфейсы внешнего мира (EventPublisher)
│   │   ├── {bounded-context}/
│   │   │   ├── model/       # Aggregate root + value objects (чистые Go-структуры)
│   │   │   ├── events/      # Domain events (структуры данных)
│   │   │   ├── *repo.go     # Repository interface (порты для persistence)
│   │   │   └── *_service.go # Stateless domain services
│   ├── application/         # APPLICATION LAYER — CQRS обработчики
│   │   ├── {bounded-context}/
│   │   │   ├── cmd/         # Command handlers (изменяют состояние)
│   │   │   └── query/       # Query handlers (только чтение)
│   ├── infra/               # INFRASTRUCTURE LAYER — адаптеры
│   │   ├── ent/             # Ent ORM: schema/ + generated code
│   │   ├── presistance/
│   │   │   ├── pg/          # PostgreSQL реализации репозиториев (Ent)
│   │   │   ├── domainmapper/# Мапперы: Ent entity → Domain model
│   │   │   ├── redis/       # Redis cache repository
│   │   │   ├── cache/       # Cache interface + ключи
│   │   │   ├── s3/          # S3 file storage repository
│   │   │   └── tx_manager.go # Транзакционный менеджер (через context)
│   │   └── events/redis/    # Event publisher (Redis Pub/Sub)
│   └── transport/           # TRANSPORT LAYER — HTTP/GRPC адаптеры
│       ├── http/
│       │   ├── *handler.go    # Регистрация маршрутов Huma + делегирование Cmd/Query
│       │   ├── dto/           # Request/Response DTO (Huma-аннотированные)
│       │   ├── dtomapper/     # Мапперы: Domain model → DTO
│       │   └── middleware/    # Auth, Recovery, CORS, Tracing
│       └── grpc/              # (пока пусто)
├── pkg/                    # SHARED KERNEL — переиспользуемые утилиты
│   ├── auth/               # JWT генерация/валидация + Telegram InitData проверка
│   ├── config/             # Server, DB (Ent), Redis, CORS конфигурация
│   ├── enums/              # Сгенерированные Ent enum-константы
│   └── logger/             # Настройка slog + Charm Bracelet
├── migrations/             # SQL-миграции данных (справочники, одноразовые преобразования)
├── docs/
│   └── openapi.json        # Сгенерированная OpenAPI 3.1 спецификация
└── docsui/                 # Scalar docs UI встраивание
```

## Миграции базы данных

Разделены два слоя (см. [migrations/README.md](migrations/README.md)):

- **Schema migrations (DDL)** — Ent auto-migrate при старте приложения (`config.RunMigrations`). Создаёт/изменяет таблицы, индексы, колонки.
- **Data migrations (SQL)** — файлы в `migrations/`, применяются вручную через `task db:migrate` (ENV=local|dev|prod). Справочники (`ref_locations`, `ref_breeds`) и одноразовые преобразования данных.

Seeds в коде не используются — заменены idempotent SQL-миграциями.

## Bounded Contexts (Domain)

| Контекст | Модель | Репозиторий (интерфейс) | Domain Service | Команды | Запросы |
|---|---|---|---|---|---|
| **User** | User, DonorPreference | UserRepository | — | register, update, delete, reset, restore | get_by_id, get_contact, get_deleted |
| **Auth** | Identity (value object) | — (через UserRepo) | AuthService | external_sign_in, mini_app_sign_in | — |
| **Pet** | Pet, PetHealth, PetTreatment, PetAnalysis | PetReadRepository + PetWriteRepository | PetService | create, update, delete, revalidate_donor | get_by_id, get_by_user |
| **BloodSearch** | BloodRequest | Repository | MatchingService | create, update, delete, select_donor, accept_response, confirm_donation, reject_donation, close_request | get_by_id, get_by_pet_id, get_donor_by_id, get_donation |
| **Donor** | DonorResponse | DonorResponseRepository | — | apply, complete_donation, cancel_donation | list_requests, get_recipient, get_planned, get_completed, get_bonuses |
| **Bonus** | Bonus (value object) | BonusRepository | BonusService | import_bonuses | — |
| **FileStorage** | — | FileRepository | FileService | get_presigned_urls, confirm_upload | — |
| **Partner** | Partner | PartnerRepository | — | — | — |
| **Reference** | Breed, Location | BreedRepository + LocationRepository | — | — | get_all_breeds, get_by_type, get_locations, get_blood_components, get_blood_groups |

## DDD: как реализовано

- **Domain Model** (`internal/domain/{ctx}/model/`) — чистые Go-структуры без тегов ORM, с методами-конструкторами (`NewPet(...)`), методами поведения (`RecalculateFactors(...)`, `UpdateFrom(...)`).
- **Repository Interface** (`internal/domain/{ctx}/*_repo.go`) — порты для persistence, разделены на write-only и read-only (см. PetWriteRepository / PetReadRepository).
- **Domain Service** (`*_service.go`) — stateless, содержит логику, требующую координации нескольких aggregate (например, `PetService.CalculateAndSetStatus` оперирует Pet + DonorResponse + BloodRequest).
- **Domain Events** (`internal/domain/{ctx}/events/`) — структуры данных событий (BloodRequestCreated, ApplyDonor, DonorSelected, DonationConfirmed, DonorCompleted и т.д.).
- **Ports** (`internal/domain/ports/`) — интерфейсы для внешних систем (EventPublisher).

## CQRS: как реализовано

**CQRS-lite** — разделены команды (изменяющие состояние) и запросы (только чтение) на уровне application layer.

- **Command handlers** (`cmd/`) — принимают входные данные, валидируют, вызывают domain-логику, сохраняют через репозиторий, публикуют события.
- **Query handlers** (`query/`) — только читают данные, возвращают DTO или domain models. Никаких side effects.
- Каждый handler — это struct с единственным методом `Handle(ctx, ...)`.
- Repository interface разделены на **write** и **read** (например, `pet.PetWriteRepository` vs `pet.PetReadRepository`), чтобы на уровне типов гарантировать, что query handler не может случайно вызвать метод записи.

## Поток данных (пример: создание питомца)

```
HTTP POST /api/v1/pets
    → middleware: Auth → проверяет JWT, кладет user_id в context
    → pet_handler_http.go:Parse DTO → petCommand := cmd.CreateHandler
        → CreateHandler.Handle(ctx, userID, petModel)
            → petRepo.Create(ctx, pet)   // EntPetRepository
                → domainmapper.PetToDomain(entPet) // Ent → Domain
            → domain events (при необходимости)
        → dtomapper.PetToDTO(pet)  // Domain → HTTP DTO
    → 201 JSON response
```

## Ключевые технологии

| Технология | Применение |
|---|---|
| **Go 1.26** | Язык |
| **Huma v2** | REST API + OpenAPI 3.1 генерация |
| **Ent** | ORM (code-first схемы, генерация типов) |
| **PostgreSQL** | Основная база данных |
| **Redis** | Кеш, Pub/Sub для событий |
| **S3 (MinIO/Cloud)** | Хранение файлов (pre-signed URLs) |
| **JWT (HS256)** | Аутентификация |
| **Telegram Mini App** | Валидация init data |
| **Scalar** | Swagger UI |

## Аутентификация и middleware (порядок)

1. `Recovery` — восстановление после паники
2. `CORS` — разрешение origin'ов
3. `BasicAuth` — защита `/docs` и `/openapi.json`
4. `Auth` — JWT-валидация (bearer token), пропускает исключённые пути
5. `Logging` — slog-http (только статусы ≥400)
6. `TraceID` — X-Trace-Id в ответ

## Обработка ошибок

Единый тип `apperrors.AppError` реализует `huma.StatusError`. Все ошибки проходят через кастомный `huma.NewError`, что даёт консистентный JSON ответ:

```json
{
  "Code": "NOT_FOUND",
  "Message": "питомец не найден",
  "Details": {"pet_id": "..."},
  "HTTPStatus": 404
}
```

## Транзакции

`presistance.TxManager` оборачивает бизнес-логику в транзакцию Ent. Транзакция передаётся через контекст, что позволяет прозрачно использовать один и тот же репозиторий внутри и вне транзакции.

## События (Domain Events → Redis)

- Domain events определяются в `internal/domain/{ctx}/events/` как plain structs.
- `EventPublisher` (интерфейс в `ports/`) публикует их в Redis Pub/Sub.
- Реализация: `internal/infra/events/redis/event_publisher.go`.
- При недоступности Redis используется `NoOpEventPublisher`.

## Важные решения

1. **Soft Delete** — реализован на уровне Ent interceptors, прозрачен для domain.
2. **Domain Mapper** — ENT сущности маппятся в domain model через `domainmapper/*.go`. Domain model никогда не зависит от ORM.
3. **DTO Mapper** — domain model маппится в HTTP DTO через `dtomapper/*.go` на транспортном уровне.
4. **Preload опции** — query handlers принимают `PreloadOptions`, позволяя клиенту гибко запрашивать связанные данные.
5. **Infra/events/redis** отделён от **infra/presistance/redis** — событийная шина и кеш путать не стоит.