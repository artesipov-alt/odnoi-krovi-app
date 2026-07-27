# Telegram Bot — Telegram Bot для «Одной Крови»

## Общий обзор

Telegram-бот для платформы «Одной Крови». Обрабатывает команды пользователей через **Grammy** (Telegram Bot API) в режиме long-polling и реагирует на события бэкенда через **Redis Pub/Sub**, отправляя уведомления о статусе донаций, запросах крови и контактах. Соединяется с Telegram через **bridge API** (`bridge.1krovi.app`).

## Технологический стек

| Технология | Применение |
|---|---|
| **Bun** | Runtime и сборщик (Bun.build) |
| **TypeScript (ESNext)** | Язык |
| **grammy** | Telegram Bot API SDK |
| **@grammyjs/runner** | Конкурентный long-polling (запуск `run(bot)`) |
| **@grammyjs/ratelimiter** | Rate limiting (активен, in-memory) |
| **@grammyjs/transformer-throttler** | API throttler (импортирован, но middleware не подключена — `bot.api.config.use()` без аргументов) |
| **Redis (ioredis)** | Pub/Sub шина событий из бэкенда |
| **Pino + pino-pretty** | Структурированное логирование |
| **shared/ts/** | OpenAPI-сгенерированный TS-клиент бэкенда |

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
tg-bot/
├── bot.ts                        # Entry point — инициализация, регистрация команд,
│                                 #   Redis-подписки, runner (long-polling), graceful shutdown
├── build.ts                      # Bun.build скрипт (→ dist/)
├── Dockerfile                    # Multi-stage: builder + production alpine
├── pm2.config.cjs                # Legacy — не используется (устарел)
├── openapitools.json             # OpenAPI Generator CLI config (v7.17.0)
├── package.json                  # name: odnoi-krovi-tg-bot
├── tsconfig.json
├── src/
│   ├── instances.ts              # Синглтоны: Bot (grammy), Pino logger, Redis, API-клиент
│   ├── telegram.ts               # sendTelegramMessage, sendTelegramContact — безопасная отправка
│   ├── telegramButtons.ts        # createOpenAppKeyboard, createUserChatButton — клавиатуры для уведомлений
│   ├── config/
│   │   └── templates.ts          # Шаблоны сообщений (start, help, profile)
│   ├── handlers/
│   │   ├── commands.ts           # Хендлеры: start, help, profile, apiTest, errCommandTest
│   │   └── errors.ts             # Глобальный error handler (GrammyError, HttpError, BotError)
│   ├── middleware/
│   │   ├── logger.ts             # Middleware логирования входящих обновлений
│   │   ├── ratelimitter.ts       # Активен — rate limit 2 запроса / 1.5s in-memory
│   │   └── throttler.ts          # API throttler (импортирован, middleware не подключена)
│   └── events/                   # Обработчики Redis Pub/Sub событий
│       ├── donor/
│       │   ├── donorCancel.ts         # Донор отказался от донации
│       │   ├── donorCompleted.ts      # Донор сообщил о завершении донации
│       │   ├── donorNotConfirmed.ts   # Реципиент не подтвердил донацию
│       │   ├── donorReject.ts         # Реципиент отклонил донацию
│       │   ├── donorSelected.ts       # Реципиент выбрал донора (контакты + уведомление обоим)
│       │   ├── helpers.ts             # Утилиты (generateMessage)
│       │   └── recipientApply.ts      # Отклик реципиента на донора (принятие заявки)
│       ├── notification/
│       │   ├── handleNotification.ts       # Обработка уведомлений с кнопками
│       │   └── handleNotificationRespond.ts # Обработка ответов на уведомления
│       ├── recipient/
│       │   ├── donorApply.ts          # Отклик донора на реципиента
│       │   ├── donationConfirmed.ts   # Подтверждение донации (от реципиента донору)
│       │   ├── handleBloodRequestCreated.ts  # Новый запрос крови → уведомление донорам
│       │   └── helpers.ts             # Утилиты (генерация сообщений)
│       └── user/
│           └── handleUserContact.ts   # Передача контакта пользователя
```

## Архитектура: поток данных

### Команды (Telegram Bot API — long polling)

```
Пользователь → /start
    → bot.ts: bot.command("start", startHandler) + runner.run(bot)
    → handlers/commands.ts:
        1. Извлечение telegramId, fullName, payload
        2. Парсинг utm-меток из payload
        3. Аутентификация через AuthV1Api (shared/ts/) с providerName: "telegram_bot"
        4. Отправка приветствия с InlineKeyboard (webApp → MINIAPP_DOMAIN)
```

### События (Redis Pub/Sub)

```
Backend → Redis PUBLISH "prod:events" { "type": "donor_response_apply", "payload": {...}, "createdAt": "..." }
    → bot.ts: redis.on("message") → channelHandlers["prod:events"]
    → JSON.parse → EventEnvelope
    → dispatchEvent(envelope, eventHandlers)  // shared/ts/events.ts
    → src/events/{domain}/{handler}.ts
        1. Валидация ProviderTelegram / ProviderTelegramID
        2. sendTelegramMessage() или sendTelegramContact() через Grammy API
```

Все каналы имеют префикс окружения: `dev:` или `prod:`, определяемый из `Bun.env.ENV`. Диспатчер типов событий (`EventHandlerMap`) объявлен в `bot.ts` — добавление нового типа требует записи там, иначе компилятор отказывает.

## Redis Pub/Sub каналы

События и уведомления приходят двумя каналами, префикс окружения (`dev:` / `prod:`) берётся из `Bun.env.ENV`:

| Канал | Содержимое | Диспатч |
|---|---|---|
| `{prefix}events` | `EventEnvelope` (см. `shared/ts/events.ts`) с полем `type` ∈ `EventType` | `dispatchEvent` в `bot.ts` по типу |
| `{prefix}notifications` | Объект уведомления для интерактивных кнопок | `handleNotification` |

Конкретные типы событий (`blood_request_created`, `donor_response_apply`, и т.д.) и их handler'ы зарегистрированы в `EventHandlerMap` в `bot.ts`. Добавление нового типа требует: (1) константу в `EVENT_TYPES`, (2) запись в `EventHandlerMap` — exhaustiveness-проверка типов поймает расхождения на этапе компиляции.

## Ключевые отличия от max-bot

| Аспект | tg-bot | max-bot |
|---|---|---|
| **Фреймворк** | grammy | @maxhub/max-bot-api |
| **Режим получения обновлений** | Long-polling через `@grammyjs/runner` | Webhook через Bun.serve |
| **Передача контактов** | `sendContact()` (нативный Telegram) | Contact attachment + VCF |
| **Кнопка «Открыть приложение»** | `createOpenAppKeyboard()` в `telegramButtons.ts` → `InlineKeyboard().webApp()` | `getAppOpenKeyboard()` в `keyboards.ts` → `Keyboard.button.link()` |
| **Bridge** | Подключение через `bridge.1krovi.app` с секретом | Нет, прямой вызов Max Bot API |
| **Provider name** | `telegram_bot` | `max_bot` |
| **ID-поля в событиях** | `ProviderTelegram`, `ProviderTelegramID` | `ProviderMaxID`, `ProviderMaxID` |
| **Rate limiting** | Активен (`@grammyjs/ratelimiter`) | Нет (legacy, закомментирован) |
| **API throttler** | Импортирован, но не подключен | Нет |
| **HTTP сервер** | Нет | Есть (Bun.serve для вебхуков) |

## API-клиент (shared/ts/)

- Использует сгенерированный OpenAPI TypeScript-клиент из `shared/ts/`.
- Инициализируется в `instances.ts` с таймаутом 5s и middleware для логирования.
- Сейчас используется только `AuthV1Api` (для аутентификации пользователей через `authUserViaService`).

## Важные решения и конвенции

1. **Bridge API** — бот подключается к Telegram не напрямую, а через прокси-сервер `bridge.1krovi.app`, добавляя `?secret=BRIDGE_TOKEN` к каждому запросу (см. `instances.ts`, `buildUrl`).
2. **sendTelegramMessage / sendTelegramContact** (src/telegram.ts) — единые точки отправки с валидацией chat_id. Всегда использовать их вместо прямого вызова `bot.api.*`.
3. **Контакты через sendContact** — для передачи контактов используется нативный Telegram `sendContact`, а не VCF-вложение (в отличие от max-bot).
4. **Кнопка «Открыть приложение»** — единая функция `createOpenAppKeyboard()` в `src/telegramButtons.ts`. Используется во **всех** уведомлениях, кроме `handleUserContact`. URL выбирается по `Bun.env.ENV`:
   - dev → `https://dev.1krovi.app`
   - prod → `https://1krovi.app`
   Кнопка передаётся как `reply_markup` в `sendTelegramMessage`.
5. **In-memory кэш токенов (`authStore`)** — JWT токены кэшируются в `Map<telegramId, AuthData>` для избежания повторной аутентификации при каждом callback-нажатии. Токен живёт 24ч, проверка — с запасом 5 минут. При 401 ответе от бэкенда — `invalidateToken()` + retry с новым токеном.
6. **Callback-кнопки в уведомлениях** — `callback_data` передаёт action и requestId через формат `notification_{yes|no}_{requestId}`. Обработка через `bot.callbackQuery()` с regex-матчингом. После обработки — `ctx.answerCallbackQuery()` (подавление прелоадера).
