# 🚀 Backend API - Одной Крови

<div align="center">

![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go) ![Fiber](https://img.shields.io/badge/Fiber-2.52+-000000?style=for-the-badge&logo=go&logoColor=white) ![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16+-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)

</div>

**Backend API для платформы "Одной Крови"** — Go сервер для Telegram Mini App.

## 🏗️ Структура проекта

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # Точка входа
├── internal/
│   ├── handlers/                # HTTP обработчики
│   ├── models/                  # Модели данных
│   ├── repositories/            # Доступ к данным
│   ├── services/                # Бизнес-логика
│   ├── middleware/              # Промежуточное ПО
│   └── utils/                   # Вспомогательные функции
├── pkg/
│   ├── config/                  # Конфигурация
│   └── logger/                  # Логирование
├── migrations/                  # Миграции БД
└── go.mod
```

## 🚀 Быстрый старт

### Предварительные требования

- **Go 1.25+**
- **PostgreSQL 16+**

### Установка и запуск

1. **Установка зависимостей**
```bash
go mod download
```

2. **Запуск в режиме разработки**
```bash
go run cmd/server/main.go
```

Сервер будет доступен по адресу: **http://localhost:3000**

3. **Сборка для production**
```bash
go build -o bin/server cmd/server/main.go
```

## 🛠️ Технологический стек

- **Go 1.25+** - Основной язык
- **Fiber** - Веб-фреймворк
- **PostgreSQL** - База данных
- **Swagger** - Документация API

## 🔧 Разработка

### Основные команды
```bash
go run cmd/server/main.go    # Запуск сервера
go test ./...                # Запуск тестов
go build                     # Сборка проекта
```

**Сделано с ❤️ для наших четвероногих друзей** 🐾