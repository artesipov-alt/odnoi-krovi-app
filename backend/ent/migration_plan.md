# План миграции GORM -> ENT (Odnoi Krovi App)

## 🎯 Текущий статус
Мы успешно перенесли основные схемы базы данных на ENT и начали процесс переписывания репозиториев. Главное достижение — стабилизация генерации кода ENT через отказ от сложных миксинов в пользу явного объявления ID и внешних ключей.

---

## 🛠 Архитектурные решения (Guidelines для следующей сессии)

При работе с ENT в этом проекте **ОБЯЗАТЕЛЬНО** придерживаться следующих правил:

1.  **Явные ID (Explicit IDs)**:
    *   Для сущностей с NanoID (User, Pet, BloodSearchRequest) поле `id` прописывается вручную в `Fields()`.
    *   Использовать `DefaultFunc` с соответствующим префиксом из `common.go`.
    *   *Пример:* `field.String("id").DefaultFunc(func() string { return generateID(UserPrefix) })`

2.  **Явные внешние ключи (Explicit Foreign Keys)**:
    *   **НИКОГДА** не полагаться на магию ENT при создании связей.
    *   Всегда объявлять поле ключа в `Fields()` (например, `field.String("pet_id")` или `field.Int("location_id")`).
    *   В определении `Edges()` всегда указывать `.Field("field_name")`. Это критично для предотвращения ошибок несоответствия типов (`string` vs `int`).

3.  **Аудит и Soft Delete**:
    *   Поля `created_at`, `updated_at`, `deleted_at` и соответствующие `Interceptors` для фильтрации удаленных записей прописываются **явно в каждой схеме** (из-за капризов генератора при работе с миксинами в одном пакете).

4.  **CamelCase**:
    *   Все поля и связи должны иметь `StructTag` с `json:"camelCase"`.

---

## 📋 Статус компонентов

### 1. Схемы (Schemas) — [DONE]
- [x] **User**: NanoID, Audit, Связь с Location и Pets.
- [x] **Location**: Int ID, Audit, справочник городов.
- [x] **Breed**: Int ID, Audit, справочник пород.
- [x] **Pet**: NanoID, Audit. Включает в себя (в одном файле `pet.go`):
    - `PetHealth`, `PetTreatment`, `PetAnalysis`, `PetBonus` (все 1:1 к Pet).
- [x] **BloodGroup** & **BloodComponent**: Справочники.
- [x] **BloodSearchRequest**: NanoID, Audit, связь 1:1 к Pet (один питомец — одна активная заявка).

### 2. Инфраструктура — [DONE]
- [x] **Config**: `backend/internal/utils/config/ent_db.go` для подключения `ent.Client`.
- [x] **Enums**: `backend/internal/utils/enums/ent_enums.go` с типизацией ENT и локализацией на русский.
- [x] **Generate**: `backend/ent/generate.go` настроен.

### 3. Репозитории (Repositories) — [IN PROGRESS]
- [x] **EntUserRepository**: Реализован, проверен диагностикой. Поддерживает Soft Delete и Restore.
- [ ] **EntPetRepository**: Начата подготовка (требуется реализация транзакционного создания Pet + Health + Analysis).
- [ ] **EntLocationRepository**: Ожидает реализации.
- [ ] **EntBreedRepository**: Ожидает реализации.

---

## 🚀 Следующие шаги (Next Steps)

1.  **Завершить репозитории**:
    *   Реализовать `EntPetRepository` в `backend/internal/repositories/pg/pet_repo_ent.go`.
    *   Реализовать остальные репозитории (Location, Breed, Blood).
2.  **Интеграция в Service Layer**:
    *   Постепенно заменять `GormRepository` на `EntRepository` в конструкторах сервисов.
3.  **Миграция данных**:
    *   Подготовить скрипт переноса данных из старых таблиц GORM в новые таблицы ENT (учитывая NanoID).
4.  **Тестирование**:
    *   Проверить работу `Interceptors` для Soft Delete на реальных запросах.

---
*Обнял, ушел, но фундамент оставил бетонный. Удачи, бро!* 🐾