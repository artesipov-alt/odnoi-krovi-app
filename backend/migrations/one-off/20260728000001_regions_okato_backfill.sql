-- One-off: переход на коды ОКАТО на существующих БД.
--
-- Заменяет старые строковые ID "MSK" -> "77" (г. Москва), "MO" -> "50" (Московская область)
-- во всех местах, где они хранятся:
--   - blood_requests.regions (jsonb array)
--   - donor_preferences.preferred_location_ids (jsonb array)
--   - users.location_id (FK -> ref_locations.id)
-- и удаляет устаревшие записи ref_locations с id IN ('MSK','MO').
--
-- На свежей БД не нужен: справочник сразу создаётся с ОКАТО-кодами
-- (bootstrap/20260728000001_ref_locations.sql).
-- Применяется runner-ом: task db:migrate (см. migrations/README.md).
-- Idempotent: повторный запуск безопасен (CASE-обновления не меняют
-- уже обновлённые значения).
--
-- Применять в одной транзакции (BEGIN/COMMIT). На prod — после бэкапа.

BEGIN;

-- 1. Обновление jsonb-массива blood_requests.regions.
--    Заменяет элементы "MSK" -> "77", "MO" -> "50".
UPDATE blood_requests
SET regions = (
  SELECT jsonb_agg(CASE
    WHEN elem = 'MSK' THEN '77'
    WHEN elem = 'MO'  THEN '50'
    ELSE elem
  END)
  FROM jsonb_array_elements_text(regions) AS elem
)
WHERE regions IS NOT NULL
  AND (regions::text LIKE '%MSK%' OR regions::text LIKE '%"MO"%');

-- 2. Обновление jsonb-массива donor_preferences.preferred_location_ids.
UPDATE donor_preferences
SET preferred_location_ids = (
  SELECT jsonb_agg(CASE
    WHEN elem = 'MSK' THEN '77'
    WHEN elem = 'MO'  THEN '50'
    ELSE elem
  END)
  FROM jsonb_array_elements_text(preferred_location_ids) AS elem
)
WHERE preferred_location_ids IS NOT NULL
  AND (preferred_location_ids::text LIKE '%MSK%' OR preferred_location_ids::text LIKE '%"MO"%');

-- 3. Обновление users.location_id (FK -> ref_locations.id).
--    Должно идти ДО удаления старых записей ref_locations, иначе FK violation.
UPDATE users
SET location_id = CASE location_id
    WHEN 'MSK' THEN '77'
    WHEN 'MO'  THEN '50'
END
WHERE location_id IN ('MSK', 'MO');

-- 4. Удаление устаревших записей справочника.
DELETE FROM ref_locations WHERE id IN ('MSK', 'MO');

COMMIT;
