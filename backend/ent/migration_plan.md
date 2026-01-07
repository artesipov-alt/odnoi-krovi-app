# План миграции с GORM на ENT

## Введение
Этот план является инструкцией для ИИ (искусственного интеллекта) по миграции моделей из GORM на ENT. ИИ должен следовать плану как блокноту: составлять план действий перед каждым шагом, выполнять шаги последовательно, проводить диагностику после каждого действия (используя инструмент diagnostics), исправлять ошибки на основе диагностики, и обновлять план с выполненными действиями. Всегда проверять совместимость с существующим кодом и зависимостями.

## Общий план
- Анализировать существующие модели GORM (user.go, pet.go).
- Создать схемы ENT для каждой сущности, используя правильные импорты и структуры ENT.
- Адаптировать генерацию ID: префикс + nanoID.
- Определить отношения с использованием edge.To и edge.From (не ent.HasMany/BelongsTo).
- Перенести soft delete, timestamps.
- Обновлять модели по одной, тестируя.

## Детали библиотеки ENT (важные уточнения для ИИ)
- **Импорты:** Всегда добавлять `"entgo.io/ent/schema/edge"` для использования edge.To, edge.From. Основные импорты: `"entgo.io/ent"`, `"entgo.io/ent/schema/field"`, `"entgo.io/ent/schema/mixin"`.
- **Поля:** 
  - Для float64 использовать `field.Float`, не `field.Float64`.
  - Для optional nullable полей (например, *time.Time) использовать `.Optional().Nillable()`, не просто `.Optional()`.
  - JSON-теги: всегда camelCase (например, `json:"fullName"`).
- **Отношения:** 
  - One-to-many: В owning схеме (User) использовать `edge.To("pets", Pet.Type)`. В inverse (Pet) использовать `edge.From("owner", User.Type).Ref("pets").Unique()`. ENT автоматически добавит FK поле (owner_id в Pet).
  - One-to-one: В owning схеме использовать `edge.To("health", PetHealth.Type)`. В inverse использовать `edge.From("pet", Pet.Type).Ref("health").Unique()`. ENT добавит FK (pet_id в PetHealth).
  - Не добавлять FK поля вручную в Fields — ENT делает это автоматически через edges.
- **Mixin:** Использовать для общих полей (ID, timestamps, deleted_at).
- **Enums:** Определять в field.Enum с .Values(...).
- **Диагностика:** После каждого изменения схем проводить diagnostics на всех файлах schema. Исправлять все ошибки и предупреждения перед генерацией.

## План действий (шаг за шагом)
ИИ должен выполнять шаги последовательно, отмечая выполненные, и проводить diagnostics после каждого.

1. Анализировать user.go и pet.go (структуры, поля, отношения).
2. Создать схему User в schema/user.go с BaseMixin, полями, edge.To("pets", Pet.Type).
3. Провести diagnostics на user.go.
4. Создать схему Pet в schema/pet.go с PetMixin, полями (без owner_id), edges (edge.From для owner, edge.To для health/treatments/analysis/bonuses).
5. Создать схемы PetHealth, PetTreatment, PetAnalysis, PetBonus с соответствующими полями и edge.From.
6. Провести diagnostics на всех schema файлах.
7. Исправить ошибки на основе diagnostics (добавить импорты, исправить поля, edges).
8. Провести diagnostics снова.
9. Сгенерировать код ENT.
10. Обновить модель User в коде (заменить GORM на ENT).
11. Тестировать создание/чтение User.
12. Обновить модель Pet и связанные.
13. Тестировать.
14. Удалить старые GORM файлы.

## Текущий прогресс
- Создан каталог backend/ent/.
- Анализированы user.go и pet.go.
- Решено: ID = префикс + nanoID (10 символов nanoID).
- Созданы схемы ENT для User, Pet, PetHealth, PetTreatment, PetAnalysis, PetBonus.
- Вынесены общие вещи: generateID в common.go, BaseMixin с NewBaseMixin(prefix).
- Обновлены схемы для использования общего BaseMixin.
- Исправлены ошибки: добавлены импорты edge, заменены ent.* на edge.*, убраны ручные FK, исправлены поля (Float, Nillable), добавлены Mixins для всех схем.
- Диагностика проведена — нет ошибок.
- Изменено отношение PetAnalysis на one-to-many (analyses), так как анализы — история; остальные (health, treatments, bonuses) — one-to-one.
- Добавлено поле deleted_at в BaseMixin для soft delete (реализуется вручную: Where(deleted_at.IsNil()) в запросах, Update().SetDeletedAt(time.Now()) для "удаления").
- Готово к генерации ENT (ent generate).

## Выполненные действия
- Анализ user.go и pet.go завершен.
- Создана схема ENT для User.
- Создана схема ENT для Pet и связанные (PetHealth, PetTreatment, PetAnalysis, PetBonus).
- Добавлены правильные edges (edge.To/edge.From).
- Исправлены типы полей и импорты.
- Вынесены общие вещи: создана common.go с generateID (10 символов nanoID) и NewBaseMixin(prefix).
- Обновлены все схемы для использования общего BaseMixin, убрано дублирование (PetMixin удален).
- Диагностика проведена на всех schema файлах — нет ошибок.
- Изменено отношение PetAnalysis на one-to-many для истории анализов.
- Добавлено поле deleted_at в BaseMixin для soft delete (реализуется вручную в коде).
- Готово к генерации ENT код (запустить ent generate вручную).

## Заметки
- Всегда использовать nanoID для ID (теперь 10 символов, префикс + nanoID).
- Общие поля (ID, timestamps, deleted_at) вынесены в BaseMixin с префиксами.
- Soft delete: Поле deleted_at добавлено; в коде для запросов добавлять .Where(deleted_at.IsNil()), для "удаления" — .Update().SetDeletedAt(time.Now()) вместо .Delete().
- Проверять после каждого шага: обновлять этот файл.
- JSON-теги: camelCase.
- После создания схем — ent generate (запустить вручную).
- Напоминание: Проводить diagnostics после каждого изменения!
- Отношения: PetAnalysis — one-to-many, остальные (PetHealth, PetTreatment, PetBonus) — one-to-one.
- Если ошибки в генерации — исправлять схемы и повторять diagnostics.
