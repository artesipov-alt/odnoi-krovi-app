# 🩸 Одной Крови - Платформа донорства крови для животных

<div align="center">

![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go) ![React](https://img.shields.io/badge/React-18+-61DAFB?style=for-the-badge&logo=react&logoColor=black) ![TypeScript](https://img.shields.io/badge/TypeScript-5.0+-3178C6?style=for-the-badge&logo=typescript&logoColor=white)
![Node.js](https://img.shields.io/badge/Node.js-18+-339933?style=for-the-badge&logo=nodedotjs&logoColor=white) ![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16+-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)

</div>

**Одной Крови** — это Telegram Mini App, который помогает находить донорскую кровь для животных и позволяет владельцам питомцев становиться донорами вместе со своими любимцами.

## 🗺️ Дорожная карта

Ознакомьтесь с [дорожной картой проекта](docs/roadmap.md) для понимания планов развития и текущих этапов разработки.

## 🎯 Цель проекта

Создать быстро удобную и безопасную платформу в формате Telegram Mini App для:
- 🔍 Поиска доноров крови для животных в экстренных ситуациях
- ❤️ Регистрации питомцев как потенциальных доноров
- 📍 Геолокации ближайших доноров через Telegram Web App
- 🔔 Системы уведомлений через Telegram Bot API
- 💬 Интеграции с Telegram для удобного общения между пользователями

## 🏗️ Архитектура проекта

Проект построен как монорепозиторий с разделением на backend и Telegram Mini App:

```
odnoi-krovi-app/
├── .github/         # Настрокий авторазвертывания (CD/CI)
├── backend/         # Go API сервер (Fiber + GORM + Swagger)
├── frontend/        # Telegram Mini App (React + TypeScript)
├── bot/             # Telegram Bot (Bun + Grammy)
├── docs/            # Документация
└── README.md
```

## 🚀 Быстрый старт

### Предварительные требования

- **Go 1.25+** для backend
- **Node.js 18+** и **npm** для Telegram Mini App
- **PostgreSQL 16+** для базы данных
- **Telegram Bot Token** от @BotFather
- **Docker** (опционально, для разработки)

### Установка и запуск

1. **Клонирование репозитория**
```bash
git clone git@github.com:artesipov-alt/odnoi-krovi-app.git
cd odnoi-krovi-app
```

2. **Настройка Backend**
```bash
cd backend
cp .env.example .env
# Настройте переменные окружения в .env
go mod download
go run cmd/server/main.go

# После запуска сервера документация API доступна по адресу:
# http://localhost:8080/swagger/index.html
```

3. **Настройка Telegram Mini App**
```bash
cd frontend
cp .env.example .env
npm install
npm run dev
```

4. **Запуск через Docker (рекомендуется)**
```bash
docker-compose up -d
```

## 📁 Структура проекта

### Backend (Go + Fiber + GORM + Swagger)

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # Точка входа приложения
├── internal/
│   ├── config/                  # Конфигурация приложения
│   ├── handlers/                # HTTP обработчики (с аннотациями Swagger)
│   ├── models/                  # Модели данных (DTO с Swagger тегами)
│   ├── repositories/            # Слой доступа к данным
│   ├── services/                # Бизнес-логика
│   ├── middleware/              # Промежуточное ПО
│   └── utils/                   # Вспомогательные функции
├── pkg/                         # Переиспользуемые пакеты
├── migrations/                  # Миграции базы данных
├── tests/                       # Тесты
└── go.mod
```

### Telegram Mini App (React + TypeScript + Telegram Web App SDK)

```
frontend/
├── public/                      # Статические файлы
├── src/
│   ├── components/              # Переиспользуемые компоненты
│   ├── pages/                   # Страницы приложения
│   ├── hooks/                   # Кастомные React хуки
│   ├── services/                # API клиенты + Telegram Web App SDK
│   ├── store/                   # Состояние приложения
│   ├── types/                   # TypeScript типы
│   ├── utils/                   # Вспомогательные функции
│   └── styles/                  # Стили
├── package.json
└── tsconfig.json
```

### Telegram Bot

```
bot/
├── src/
│   ├── handlers/               # Обработчики сообщений
│   ├── middleware/             # Промежуточное ПО
│   └── utils/                  # Вспомогательные функции
├── bot.ts
├── package.json
└── tsconfig.json
```

## 🛠️ Технологический стек

### Backend
- **Go 1.25+** - Основной язык программирования
- **Fiber** - Быстрый веб-фреймворк
- **GORM** - ORM для работы с базой данных
- **PostgreSQL** - Основная база данных
- **Redis** - Кэширование и сессии
- **Validator** - Валидация данных
- **Swagger** - Документация API (swaggo/swag)

### Telegram Mini App
- **React 18+** - UI библиотека
- **TypeScript** - Статическая типизация
- **Vite** - Сборщик и dev сервер
- **Telegram Web App SDK** - Интеграция с Telegram
- **React Query** - Управление состоянием сервера
- **React Hook Form** - Управление формами
- **Axios** - HTTP клиент

### Telegram Bot
- **Bun** - Серверная платформа
- **Grammy** - Фреймворк для Telegram Bot API
- **TypeScript** - Статическая типизация

## 📊 Модели данных

### Основные сущности:
- **Пользователи** - Владельцы животных (с привязкой к Telegram ID)
- **Животные** - Питомцы (доноры/реципиенты)
- **Запросы на кровь** - Заявки на поиск доноров
- **Донорские записи** - История донаций
- **Ветеринарные клиники** - Партнерские организации

## 🔧 Разработка

### Backend разработка

```bash
cd backend
# Запуск в режиме разработки
go run cmd/server/main.go

# Генерация Swagger документации
swag init -g cmd/api/main.go

# Запуск тестов
go test ./...
```

### Telegram Mini App разработка

```bash
cd frontend
# Запуск dev сервера
npm run dev

# Сборка для production
npm run build
```

### Telegram Bot разработка

```bash
cd bot
# Запуск бота
npm start

# Запуск в режиме разработки
npm run dev
```

## 🔗 Интеграция с Telegram

### Telegram Web App Features:
- **Init Data** - Авторизация через Telegram
- **Theme Params** - Адаптация под тему Telegram
- **Viewport** - Адаптация под размеры экрана
- **Back Button** - Навигация внутри приложения
- **Main Button** - Основное действие на экране
- **Haptic Feedback** - Тактильная обратная связь

### Telegram Bot Features:
- **Команды** - /start, /help, /search, /donate
- **Inline режим** - Быстрый поиск доноров
- **Уведомления** - Оповещения о новых запросах
- **Чат** - Общение между пользователями

## 🧪 Тестирование

Проект включает комплексное тестирование:

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

- Кэширование частых запросов в Redis
- Оптимизированные SQL запросы
- CDN для статических ресурсов Telegram Mini App
- Мониторинг производительности

**Сделано с ❤️ для наших четвероногих друзей** 🐾
