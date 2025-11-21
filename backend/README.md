# 🚀 Backend API - Одной Крови

<div align="center">

![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?style=for-the-badge&logo=go) ![Fiber](https://img.shields.io/badge/Fiber-2.52+-000000?style=for-the-badge&logo=go&logoColor=white) ![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16+-4169E1?style=for-the-badge&logo=postgresql&logoColor=white) ![Swagger](https://img.shields.io/badge/Swagger-API-85EA2D?style=for-the-badge&logo=swagger&logoColor=black)

</div>

**Backend API для платформы "Одной Крови"** — Go сервер для Telegram Mini App.

## 📁 Структура проекта

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

Примечание: в корне монорепозитория есть папки `shared/` и `microservices/`.  
- `shared/` — автосгенерированные TypeScript типы и клиенты, сгенерированные из Swagger/OpenAPI (используются frontend и bot).  
- `microservices/` — дополнительная логика и отдельные сервисы для backend (в перспективе — отдельные деплои/контейнеры).

## 🚀 Быстрый старт

### Предварительные требования

- **Go 1.25+**
- **PostgreSQL 16+**

### Установка и запуск

#### 1. Установка Go

##### macOS

**Вариант 1: Через Homebrew**
```bash
brew install go
```

**Вариант 2: Обычная установка**
1. Скачайте установщик `.pkg` с [официального сайта Go](https://go.dev/dl/)
2. Откройте скачанный файл и следуйте инструкциям установщика
3. Перезапустите терминал

Проверьте установку:
```bash
go version
```

##### Windows
1. Скачайте установщик с [официального сайта Go](https://go.dev/dl/)
2. Запустите `.msi` файл и следуйте инструкциям установщика
3. Перезапустите терминал
4. Проверьте установку:
```bash
go version
```

##### Linux (Ubuntu/Debian)
```bash
# Скачайте последнюю версию
wget https://go.dev/dl/go1.25.linux-amd64.tar.gz

# Распакуйте в /usr/local
sudo tar -C /usr/local -xzf go1.25.linux-amd64.tar.gz

# Добавьте в PATH (добавьте в ~/.profile или ~/.bashrc)
export PATH=$PATH:/usr/local/go/app

# Примените изменения
source ~/.profile

# Проверьте установку
go version
```

#### 2. Клонирование репозитория

```bash
git clone git@github.com:artesipov-alt/odnoi-krovi-app.git
cd odnoi-krovi-app
# Рекомендуется использовать Taskfile для автоматизации: task -l  (см. корневой Taskfile.yaml)
```

#### 3. Переход в директорию backend

```bash
cd odnoi-krovi-app/backend
```

#### 4. Настройка переменных окружения

> ⚠️ **Важно:** Перед запуском backend убедитесь, что вы создали файл `.env` в корневой директории проекта согласно [основной инструкции](../README.md#2-настройка-переменных-окружения).

Файл `.env` должен находиться в `odnoi-krovi-app/.env` (не в `odnoi-krovi-app/backend/.env`), так как все сервисы используют единый конфигурационный файл.

#### 5. Установка зависимостей

```bash
go mod download
```

#### 6. Запуск в режиме разработки

Вариант A — локально (Go):
```bash
go run cmd/server/main.go
```
Вариант B — через Taskfile (рекомендуется для монорепозитория; оркестрация зависит от Taskfile.yaml в корне):
```bash
# из корня репозитория
task dev-backend     # запустит backend в dev-режиме
# или
task dev             # запустит все сервисы (backend, frontend, bot) согласно Taskfile.yaml
```
Вариант C — через Docker Compose (локальное продакшн-подобное окружение):
```bash
# в корне репозитория
docker-compose up -d backend
# или поднять весь стек
docker-compose up -d
# Просмотр логов
docker-compose logs -f backend
```

Сервер будет доступен по адресу: **http://localhost:3000**

#### 7. Сборка для production

```bash
go build -o bin/server cmd/server/main.go
```

Или собрать и запустить контейнер через docker-compose:
```bash
docker-compose build backend
docker-compose up -d backend
```

Запуск собранного бинарника:
```bash
./app/server
```

## 📚 Документация API

После запуска сервера, документация Swagger будет доступна по адресу:

**http://localhost:3000/api/swagger**

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
go fmt ./...                 # Форматирование кода
go vet ./...                 # Статический анализ
```

**Сделано с ❤️ для наших четвероногих друзей** 🐾
