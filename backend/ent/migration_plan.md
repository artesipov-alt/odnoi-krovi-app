# План миграции и разработки схем ENT

Этот файл является основным руководством и трекером прогресса для разработки базы данных. Каждый раз после внесения изменений в схемы или выполнения этапа миграции, необходимо обновлять этот файл.

## 🛠 Общие правила разработки (Guidelines)

При создании или изменении любой схемы ENT необходимо строго придерживаться следующих правил:

1.  **JSON Аннотации**: Все поля должны иметь тег `StructTag` с именованием в стиле **camelCase** для фронтенд-разработчиков.
    *   *Пример:* `field.String("full_name").StructTag("json:\"fullName\"")`
2.  **Идентификация (ID)**:
    *   Для сущностей, требующих уникальный строковый ID, использовать `NewBaseMixin(Prefix)`.
    *   Префиксы должны быть определены в `common.go` (например, `USR`, `PET`).
3.  **Аудит и Мягкое удаление**:
    *   Всегда подключать `AuditMixin` или `BaseMixin`.
    *   Это обеспечивает наличие полей `created_at`, `updated_at`, `deleted_at` и автоматическую фильтрацию удаленных записей.
4.  **Валидация**:
    *   Использовать `MaxLen()` для строковых полей, где это уместно.
    *   Использовать `Optional()` и `Nillable()` для полей, которые могут быть пустыми в БД и должны быть указателями в Go.
5.  **Диагностика**:
    *   Перед фиксацией любых изменений в схемах **обязательно** запускать генерацию кода (`go generate ./ent`) для проверки целостности и отсутствия ошибок.

---

## 📋 Статус разработки схем

### 1. Базовые компоненты (`common.go`)
- [x] Определение префиксов ID
- [x] Реализация генератора NanoID
- [x] `AuditMixin` (Soft Delete + Timestamps)
- [x] `BaseMixin` (ID + Audit)
- [x] Интерцептор для фильтрации `deleted_at`

### 2. Сущности (Schemas)
- [x] **User**
    - [x] Поля (telegram_id, role, etc.)
    - [x] Mixin (Base)
    - [x] Связи (pets)
- [x] **Pet**
    - [x] Поля (name, type, status, etc.)
    - [x] Mixin (Base)
    - [x] Связи (owner, health, treatments, analyses, bonuses)
- [x] **PetHealth**
    - [x] Поля
    - [x] Mixin (Audit)
    - [x] Связи (pet)
- [x] **PetTreatment**
    - [x] Поля
    - [x] Mixin (Audit)
    - [x] Связи (pet)
- [x] **PetAnalysis**
    - [x] Поля
    - [x] Mixin (Audit)
    - [x] Связи (pet)
- [x] **PetBonus**
    - [x] Поля
    - [x] Mixin (Audit)
    - [x] Связи (pet)
- [x] **Location**
    - [x] Поля (name)
    - [x] Mixin (Audit)
    - [x] Связи (users)
- [x] **Breed**
    - [x] Поля (name, type)
    - [x] Mixin (Audit)
    - [x] Связи (pets)
- [x] **BloodGroup**
    - [x] Поля (pet_type, blood_group, description)
    - [x] Mixin (Audit)
    - [x] Связи (search_requests)
- [x] **BloodComponent**
    - [x] Поля (name)
    - [x] Mixin (Audit)
    - [x] Связи (search_requests)
- [x] **BloodSearchRequest**
    - [x] Поля (blood_volume, regions, photo_urls, etc.)
    - [x] Mixin (Base)
    - [x] Связи (pet [1:1], blood_components, blood_group)

---

## 🚀 План действий (Next Steps)

1.  [ ] **Проверка текущего состояния**: Запустить `go generate ./ent` для подтверждения валидности текущих схем.
2.  [ ] **Расширение схем**: (Добавить новые сущности по мере необходимости, например, `Donation`, `Request`).
3.  [ ] **Миграция**: Подготовка и запуск миграции в БД.
4.  [ ] **Тестирование**: Написание тестов для проверки Soft Delete и генерации ID.

---
*Последнее обновление: 24.05.2024*