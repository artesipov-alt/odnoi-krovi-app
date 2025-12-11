# 🩸 Одной Крови - Платформа донорства крови для животных

<div align="center">

![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go) ![React](https://img.shields.io/badge/React-18+-61DAFB?style=for-the-badge&logo=react&logoColor=black) ![TypeScript](https://img.shields.io/badge/TypeScript-5.0+-3178C6?style=for-the-badge&logo=typescript&logoColor=white)
![Node.js](https://img.shields.io/badge/Node.js-18+-339933?style=for-the-badge&logo=nodedotjs&logoColor=white) ![Bun](https://img.shields.io/badge/Bun-1.3+-000000?style=for-the-badge&logo=bun&logoColor=white) ![Telegram](https://img.shields.io/badge/Telegram-MiniApp-26A5E4?style=for-the-badge&logo=telegram&logoColor=white) ![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16+-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)
![Swagger](https://img.shields.io/badge/Swagger-85EA2D?style=for-the-badge&logo=Swagger&logoColor=black)

</div>

**Одной Крови** — это Telegram Mini App, который помогает находить донорскую кровь для животных и позволяет владельцам питомцев становиться донорами вместе со своими любимцами.

## 🗺️ Дорожная карта

Ознакомьтесь с [дорожной картой проекта](ROADMAP.md) для понимания планов развития и текущих этапов разработки.
Инструкцию по установке Task CLI можно найти в [docs/taskfile-install.md](docs/taskfile-install.md).

## 🎯 Цель проекта

Создать удобную и безопасную платформу в формате Telegram Mini App для:
- 🔍 Поиска доноров крови для животных в экстренных ситуациях
- ❤️ Регистрации питомцев как потенциальных доноров
- 📍 Геолокации ближайших доноров через Telegram Web App
- 🔔 Системы уведомлений через Telegram Bot API
- 💬 Интеграции с Telegram для удобного общения между пользователями

## 🏗️ Архитектура проекта

Проект построен как монорепозиторий с разделением на backend, bot и frontend - Telegram Mini App:

```
odnoi-krovi-app/
├── .github/         # Настрокий авторазвертывания (CD/CI)
├── backend/         # Go API сервер (Fiber + GORM + Swagger)
├── frontend/        # Telegram Mini App (React + TypeScript)
├── bot/             # Telegram Bot (Bun + Grammy)
├── microservices/   # Микросервисы
│   └── blood-microservice/  # Микросервис управления пулом поиска крови
├── shared/          # Автосгенерированные TypeScript типы c бэкенда (Swagger -> types)
├── docs/            # Документация
└── README.md
```

## 🚀 Быстрый старт

### Предварительные требования

- **Go 1.25+** для backend и микросервисов
- **Node.js 18+** и **npm** для Telegram Mini App
- **Bun 1.3+** для бота (используется в `bot/`)
- **PostgreSQL 16+** для баз данных (основная + bloodsearch)
- **Telegram Bot Token** от @BotFather
- **Docker** (опционально, для разработки)

### Установка и запуск

#### 1. Клонирование репозитория
```bash
git clone git@github.com:artesipov-alt/odnoi-krovi-app.git
cd odnoi-krovi-app
```

#### 2. Настройка переменных окружения

Создайте файл `.env` в **корневой директории проекта** со следующими переменными:

```env
# Database Configuration (Main Backend)
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=odnoi_krovi
DB_SSLMODE=disable

# Database Configuration (Blood Search Microservice)
DB_NAME_BLOOD_SEARCH=bloodsearch

# Server Configuration (Backend)
SERVER_PORT=3000
ENVIRONMENT=development

# Microservices URLs
BLOOD_MICROSERVICE_URL=http://localhost:8081

# S3 Storage Configuration
S3_ENDPOINT=your_s3_endpoint
S3_ACCESS_KEY=your_s3_access_key
S3_SECRET_KEY=your_s3_secret_key
S3_BUCKET_NAME=your_s3_bucket_name
S3_REGION=ru-central1

# Telegram Bot Configuration
BOT_TOKEN=your_telegram_bot_token
PROD_BOT_API_KEY=your_prod_bot_token
MINIAPP_DOMAIN=https://your-miniapp-domain.com

# API Configuration
API_BASE_URL=http://localhost:3000/api/v1

# Logging
LOG_LEVEL=info
APP_ENV=development
```

**Получение Telegram Bot Token:**
1. Напишите [@BotFather](https://t.me/botfather) в Telegram
2. Отправьте команду `/newbot`
3. Следуйте инструкциям для создания бота
4. Скопируйте полученный токен в `BOT_TOKEN` и `PROD_BOT_API_KEY`

**Важные замечания:**
- 📁 Файл `.env` должен находиться в корне проекта: `odnoi-krovi-app/.env`
- 🔒 Все сервисы (backend, bot, frontend) читают переменные из этого единого файла
- ⚠️ Никогда не коммитьте `.env` файл в git (он уже добавлен в `.gitignore`)
- 📋 Для production используйте отдельный `.env` с реальными учетными данными

#### 3. Установка компонентов

Следуйте инструкциям по установке для каждого компонента:
- [Backend Setup](backend/README.md) - Go API сервер
- [Blood Microservice Setup](microservices/blood-microservice/README.md) - Микросервис поиска крови
- [Bot Setup](bot/README.md) - Telegram Bot
- [Frontend Setup](frontend/README.md) - Telegram Mini App

### Taskfile (быстрый запуск задач)
В корне проекта есть `Taskfile.yaml` (формат Task v3). Это удобный способ запускать часто используемые команды и цепочки задач.

- Просмотреть список доступных задач:
  - `task -l` или `task --list`
- Запустить задачу:
  - `task <название_задачи>` (например `task generate-api` или `task dev`)
- Примеры:
```bash
task -l
task generate-api
task dev
```
(Для работы Taskfile требуется [установленный CLI](docs/taskfile-install.md))

#### 4. Запуск через Docker (рекомендуется)
```bash
# Запуск всех сервисов
docker-compose up -d

# Или только нужных сервисов
docker-compose up backend blood-microservice -d
```

## 📁 Структура проекта

### 🚀 Backend (Go + Fiber + GORM + Swagger)
Мощный API сервер с современным стеком технологий. Подробное описание архитектуры и возможностей доступно в [документации backend-сервера](backend/README.md).

### 🩸 Blood Microservice (Go + Connect + GORM)
Микросервис для управления пулом поиска крови. Обрабатывает запросы на поиск доноров, управляет заявками и статусами. Использует Connect протокол для межсервисного взаимодействия. Подробнее в [документации микросервиса](microservices/blood-microservice/README.md).

### 💻 Telegram Mini App (React + TypeScript + Telegram Web App SDK)
Интуитивный интерфейс для пользователей с полной интеграцией в Telegram. Узнайте больше о фронтенд-архитектуре в [документации по фронтенду](frontend/README.md).

### 🤖 Telegram Bot (Bun + Grammy)
Быстрый бот для уведомлений и коммуникации между пользователями. Полное описание функциональности и настройки смотрите в [документации бота](bot/README.md).

### 📦 Shared
- `shared/` — содержит автосгенерированные TypeScript типы и утилиты, сгенерированные из Swagger (используется фронтом и ботом для согласованных типов).

## 🛠️ Технологический стек

### Backend
- **Go 1.25+** - Основной язык программирования
- **Fiber** - Быстрый веб-фреймворк
- **GORM** - ORM для работы с базой данных
- **PostgreSQL** - Основная база данных
- **Redis** - Кэширование и сессии (планируется)
- **Validator** - Валидация данных
- **Swagger** - Документация API (swaggo/swag)

### Microservices
- **Go 1.25+** - Основной язык программирования
- **Connect** - RPC фреймворк (вместо gRPC)
- **Chi** - HTTP роутер
- **GORM** - ORM для работы с базой данных
- **PostgreSQL** - База данных микросервиса
- **Zap** - Структурированное логирование
- **Buf** - Генерация кода из protobuf

### Telegram Mini App
- **React 18+** - UI библиотека
- **TypeScript** - Статическая типизация
- **Vite** - Сборщик и dev сервер
- **Telegram Web App SDK** - Интеграция с Telegram
- **React Query** - Управление состоянием сервера
- **React Hook Form** - Управление формами
- **Axios** - HTTP клиент

### Telegram Bot
- **TypeScript** - Статическая типизация
- **Bun** - Серверная платформа (современная альтернатива Node.js и npm)
- **Grammy** - Фреймворк для Telegram Bot API

## 📊 Модели данных (в работе)

### Основные сущности:
- **Пользователи** - Владельцы животных (с привязкой к Telegram ID)
- **Животные** - Питомцы (доноры/реципиенты)
- **Запросы на кровь** - Заявки на поиск доноров
- **Донорские записи** - История донаций
- **Ветеринарные клиники** - Партнерские организации

## 🔗 Интеграция с Telegram

### Telegram Web App Features:
- **Init Data** - Авторизация через Telegram
- **Theme Params** - Адаптация под тему Telegram
- **Viewport** - Адаптация под размеры экрана
- **Back Button** - Навигация внутри приложения

### Telegram Bot Features:
- **Уведомления** - Оповещения о новых запросах
- **Чат** - Общение между пользователями

## 🧪 Тестирование

Проект планирует включение комплексного тестирования:

- **Unit тесты** для бизнес-логики
- **Интеграционные тесты** для API
- **E2E тесты** для критических пользовательских сценариев
- **Тестирование Telegram Web App SDK**

## 🔒 Безопасность

- Валидация всех входящих данных
- Проверка подписи Telegram Web App Init Data
- Rate limiting для API endpoints
- HTTPS в production

## 📈 Производительность

- Кэширование частых запросов в Redis (планируется)
- Оптимизированные SQL запросы
- CDN для статических ресурсов Telegram Mini App
- Мониторинг производительности (планируется)

**Сделано с ❤️ для наших четвероногих друзей** 🐾
