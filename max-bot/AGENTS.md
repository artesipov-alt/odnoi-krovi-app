# Max Bot — Telegram Bot для Max ecosystem

## Общий обзор

Telegram-бот для платформы «Одной Крови», работающий в экосистеме **Max** (через `@maxhub/max-bot-api`). Обрабатывает команды пользователей через Max Bot API и реагирует на события бэкенда через **Redis Pub/Sub**, отправляя уведомления пользователям о статусе донаций, запросах крови и контактах.

## Технологический стек

| Технология | Применение |
|---|---|
| **Bun** | Runtime и сборщик (Bun.build) |
| **TypeScript (ESNext)** | Язык |
| **@maxhub/max-bot-api** | Max Bot API SDK (бот, клавиатуры, контекст) |
| **Redis (ioredis)** | Pub/Sub шина событий из бэкенда |
| **Pino + pino-pretty** | Структурированное логирование |
| **Bun.serve** | HTTP-сервер для вебхуков Max |
| **shared/ts/** | OpenAPI-сгенерированный TS-клиент бэкенда |

> **Важно:** Бот использует `@maxhub/max-bot-api`, а не `grammy`. Пакет `grammy` присутствует в `package.json`, но не используется (legacy). Middleware `ratelimitter.ts` и `throttler.ts` — закомментированный legacy-код от grammy.

## Сборка и запуск

```bash
# dev (из корня монорепозитория)
bun --env-file=../.env bot.ts

# build
bun run build.ts   # → dist/bot.js

# prod
bun run dist/bot.js

# Docker (multi-stage)
docker build -f Dockerfile .
```

## Организация файлов

```
max-bot/
├── bot.ts                        # Entry point — инициализация, регистрация команд,
│                                 #   Redis-подписки, HTTP-сервер, graceful shutdown
├── build.ts                      # Bun.build скрипт (→ dist/)
├── Dockerfile                    # Multi-stage: builder + production alpine
├── pm2.config.cjs                # Legacy — не используется (устарел)
├── openapitools.json             # OpenAPI Generator CLI config (v7.17.0)
├── package.json                  # name: odnoi-krovi-max-bot
├── tsconfig.json
├── src/
│   ├── instances.ts              # Синглтоны: Bot, Pino logger, Redis, API-клиент
│   ├── keyboards.ts              # Inline-клавиатуры (главная, открыть приложение)
│   ├── max.ts                    # sendMessageToUser — безопасная отправка через Max API
│   ├── server.ts                 # Bun.serve HTTP-сервер (webhook endpoint)
│   ├── config/
│   │   └── templates.ts          # Шаблоны сообщений (start, help, profile)
│   ├── handlers/
│   │   ├── commands.ts           # Хендлеры: start, help, profile, apiTest, errCommandTest
│   │   └── errors.ts             # Глобальный error handler (MaxError)
│   ├── middleware/
│   │   ├── logger.ts             # Middleware логирования входящих обновлений
│   │   ├── ratelimitter.ts       # ⛔ Закомментирован (grammy legacy)
│   │   └── throttler.ts          # ⛔ Закомментирован (grammy legacy)
│   └── events/                   # Обработчики Redis Pub/Sub событий
│       ├── donor/
│       │   ├── donorCancel.ts         # Донор отказался от донации
│       │   ├── donorCompleted.ts      # Донор сообщил о завершении донации
│       │   ├── donorNotConfirmed.ts   # Реципиент не подтвердил донацию
│       │   ├── donorReject.ts         # Реципиент отклонил донацию
│       │   ├── recipientApply.ts      # Отклик реципиента на донора (принятие заявки)
│       │   └── helpers.ts             # Утилиты (generateMessage)
│       ├── recipient/
│       │   ├── donorApply.ts          # Отклик донора на реципиента
│       │   ├── donationConfirmed.ts   # Подтверждение донации (от реципиента донору)
│       │   ├── handleBloodRequestCreated.ts  # Новый запрос крови → уведомление донорам
│       │   └── helpers.ts             # Утилиты (генерация сообщений, VCF)
│       └── user/
│           └── handleUserContact.ts   # Передача контакта пользователя
```

## Архитектура: поток данных

### Команды (Max Bot API)

```
Пользователь → /start (или bot_started event)
    → bot.ts: bot.command("start", startHandler)
    → handlers/commands.ts:
        1. Извлечение maxId, fullName, payload
        2. Парсинг utm-меток из payload
        3. Аутентификация через AuthV1Api (shared/ts/) с таймаутом 5s
        4. Отправка приветствия с inline-клавиатурой
```

### События (Redis Pub/Sub)

```
Backend → Redis PUBLISH "prod:events" { "type": "donor_response_apply", "payload": {...}, "createdAt": "..." }
    → bot.ts: redis.on("message") → channelHandlers["prod:events"]
    → JSON.parse → EventEnvelope
    → dispatchEvent(envelope, eventHandlers)  // shared/ts/events.ts
    → src/events/{domain}/{handler}.ts
        1. Валидация ProviderMaxID
        2. sendMessageToUser() через Max Bot API
        3. При необходимости — contact attachment (VCF)
```

Все каналы имеют префикс окружения: `dev:` или `prod:`, определяемый из `Bun.env.ENV`. Диспатчер типов событий (`EventHandlerMap`) объявлен в `bot.ts` — добавление нового типа требует записи там, иначе компилятор отказывает.

### HTTP-сервер (вебхуки)

```
Max Platform → POST /webhook (X-Max-Bot-Api-Secret)
    → server.ts:
        1. Проверка WEBHOOK_SECRET
        2. bot.handleUpdate(update) — передача в Max Bot SDK
```

**Подписка на обновления.** В отличие от Telegram, Max API требует явной
регистрации вебхука через `POST https://platform-api2.max.ru/subscriptions`
с массивом `update_types`. Без `message_callback` в этом списке события
о нажатиях на inline-кнопки не доставляются боту.

Регистрация выполняется автоматически на старте через
`registerWebhook()` в `src/maxbot.ts`, если задан `MAX_BOT_WEBHOOK_URL`:

- prod: `https://prodbot.1krovi.app/webhook`
- dev: `https://devbot.1krovi.app/webhook`

Оба проксируются Caddy на `max-bot:6000` (см. `frontend/Caddyfile{,*.dev}`).

## Redis Pub/Sub каналы

События и уведомления приходят двумя каналами, префикс окружения (`dev:` / `prod:`) берётся из `Bun.env.ENV`:

| Канал | Содержимое | Диспатч |
|---|---|---|
| `{prefix}events` | `EventEnvelope` (см. `shared/ts/events.ts`) с полем `type` ∈ `EventType` | `dispatchEvent` в `bot.ts` по типу |
| `{prefix}notifications` | Объект уведомления для интерактивных кнопок | `handleNotification` |

Конкретные типы событий (`blood_request_created`, `donor_response_apply`, и т.д.) и их handler'ы зарегистрированы в `EventHandlerMap` в `bot.ts`. Добавление нового типа требует: (1) константу в `EVENT_TYPES`, (2) запись в `EventHandlerMap` — exhaustiveness-проверка типов поймает расхождения на этапе компиляции.

## API-клиент (shared/ts/)

- Использует сгенерированный OpenAPI TypeScript-клиент из `shared/ts/`.
- Инициализируется в `instances.ts` с таймаутом 5s и middleware для логирования.
- Сейчас используется только `AuthV1Api` (для аутентификации пользователей через `authUserViaService`).

## Важные решения и конвенции

1. **Max Bot API вместо grammy** — проект перешёл с grammy на `@maxhub/max-bot-api`. Код grammy (включая `pm2.config.cjs`, `ratelimitter.ts`, `throttler.ts`) оставлен как legacy reference и не используется.
2. **sendMessageToUser** (src/max.ts) — единая точка отправки сообщений с валидацией MaxID (числовой, не пустой). Всегда использовать её вместо прямого вызова `bot.api.sendMessageToUser`.
3. **Contact attachments** — для передачи контактов используется `generateVCF()` и `type: "contact"` attachment.
4. **Кнопка «Открыть приложение»** — во всех уведомлениях (кроме `handleUserContact`) добавляется `getAppOpenKeyboard()` из `src/keyboards.ts`. Это Link Button, открывающая Mini App Max по ссылке `https://max.ru/{BOT_ID}_bot?startapp`.
5. **Разделение контакта и кнопки** — в сценариях, где нужен и контакт, и кнопка (donorApply, donorNotConfirmed, handleUserContact), отправка разбивается на **два последовательных вызова** `sendMessageToUser`:
   - Сначала текст-уведомление + кнопка (`getAppOpenKeyboard()`)
   - Затем пустое сообщение + contact attachment (VCF)
   Это сделано, чтобы избежать проблем с отображением в Max.
9. **In-memory кэш токенов (`authStore`)** — JWT токены кэшируются в `Map<MaxId, AuthData>` для избежания повторной аутентификации при каждом callback-нажатии. Токен живёт 24ч, проверка — с запасом 5 минут. При 401 ответе от бэкенда — `invalidateToken()` + retry.
10. **Callback-кнопки в уведомлениях** — `callbackData` передаёт action и requestId через формат `notification_{yes|no}_{requestId}`. Обработка через `bot.action(regex, handler)`.
11. **action вместо callbackQuery** — в Max API используется `bot.action(triggers, handler)`, а не `bot.callbackQuery()`. Тип контекста — `FilteredContext<Ctx, 'message_callback'>`. `ctx.match` — `RegExpExecArray`.
12. **answerOnCallback** — метод `ctx.answerOnCallback({ notification: string, message?: ... })`. Нет `show_alert`. notification — текст всплывающего уведомления.
13. **ctx.user вместо ctx.from** — нет `ctx.from`. Пользователь: `ctx.user` с полями `user_id`, `first_name`, `last_name`, `username`.
14. **editMessage()** — в Max API `ctx.editMessage({ text, format, attachments })`, а не `ctx.editMessageText(text, extra)`. `attachments` = массив клавиатур. Для удаления кнопок — `attachments: [Keyboard.inlineKeyboard([])]` (пустой массив `attachments` не снимает уже отрисованные кнопки).