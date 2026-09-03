-- Бэкфилл users.verified_at для уже верифицированных пользователей.
--
-- Поле verified_at добавлено в 3.29.0 для уведомления «добавьте питомца»
-- (тип user_verified_no_pets, через 24ч после верификации). Для пользователей,
-- верифицированных до появления поля, лучшее доступное приближение момента
-- верификации — created_at.
--
-- Idempotent: обновляются только строки с verified_at IS NULL.
-- Мягко удалённые пользователи не обновляются (запрос уведомлений их и так не выберет).
BEGIN;

UPDATE users
SET verified_at = created_at
WHERE verified = true
  AND verified_at IS NULL
  AND deleted_at IS NULL;

COMMIT;
