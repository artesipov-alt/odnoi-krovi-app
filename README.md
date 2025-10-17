# 🩸 Одной Крови - Платформа донорства крови для животных

**Одной Крови** — это Telegram Mini App, который помогает находить донорскую кровь для животных и позволяет владельцам питомцев становиться донорами вместе со своими любимцами.

## 🎯 Цель проекта

Создать удобную и безопасную платформу в формате Telegram Mini App для:
- 🔍 Поиска доноров крови для животных в экстренных ситуациях
- ❤️ Регистрации питомцев как потенциальных доноров
- 📍 Геолокации ближайших доноров через Telegram Web App
- 🔔 Системы уведомлений через Telegram Bot API
- 💬 Интеграции с Telegram для удобного общения между пользователями

## 🏗️ Архитектура проекта

Проект построен как монорепозиторий с разделением на backend и Telegram Mini App:

```
odnoi-krovi-app/
├── backend/         # Go API сервер (Fiber + GORM)
├── frontend/        # Telegram Mini App (React + TypeScript)
├── bot/             # Telegram Bot (Node.js + Grammy)
├── docs/            # Документация
├── ROADMAP.md/      # Дорожная карта
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
git clone https://github.com/your-username/odnoi-krovi-app.git
cd odnoi-krovi-app
```

2. **Настройка Backend**
```bash
cd backend
cp .env.example .env
# Настройте переменные окружения в .env
go mod download
go run main.go
```

3. **Настройка Telegram Mini App**
```bash
cd tma
cp .env.example .env
npm install
npm run dev
```

4. **Настройка Telegram Bot**
```bash
cd bot
cp .env.example .env
npm install
npm start
```

5. **Запуск через Docker (рекомендуется)**
```bash
docker-compose up -d
```

## 📁 Структура проекта

### Backend (Go + Fiber + GORM)

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # Точка входа приложения
├── internal/
│   ├── config/                  # Конфигурация приложения
│   ├── handlers/                # HTTP обработчики
│   ├── models/                  # Модели данных (GORM)
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

### Telegram Mini App
- **React 18+** - UI библиотека
- **TypeScript** - Статическая типизация
- **Vite** - Сборщик и dev сервер
- **Telegram Web App SDK** - Интеграция с Telegram
- **React Query** - Управление состоянием сервера
- **React Hook Form** - Управление формами
- **Axios** - HTTP клиент

### Telegram Bot
- **Node.js** - Серверная платформа
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
go run cmd/api/main.go

# Запуск тестов
go test ./...
```

### Telegram Mini App разработка

```bash
cd tma
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
