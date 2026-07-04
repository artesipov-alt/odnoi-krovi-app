# Changelog

Все заметные изменения в **max-bot** фиксируются в этом файле.

Формат основан на [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
и проект следует [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.9.3] — 2026-07-04

### Fixed

- **Ошибка при нажатии «Да» в уведомлении `recipient_empty_showcase`.** Max API отклонял `editMessage` с ошибкой `Field 'buttons' size (0) must be at least 1`, потому что `Keyboard.inlineKeyboard([])` создаёт клавиатуру с пустым массивом `buttons`. Исправлено: передаём `attachments: null` при редактировании сообщения, что снимает inline-кнопки.

## [0.9.2] — 2026-07-04

### Fixed

- **TLS-валидация `platform-api2.max.ru` падала с `UNABLE_TO_GET_ISSUER_CERT_LOCALLY`.**
  В `certs/` лежали только **промежуточные** CA Минцифры (`*.cer`), но не
  было **корневого**. Без корневого CA цепочка доверия не строится,
  даже если все промежуточные CA доверены системе. Проблема **не связана**
  с Bun: TLS-стек Bun и `NODE_EXTRA_CA_CERTS` работают штатно при наличии
  полной цепочки CA. Добавлен `russian_trusted_root_ca.crt` в
  `max-bot/certs/`. В `Dockerfile` копирование сертификатов расширено:
  `*.crt` (root) копируется напрямую, `*.cer` (sub) переименовываются
  в `.crt` для совместимости с `update-ca-certificates` Alpine.


## [0.9.1] — 2026-07-04

### Fixed

- **Inline-кнопки в Max не реагировали на нажатия.** Бот не получал
  `message_callback` от Max API, потому что вебхук не был подписан на этот
  тип обновлений. Добавлена автоматическая регистрация вебхука на старте
  через `registerWebhook()` в `src/maxbot.ts` (с `update_types`, включающим
  `message_callback` и остальные нужные типы). Управляется переменной
  окружения `MAX_BOT_WEBHOOK_URL` (для prod и dev прописана в
  `docker-compose.yml` / `docker-compose.dev.yml`).
- **Кнопки «Да»/«Нет» не убирались после нажатия «Да».** В Max API
  `attachments: []` при `editMessage` не сбрасывает уже отрисованную
  inline-клавиатуру — нужно явно передать пустую:
  `attachments: [Keyboard.inlineKeyboard([])]`.
- **`UNABLE_TO_GET_ISSUER_CERT_LOCALLY` при `POST /subscriptions`.**
  Домен `platform-api2.max.ru` подписан промежуточным CA Минцифры,
  которого нет в стандартном `ca-certificates` Alpine. В runtime-стадию
  `Dockerfile` добавлены `certs/*.cer` (Russian Trusted Sub CA) +
  `NODE_EXTRA_CA_CERTS` как страховка от собственного CA-bundle Bun.
- **Миграция Max API на `platform-api2.max.ru` (дедлайн 19.07.2026).**
  SDK `@maxhub/max-bot-api` по умолчанию ходит на старый домен
  `platform-api.max.ru`, который отключат. В `src/instances.ts` конструктору
  `Bot` теперь передаётся `clientOptions.baseUrl =
  "https://platform-api2.max.ru"`. Это влияет на все API-вызовы:
  `sendMessage`, `sendMessageToUser`, `editMessage`, `deleteMessage`,
  `getMyInfo`, `setMyCommands`, `answerOnCallback` и т.д.
  Константа `MAX_API_BASE_URL` экспортируется и переиспользуется в
  `src/maxbot.ts` для прямого `fetch` на `/subscriptions`.

## [0.9.0] — 2026-07-03

### Changed

- **Рефакторинг подписки на Redis-каналы:**
  - Вместо 10 отдельных каналов бот подписывается только на два: `events` и `notifications`.
  - Для канала `events` добавлен парсинг `EventEnvelope` и диспатч по полю `type`.
  - Все event handler'ы обновлены: поля переведены на camelCase в соответствии с новыми JSON-тегами от бэкенда.
  - `handleBloodRequestCreated` переписан: вместо массива `AvilableDonors` — один донор с полем `maxId`.

## [0.8.0] — 2026-07-02

### Added

- **Уведомление `recipient_empty_showcase`** — если у реципиента пустая витрина и он не заходил 24ч, приходит сообщение с кнопками «Да» / «Нет». Ответ обрабатывается через `POST /v1/blood-request/notification/respond` (с авторизацией через кэшированный JWT).
- **Уведомление `recipient_search_closed_inactive`** — если реципиент не заходил 48ч и не нажал «Да», поиск закрывается, приходит уведомление с кнопкой «Открыть приложение».
- **In-memory кэш токенов (`authStore`)** — `getOrCreateToken` кэширует JWT в Map, избегая повторной аутентификации при каждом нажатии callback-кнопки. Retry при 401.
- **Callback-обработчик `notification_yes` / `notification_no`** — обрабатывает нажатия на кнопки в уведомлениях с автоподстановкой Authorization header.

## [0.7.0] — 2026-07-01

### Added

- **Подписка на канал `notifications`** — новый обработчик `handleNotification` слушает Redis-канал `{env:}notifications`.
  Поддерживаемые типы уведомлений:
  - `recipient_donor_waiting` — донор откликнулся на запрос реципиента
  - `recipient_inactive_warning` — реципиент не заходил 6ч, есть невыбранные доноры
  - `donor_not_accepted` — донору о том, что реципиент не принял предложение

## [0.6.0] — 2026-06-24

### Added

- **Кнопка «🩸 Открыть приложение»** добавлена на все уведомления (7 событий):
  - Отклик донора на реципиента (`donorApply`)
  - Отклик реципиента на донора (`recipientApply`)
  - Отмена донации донором (`donorCancel`)
  - Отклонение донации реципиентом (`donorReject`)
  - Неподтверждение донации (`donorNotConfirmed`)
  - Завершение донации донором (`donorCompleted`)
  - Передача контакта пользователя (`handleUserContact`)
- Используется существующая функция `getAppOpenKeyboard()` из `src/keyboards.ts`.

### Changed

- **Разделение контакта и кнопки** — в сценариях, где уведомление содержит и contact-attachment, и кнопку, отправка разбита на два последовательных вызова `sendMessageToUser`:
  1. Текст-уведомление + кнопка открытия приложения
  2. Пустое сообщение + contact (VCF)
  Затронутые файлы: `donorApply.ts`, `donorNotConfirmed.ts`, `handleUserContact.ts`.
