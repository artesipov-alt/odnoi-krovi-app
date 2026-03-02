# Архитектура проекта "Одной крови"

## Общее описание

"Одной крови" — платформа для поиска доноров крови для домашних животных. Регистрация питомцев, создание заявок, отклики доноров, управление статусами.

## Технологический стек

- **Go 1.21+**, **Huma v2** (HTTP), **ENT** (ORM), **PostgreSQL 15+**

## Структура проекта (Clean Architecture)

```
backend/
├── internal/
│   ├── domain/              # Чистый домен, без зависимостей
│   │   ├── pet/            # Агрегат Питомец
│   │   ├── user/           # Агрегат Пользователь
│   │   ├── bloodsearch/    # Агрегаты Заявка, Отклик
│   │   └── reference/      # Справочники
│   ├── application/        # Use Cases (CQRS)
│   │   └── {domain}/
│   │       ├── cmd/        # Команды (изменяют состояние)
│   │       └── query/      # Запросы (только чтение)
│   ├── transport/          # HTTP адаптеры (Huma)
│   │   ├── dto/            # Input/Output структуры
│   │   └── handlers/       # HTTP handlers
│   └── infra/              # Инфраструктура
│       └── presistance/pg/ # Реализация репозиториев
```

## Ключевые архитектурные принципы

### 1. Clean Architecture

Зависимости направлены внутрь: `transport → application → domain → (ничего)`

### 2. CQRS

- **Command** — изменяют состояние, возвращают минимум данных (ID, timestamps)
- **Query** — только чтение, возвращают полные данные

### 3. Агрегаты (Aggregates)

**Pet** — основной агрегат с Health/Treatments/Analyses (жизненный цикл связан).

**User**, **BloodRequest**, **DonorResponse** — отдельные агрегаты.

### 4. Repository Pattern

Интерфейсы в домене, реализация в инфраструктуре. Репозитории принимают и возвращают доменные модели.

### 5. Конструкторы с валидацией

Каждый агрегат создаётся через конструктор (`NewPet`, `NewUser`, `NewBloodRequest`, `NewDonorResponse`), который валидирует инварианты.

### 6. Маппинг DTO ↔ Domain

Только в `mapper/` пакете. Transport слой использует мапперы для конвертации DTO в доменные модели (через конструкторы).

## ⚠️ TODO

### 🔴 Высокий приоритет

1. **Update Pattern для User**
   - Сейчас: partial update через поля в handler
   - Нужно: `UpdateFrom()` метод в домене (как у Pet)

2. **Update Pattern для BloodRequest**
   - Добавить `UpdateFrom()` в доменную модель

3. **Update Pattern для DonorResponse**
   - Добавить `UpdateFrom()` в доменную модель

### 🟡 Средний приоритет

4. **Command Result структуры**
   - Сейчас: `CreateHandler.Handle()` возвращает `*model.Pet`
   - Нужно: Возвращать `CreatePetResult{ID, CreatedAt}` — минимум данных

5. **ReadModel/View для Query**
   - Сейчас: Query возвращает доменную модель `*model.Pet`
   - Нужно: Создать `PetView` с только нужными полями для чтения

6. **Value Objects**
   - Сейчас: `ChipNumber`, `Email`, `Phone` — простые строки
   - Нужно: Отдельные типы с валидацией в конструкторе

### 🟢 Низкий приоритет

7. **Domain Events**
   - `PetBecameDonorEvent`, `BloodRequestCreatedEvent`

8. **Unit of Work**
   - Для транзакций spanning multiple aggregates

9. **Read Models для сложных запросов**
   - Проекции для списков с фильтрацией

## Changelog

### [2024-XX-XX] Рефакторинг: Fat Services → CQRS + Clean Architecture

#### ✅ Добавлено
- Разделение на слои: domain → application → transport
- CQRS: Command и Query handlers
- Конструкторы с валидацией: `NewPet()`, `NewUser()`, `NewBloodRequest()`, `NewDonorResponse()`
- Разделение стоп-факторов на статические и динамические
- Метод `CalculateStatus()` для вычисления статуса питомца
- Явный маппинг без `copier.Copy`
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