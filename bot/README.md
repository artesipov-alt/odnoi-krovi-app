# 🤖 Telegram Bot - Одной Крови

<div align="center">

![Bun](https://img.shields.io/badge/Bun-1.3+-000000?style=for-the-badge&logo=bun&logoColor=white) ![TypeScript](https://img.shields.io/badge/TypeScript-5.0+-3178C6?style=for-the-badge&logo=typescript&logoColor=white) ![Grammy](https://img.shields.io/badge/Grammy-1.38+-00ADD8?style=for-the-badge&logo=telegram&logoColor=white)

</div>

**Telegram Bot для платформы "Одной Крови"** — помощник для поиска донорской крови для животных и управления профилем пользователя через Telegram.

## 🎯 Функциональность

### Основные возможности:
- **📱 Интеграция с Mini App** - Быстрый доступ к основному приложению
- **👤 Управление профилем** - Просмотр и редактирование данных пользователя
- **🔍 Поиск доноров** - Быстрый поиск через inline режим (в разработке)
- **🔔 Уведомления** - Оповещения о новых запросах крови
- **📊 Статистика** - Отслеживание активности пользователей

### Команды бота:
- `/start` - Запуск бота и основное меню
- `/help` - Справка по использованию
- `/profile` - Просмотр профиля пользователя

## 🏗️ Архитектура

```
bot/
├── src/
│   ├── handlers/              # Обработчики сообщений и команд
│   │   ├── commands.ts        # Основные команды бота
│   │   └── errors.ts          # Обработка ошибок
│   ├── middleware/            # Промежуточное ПО
│   │   ├── logger.ts          # Логирование запросов
│   │   └── ratelimitter.ts    # Ограничение запросов
│   ├── config/                # Конфигурация
│   │   └── templates.ts       # Шаблоны сообщений
│   ├── utils/                 # Вспомогательные функции
│   └── instances.ts           # Инициализация бота и логгера
├── bot.ts                     # Основной файл бота
├── build.ts                   # Скрипт сборки
├── pm2.config.cjs             # Конфигурация для PM2
└── package.json
```

## 🚀 Быстрый старт

### Предварительные требования

- **Bun 1.3+** - JavaScript runtime
- **Telegram Bot Token** от @BotFather
- **OpenAPI Generator** - Для генерации TypeScript типов из API
- **Node.js 18+** (опционально, для PM2)

### Установка и запуск

#### 1. Установка Bun

##### macOS

**Вариант 1: Через Homebrew**
```bash
brew install oven-sh/bun/bun
```

**Вариант 2: Через curl**
```bash
curl -fsSL https://bun.sh/install | bash
```

##### Windows
```bash
powershell -c "irm bun.sh/install.ps1 | iex"
```

##### Linux
```bash
curl -fsSL https://bun.sh/install | bash
```

Проверьте установку:
```bash
bun --version
```

#### 2. Установка OpenAPI Generator

OpenAPI Generator нужен для автоматической генерации TypeScript типов из спецификации Backend API. Это избавляет от необходимости вручную писать типы и гарантирует их синхронизацию с Backend.

##### macOS
```bash
brew install openapi-generator
```

##### Windows
1. Установите Java (требуется для OpenAPI Generator)
2. Скачайте JAR файл с [официального сайта](https://openapi-generator.tech/docs/installation)
3. Или используйте npm:
```bash
npm install -g @openapitools/openapi-generator-cli
```

##### Linux (Ubuntu/Debian)
```bash
# Через npm (рекомендуется)
npm install -g @openapitools/openapi-generator-cli

# Или через snap
sudo snap install openapi-generator
```

Проверьте установку:
```bash
openapi-generator-cli version
```

#### 3. Клонирование репозитория (если еще не сделали)
```bash
git clone https://github.com/your-username/odnoi-krovi-app.git
cd odnoi-krovi-app
```

#### 4. Настройка переменных окружения

> ⚠️ **Важно:** Перед запуском бота убедитесь, что вы создали файл `.env` в корневой директории проекта согласно [основной инструкции](../README.md#2-настройка-переменных-окружения).

Файл `.env` должен находиться в `odnoi-krovi-app/.env` (не в `odnoi-krovi-app/bot/.env`), так как все сервисы используют единый конфигурационный файл.

#### 5. Переход в директорию bot

```bash
cd bot
```

#### 6. Установка зависимостей
```bash
bun install
```

#### 7. Генерация TypeScript типов из API

После того как Backend запущен и сгенерировал Swagger документацию, выполните:

**Вариант 1: Через npm script (рекомендуется)**
```bash
bun run generate-api
```

**Вариант 2: Полная команда вручную**

Если хотите больше контроля или нужно настроить генерацию:

```bash
openapi-generator generate \
  -i ../backend/docs/swagger.json \
  -g typescript-fetch \
  -o src/api
```

Или с дополнительными параметрами:

```bash
openapi-generator generate \
  -i ../backend/docs/swagger.json \
  -g typescript-fetch \
  -o src/api \
  --additional-properties=supportsES6=true,npmName=odnoi-krovi-api,npmVersion=1.0.0
```

**Параметры команды:**
- `-i` - путь к OpenAPI спецификации (swagger.json или swagger.yaml)
- `-g` - генератор (typescript-fetch, typescript-axios, typescript-node и др.)
- `-o` - папка для сохранения сгенерированных файлов
- `--additional-properties` - дополнительные настройки генерации

**Что делает команда:**
- Читает OpenAPI спецификацию из Backend
- Генерирует TypeScript типы, интерфейсы и API клиент
- Сохраняет сгенерированные файлы в указанную папку
- Эти типы могут использоваться как в Bot, так и в **Frontend приложении**

**Когда нужно запускать генерацию:**
- После первого клонирования проекта
- После обновления API на Backend
- После добавления новых эндпоинтов или изменения существующих

> 💡 **Совет для Frontend разработчиков:** Вы можете использовать ту же команду для генерации типов в вашем Frontend приложении, просто измените путь `-o` на нужную директорию в вашем проекте.

#### 8. Запуск в режиме разработки
```bash
bun run dev
```

#### 9. Сборка для production
```bash
bun run build
```

#### 10. Запуск собранной версии
```bash
bun run preview
```

## 🛠️ Технологический стек

### Основные технологии:
- **Bun** - Быстрый JavaScript runtime
- **TypeScript** - Статическая типизация
- **Grammy** - Фреймворк для Telegram Bot API
- **Pino** - Структурированное логирование

### Дополнительные пакеты:
- `@grammyjs/ratelimiter` - Ограничение запросов
- `@grammyjs/runner` - Параллельное выполнение
- `@grammyjs/transformer-throttler` - Ограничение API вызовов
- `pino-pretty` - Красивое форматирование логов


## 🚀 Production развертывание

### Использование PM2

```bash
# Установка PM2
npm install -g pm2

# Запуск бота
pm2 start pm2.config.cjs

# Мониторинг
pm2 dash

# Логи
pm2 logs odnoi-krovi-bot
```

## 🔒 Безопасность

### Rate Limiting
- Ограничение: 2 запроса в 1.5 секунды
- Хранилище: In-memory (планируется Redis)
- Сообщение при превышении лимита

### Валидация данных
- Проверка всех входящих сообщений
- Валидация callback данных
- Обработка ошибок на всех уровнях
