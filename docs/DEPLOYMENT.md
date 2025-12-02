# 🚀 Deployment Guide

Руководство по развертыванию проекта "Одной Крови" в различных окружениях.

## 📋 Содержание

- [Локальная разработка](#локальная-разработка)
- [Development окружение (Docker)](#development-окружение-docker)
- [Production развертывание](#production-развертывание)
- [CI/CD Pipeline](#cicd-pipeline)
- [Troubleshooting](#troubleshooting)

## 🏠 Локальная разработка

### Без Docker

Каждый сервис запускается отдельно на localhost.

#### 1. Подготовка

```bash
# Клонируем репозиторий
git clone git@github.com:artesipov-alt/odnoi-krovi-app.git
cd odnoi-krovi-app

# Копируем .env файл
cp .env.example .env
# Отредактируйте .env с локальными настройками
```

#### 2. База данных

```bash
# PostgreSQL для основного backend
psql -U postgres -c "CREATE DATABASE odnoi_krovi;"

# PostgreSQL для микросервиса
psql -U postgres -c "CREATE DATABASE bloodsearch;"
```

#### 3. Запуск сервисов

**Терминал 1: Blood Microservice**
```bash
cd microservices/blood-microservice
export DATABASE_URL="host=localhost user=postgres password=postgres dbname=bloodsearch port=5432 sslmode=disable"
export SERVER_PORT=8081
export APP_ENV=development
go run cmd/microservice/main.go
```

**Терминал 2: Backend**
```bash
cd backend
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=odnoi_krovi
export SERVER_PORT=3000
export BLOOD_MICROSERVICE_URL=http://localhost:8081
go run cmd/server/main.go
```

**Терминал 3: Bot**
```bash
cd bot
export BOT_TOKEN=your_bot_token
export API_BASE_URL=http://localhost:3000/api/v1
export MINIAPP_DOMAIN=http://localhost:5173
bun run dev
```

**Терминал 4: Frontend**
```bash
cd frontend
npm install
npm run dev
```

#### 4. Проверка

```bash
# Health check микросервиса
curl http://localhost:8081/health

# Health check backend
curl http://localhost:3000/health

# Frontend
open http://localhost:5173
```

## 🐳 Development окружение (Docker)

Используется для разработки с Docker, но с возможностью отладки.

### Запуск

```bash
# Использовать dev конфигурацию
docker-compose -f docker-compose.dev.yml up -d

# Или через Taskfile
task docker-dev
```

### Особенности dev окружения

- ✅ Порты микросервиса **exposed** (8081) для debugging
- ✅ Hot reload (через volumes)
- ✅ Подробное логирование
- ✅ Быстрые healthchecks
- ✅ Development базы данных

### Логи

```bash
# Все логи
docker-compose -f docker-compose.dev.yml logs -f

# Конкретный сервис
docker-compose -f docker-compose.dev.yml logs -f blood-microservice

# Логи за последние 10 минут
docker-compose -f docker-compose.dev.yml logs --since 10m
```

### Остановка

```bash
docker-compose -f docker-compose.dev.yml down

# С удалением volumes
docker-compose -f docker-compose.dev.yml down -v
```

## 🚀 Production развертывание

### Архитектура

```
Internet
   │
   ↓
[Nginx/Traefik]
   │
   ├─→ Frontend (80, 443)
   │
   ├─→ Backend (3000)
   │     │
   │     └─→ Blood Microservice (8081, internal only)
   │
   └─→ Bot (webhook)
```

### 1. Подготовка сервера

```bash
# Обновляем систему
sudo apt update && sudo apt upgrade -y

# Устанавливаем Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh

# Устанавливаем Docker Compose
sudo apt install docker-compose-plugin

# Клонируем репозиторий
git clone git@github.com:artesipov-alt/odnoi-krovi-app.git
cd odnoi-krovi-app
```

### 2. Настройка переменных окружения

```bash
# Создаем .env файл
cat > .env << EOF
# Database (Main)
DB_HOST=10.0.0.129
DB_PORT=5432
DB_USER=admin
DB_PASSWORD=secure_password_here
DB_NAME=odnoi_krovi
DB_SSLMODE=disable

# Database (Blood Search Microservice)
DB_NAME_BLOOD_SEARCH=bloodsearch

# Redis
REDIS_HOST=10.0.0.73
REDIS_PORT=6379

# Telegram Bot
PROD_BOT_API_KEY=your_production_bot_token
MINIAPP_DOMAIN=https://1krovi.app
API_BASE_URL=https://1krovi.app/api/v1

# Server
SERVER_PORT=3000
ENVIRONMENT=production
APP_ENV=production
EOF

# Защищаем .env
chmod 600 .env
```

### 3. Сборка и запуск

```bash
# Собираем все сервисы
docker-compose build

# Запускаем в фоновом режиме
docker-compose up -d

# Проверяем статус
docker-compose ps

# Проверяем логи
docker-compose logs -f
```

### 4. Проверка работоспособности

```bash
# Health check микросервиса (изнутри backend контейнера)
docker-compose exec backend wget -O- http://blood-microservice:8081/health

# Health check backend
curl http://localhost:3000/health

# Проверка всех контейнеров
docker-compose ps
```

### 5. Обновление

```bash
# Получаем последние изменения
git pull origin main

# Пересобираем измененные сервисы
docker-compose build

# Перезапускаем с zero-downtime (если настроен)
docker-compose up -d

# Или перезапускаем конкретный сервис
docker-compose up -d --no-deps --build backend
```

### Особенности production

- ❌ Порт микросервиса **НЕ exposed** наружу (только internal)
- ✅ Health checks с retry
- ✅ Автоматический перезапуск (`restart: unless-stopped`)
- ✅ Production базы данных
- ✅ HTTPS через Nginx
- ✅ Graceful shutdown

## 🔄 CI/CD Pipeline

### GitHub Actions Workflow

Pipeline автоматически:
1. Отслеживает изменения в конкретных директориях
2. Пересобирает только измененные сервисы
3. Применяет изменения с минимальным downtime

### Триггеры

```yaml
# При push в dev ветку
on:
  push:
    branches:
      - dev
```

### Отслеживаемые директории

- `backend/**` → пересборка backend
- `bot/**` → пересборка bot
- `frontend/**` → пересборка frontend
- `microservices/**` → пересборка blood-microservice

### Последовательность развертывания

1. **Microservices** (первым, так как backend зависит от него)
2. **Backend**
3. **Bot**
4. **Frontend**

### Ручной деплой

```bash
# На сервере
cd ~/odnoi-krovi-app
git pull origin dev

# Пересобрать конкретный сервис
docker-compose build blood-microservice
docker-compose up -d blood-microservice

# Или все сразу
docker-compose build
docker-compose up -d
```

## 🔍 Мониторинг

### Проверка статуса сервисов

```bash
# Статус всех контейнеров
docker-compose ps

# Использование ресурсов
docker stats

# Логи в реальном времени
docker-compose logs -f
```

### Health Checks

```bash
# Микросервис (изнутри)
docker-compose exec backend curl http://blood-microservice:8081/health

# Backend
curl http://localhost:3000/health

# Проверка health status
docker-compose ps | grep healthy
```

### Метрики (планируется)

- Prometheus для сбора метрик
- Grafana для визуализации
- AlertManager для алертов

## 🛠️ Troubleshooting

### Микросервис недоступен для backend

**Симптомы:**
```
Error: dial tcp: lookup blood-microservice: no such host
```

**Решение:**
```bash
# Проверить, что микросервис запущен
docker-compose ps blood-microservice

# Проверить логи
docker-compose logs blood-microservice

# Проверить сеть
docker network inspect odnoi-krovi-app_odnoi_net

# Проверить доступность изнутри backend
docker-compose exec backend ping blood-microservice
```

### Контейнер постоянно перезапускается

**Симптомы:**
```bash
docker-compose ps
# blood-microservice   Restarting
```

**Решение:**
```bash
# Смотрим логи
docker-compose logs --tail=100 blood-microservice

# Проверяем health check
docker inspect odnoi-krovi-app-blood-microservice-1 | grep -A 10 Health

# Запускаем в interactive режиме для отладки
docker-compose run --rm blood-microservice sh
```

### База данных недоступна

**Симптомы:**
```
Error: failed to connect to database
```

**Решение:**
```bash
# Проверить подключение к PostgreSQL
docker-compose exec backend nc -zv 10.0.0.129 5432

# Проверить credentials
docker-compose exec backend env | grep DB_

# Проверить из контейнера микросервиса
docker-compose exec blood-microservice nc -zv 10.0.0.129 5432
```

### Порты уже заняты

**Симптомы:**
```
Error: bind: address already in use
```

**Решение:**
```bash
# Найти процесс на порту
sudo lsof -i :8081
sudo lsof -i :3000

# Остановить старые контейнеры
docker-compose down

# Убить процесс (если не Docker)
kill -9 <PID>
```

### Нет места на диске

**Симптомы:**
```
Error: no space left on device
```

**Решение:**
```bash
# Очистить неиспользуемые образы
docker system prune -a

# Удалить старые volumes
docker volume prune

# Удалить stopped контейнеры
docker container prune

# Проверить размер
docker system df
```

## 📊 Полезные команды

```bash
# Быстрый рестарт всех сервисов
docker-compose restart

# Рестарт конкретного сервиса
docker-compose restart blood-microservice

# Пересобрать и перезапустить без кэша
docker-compose build --no-cache
docker-compose up -d --force-recreate

# Посмотреть переменные окружения
docker-compose exec backend env

# Зайти в контейнер
docker-compose exec backend sh

# Экспорт логов
docker-compose logs > logs.txt

# Проверить конфигурацию
docker-compose config

# Обновить один сервис без downtime остальных
docker-compose up -d --no-deps --build blood-microservice
```

## 🔐 Безопасность

### Чеклист для production

- [ ] `.env` файл в `.gitignore`
- [ ] Сильные пароли для баз данных
- [ ] HTTPS настроен
- [ ] Firewall правила настроены
- [ ] Микросервис не exposed наружу
- [ ] SSH ключи вместо паролей
- [ ] Регулярные бэкапы БД
- [ ] Мониторинг и алерты настроены
- [ ] Rate limiting настроен
- [ ] Логи ротируются

### Рекомендации

1. Используйте секреты Docker для sensitive данных
2. Регулярно обновляйте базовые образы
3. Сканируйте образы на уязвимости
4. Ограничивайте ресурсы контейнеров
5. Используйте read-only файловые системы где возможно

---

**Документация обновлена:** 2025-12-02

**Сделано с ❤️ для наших четвероногих друзей** 🐾