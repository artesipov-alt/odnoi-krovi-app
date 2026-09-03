-- Дроп старых полных UNIQUE-индексов users.phone / users.email.
--
-- В 3.28.1 Unique() убран с полей phone и email: вместо полных индексов
-- используются partial unique индексы user_phone / user_email
-- (WHERE deleted_at IS NULL), которые Ent auto-migrate создаёт при старте.
-- Старые полные индексы users_phone_key / users_email_key Ent не дропает,
-- и пока они существуют, запись с телефоном/почтой умягко удалённого
-- пользователя продолжает блокировать повторную регистрацию с теми же
-- phone/email (duplicate key value violates unique constraint).
--
-- Применять ТОЛЬКО после деплоя 3.28.1+ и auto-migrate, когда partial-индексы
-- user_phone / user_email уже созданы.
--
-- Idempotent: DROP INDEX IF EXISTS.
BEGIN;

-- Освобождаем телефоны/почты мягко удалённых пользователей
DROP INDEX IF EXISTS users_phone_key;
DROP INDEX IF EXISTS users_email_key;

COMMIT;
