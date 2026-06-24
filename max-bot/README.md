# 🤖 Telegram Bot - Одной Крови

<div align="center">

![Bun](https://img.shields.io/badge/Bun-1.3+-000000?style=for-the-badge&logo=bun&logoColor=white) ![TypeScript](https://img.shields.io/badge/TypeScript-5.0+-3178C6?style=for-the-badge&logo=typescript&logoColor=white) ![Grammy](https://img.shields.io/badge/Grammy-1.38+-00ADD8?style=for-the-badge&logo=telegram&logoColor=white)

</div>

**Telegram Bot для платформы "Одной Крови"** — помощник для поиска донорской крови для животных и управления профилем пользователя через Telegram.

## 🎯 Функциональность

### Основные возможности:
- **📱 Интеграция с Mini App** — быстрый доступ к основному приложению
- **👤 Управление профилем** — просмотр и редактирование данных пользователя
- **🔍 Поиск доноров** — быстрый поиск через inline режим (в разработке)
- **🔔 Уведомления** — оповещения о новых запросах крови
- **📊 Статистика** — отслеживание активности пользователей

### Команды бота:
- `/start` — запуск бота и основное меню
- `/help` — справка по использованию
- `/profile` — просмотр профиля пользователя

## 🏗️ Архитектура

```odnoi-krovi-app/bot/README.md#L241-280
bot/
├── src/
│   ├── handlers/              # Обработчики сообщений и команд
│   │   ├── commands.ts        # Основные команды бота
│   │   └── errors.ts          # Обработка ошибок
│   ├── middleware/            # Промежуточное ПО
│   │   ├── logger.ts          # Логирование запросов
│   │   └── ratelimiter.ts     # Ограничение запросов
│   ├── config/                # Конфигурация
│   │   └── templates.ts       # Шаблоны сообщений
│   ├── utils/                 # Вспомогательные функции
│   └── instances.ts           # Инициализация бота и логгера
├── bot.ts                     # Основной файл бота
├── build.ts                   # Скрипт сборки
├── package.json
└── Dockerfile                 # Образ для docker-compose
```

> Примечание: `pm2.config.cjs` больше не используется — проект перешёл на запуск через `docker-compose` и `Taskfile.yaml` (корневой `Taskfile.yaml` содержит задачи для установки, разработки и сборки всех сервисов, включая bot).

## 🚀 Быстрый старт

### Предварительные требования
- **Bun 1.3+** (локально для разработки)
- **Docker & docker-compose** (для локального продакшн-подобного запуска)
- **Task CLI (go-task)** — рекомендуется для запуска задач из `Taskfile.yaml`. Инструкции по установке: `docs/taskfile-install.md`.
- **OpenAPI Generator** (опционально) — если нужно вручную генерировать TS типы из backend Swagger (альтернативно используйте `task generate-api`).
- **Файл .env** в корне проекта (в корне монорепозитория — `odnoi-krovi-app/.env`).

### Установка и запуск (варианты)

1) Быстрый способ — через Taskfile (рекомендуется)
```odnoi-krovi-app/bot/README.md#L281-320
# В корне проекта
task install-bot        # установит зависимости bot (выполняет `bun install` в папке bot)
task generate-api       # сгенерирует shared типы из backend/swagger и положит в shared/
task dev-bot            # запустит bot в режиме разработки (bun run dev)
# Или для запуска всех сервисов:
task dev                # запустит backend, frontend и bot (см. Taskfile.yaml)
```

2) Локальная разработка вручную
```odnoi-krovi-app/bot/README.md#L321-360
# 1) Клонируем и переходим в корень репозитория
git clone git@github.com:artesipov-alt/odnoi-krovi-app.git
cd odnoi-krovi-app

# 2) Убедитесь, что в корне есть .env (все сервисы читают его)
# 3) Перейдите в директорию бота
cd bot

# 4) Установите зависимости (Bun)
bun install

# 5) Сгенерируйте API типы (если backend запущен и swagger доступен)
# рекомендуется запускать из корня:
task generate-api
# или вручную:
openapi-generator generate -i ../backend/docs/swagger.json -g typescript-fetch -o ../shared

# 6) Запуск в dev режиме
bun run dev
```

3) Запуск через Docker Compose (локально/простое продакшн-окружение)
```odnoi-krovi-app/bot/README.md#L361-420
# В корне проекта:
# Собрать образ bot и запустить сервисы (docker-compose.yml в корне репозитория)
docker-compose build bot
docker-compose up -d bot

# Запустить весь стек
docker-compose up -d

# Просмотр логов бота
docker-compose logs -f bot

# Остановить
docker-compose stop bot
# Или полностью:
docker-compose down
```

## 📦 Генерация типов API (shared)
- В этом репозитории есть папка `shared/` — автосгенерированные TypeScript типы, используемые frontend и bot.
- Рекомендуется запускать `task generate-api` из корня репозитория. В `Taskfile.yaml` для `generate-api` уже прописан порядок: сначала `generate-swagger`, затем `openapi-generator generate ... -o shared`.

## 🛠️ Команды для разработчика

```odnoi-krovi-app/bot/README.md#L421-480
# Просмотр задач Taskfile
task -l

# Установка зависимостей
task install-bot   # через Taskfile (рекомендуется)
# или локально:
cd bot && bun install

# Разработка
task dev-bot       # через Taskfile (рекомендуется)
# или локально:
cd bot && bun run dev

# Сборка для production
task build-bot
# или локально:
cd bot && bun run build

# Запуск через docker-compose
docker-compose build bot
docker-compose up -d bot
```

## 🔒 Безопасность и конфигурация
- Все секреты (BOT_TOKEN, PROD_BOT_API_KEY и пр.) хранятся в `.env` в корне репозитория. Никогда не коммитьте `.env`.
- При запуске в контейнерах Docker убедитесь, что `docker-compose.yml` передаёт переменные окружения в контейнер бота.

## 🔧 Технологический стек
- Bun — runtime и менеджер пакетов
- TypeScript — типизация
- Grammy — фреймворк для Telegram Bot API
- Pino — структурированное логирование

## 🚩 Примечания
- PM2 больше не используется в репозитории — переход к docker-compose даёт воспроизводимое окружение и более простую интеграцию с CI/CD.
- Для удобства локальной разработки используйте `task` (см. `docs/taskfile-install.md`), он оркестрирует установку, генерацию типов и запуск сервисов.

**Сделано с ❤️ для наших четвероногих друзей** 🐾
