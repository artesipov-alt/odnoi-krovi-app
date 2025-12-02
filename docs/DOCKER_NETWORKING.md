# 🐳 Docker Networking и Межсервисное Взаимодействие

## 📋 Обзор

В проекте "Одной Крови" используется Docker Compose для оркестрации нескольких сервисов, которые взаимодействуют друг с другом через внутреннюю Docker сеть.

## 🌐 Архитектура сети

```
┌─────────────────────────────────────────────────────────────┐
│                      odnoi_net (bridge)                      │
│                                                              │
│  ┌──────────┐    ┌──────────────────┐    ┌──────────────┐  │
│  │ Frontend │    │     Backend      │    │     Bot      │  │
│  │  :80     │◄───┤      :3000       ├───►│              │  │
│  │  :443    │    │                  │    │              │  │
│  └──────────┘    └────────┬─────────┘    └──────────────┘  │
│                           │                                  │
│                           ↓ HTTP                             │
│                  ┌─────────────────────┐                     │
│                  │ Blood Microservice  │                     │
│                  │       :8081         │                     │
│                  └─────────────────────┘                     │
│                                                              │
└─────────────────────────────────────────────────────────────┘
                           │
                           ↓
                  ┌─────────────────────┐
                  │   PostgreSQL DB     │
                  │  10.0.0.129:5432   │
                  │                    │
                  │ - odnoi_krovi (main)│
                  │ - bloodsearch (ms) │
                  └─────────────────────┘
```

## 🔧 Конфигурация Docker Compose

### Сеть

```yaml
networks:
  odnoi_net:
    driver: bridge
```

Все сервисы подключены к одной bridge-сети `odnoi_net`, что позволяет им общаться друг с другом.

### Сервисы

| Сервис | Порт (host:container) | DNS имя в сети | Описание |
|--------|----------------------|----------------|----------|
| `frontend` | 80:80, 443:443 | frontend | Nginx с Telegram Mini App |
| `backend` | 3000:3000 | backend | REST API (Fiber) |
| `bot` | - | bot | Telegram Bot (Grammy) |
| `blood-microservice` | 8081:8081 | blood-microservice | Микросервис поиска крови |

## 🔗 Межсервисное взаимодействие

### 1. Backend → Blood Microservice

**Внутри Docker сети:**
```yaml
# docker-compose.yml
backend:
  environment:
    BLOOD_MICROSERVICE_URL: http://blood-microservice:8081
```

**В коде (backend):**
```go
// Использует имя сервиса из docker-compose
bloodSearchClient := services.NewBloodSearchClient("http://blood-microservice:8081")

// Docker автоматически резолвит blood-microservice в IP контейнера
pets, err := bloodSearchClient.GetPets(ctx, filter)
```

**Локальная разработка:**
```env
# .env для локальной разработки
BLOOD_MICROSERVICE_URL=http://localhost:8081
```

### 2. Bot → Backend

```yaml
# docker-compose.yml
bot:
  environment:
    API_BASE_URL: https://1krovi.app/api/v1  # Production через HTTPS
```

В production бот обращается к backend через публичный домен, так как webhook от Telegram приходит извне.

### 3. Frontend → Backend

Frontend работает в браузере пользователя и обращается к backend через публичный URL:

```typescript
// frontend/src/config.ts
const API_BASE_URL = 'https://1krovi.app/api/v1';
```

## 🏠 Локальная разработка vs Docker

### Локальная разработка (без Docker)

```env
# Backend
BLOOD_MICROSERVICE_URL=http://localhost:8081

# Bot  
API_BASE_URL=http://localhost:3000/api/v1

# Frontend
VITE_API_BASE_URL=http://localhost:3000/api/v1
```

Все сервисы работают на localhost с разными портами.

### Docker Compose

```yaml
# Backend
BLOOD_MICROSERVICE_URL: http://blood-microservice:8081
#                              ↑ DNS имя сервиса

# Backend внутри контейнера НЕ использует localhost!
```

## 📊 DNS Resolution в Docker

Docker автоматически создает DNS записи для каждого сервиса:

```bash
# Внутри контейнера backend можно резолвить другие сервисы:
docker exec -it odnoi-krovi-app-backend-1 /bin/sh
ping blood-microservice
# PING blood-microservice (172.18.0.4): 56 data bytes
```

Docker DNS резолвит имена сервисов в IP адреса контейнеров в той же сети.

## 🔍 Проверка связности

### 1. Проверить сеть

```bash
# Список сетей
docker network ls

# Детали сети odnoi_net
docker network inspect odnoi-krovi-app_odnoi_net

# Какие контейнеры в сети
docker network inspect odnoi-krovi-app_odnoi_net -f '{{range .Containers}}{{.Name}} - {{.IPv4Address}}{{"\n"}}{{end}}'
```

### 2. Проверить связь между контейнерами

```bash
# Зайти в контейнер backend
docker exec -it odnoi-krovi-app-backend-1 sh

# Проверить доступность микросервиса
wget -O- http://blood-microservice:8081/health
# или
curl http://blood-microservice:8081/health
```

### 3. Проверить логи

```bash
# Логи backend
docker compose logs -f backend

# Логи микросервиса
docker compose logs -f blood-microservice

# Все логи вместе
docker compose logs -f
```

## 🚨 Типичные проблемы

### 1. "Connection refused" ошибка

**Проблема:**
```
Error: dial tcp: lookup blood-microservice: no such host
```

**Решение:**
- Убедитесь, что все сервисы в одной сети
- Проверьте, что микросервис запущен: `docker compose ps`
- Используйте имя сервиса из docker-compose.yml, а не `localhost`

### 2. Timeout при обращении к микросервису

**Проблема:**
```
Error: context deadline exceeded
```

**Решение:**
- Проверьте, что микросервис запущен: `docker compose logs blood-microservice`
- Увеличьте timeout в HTTP клиенте
- Проверьте health check: `docker exec backend curl http://blood-microservice:8081/health`

### 3. Backend не может найти микросервис

**Проблема:**
```
BLOOD_MICROSERVICE_URL=http://localhost:8081  # ❌ Неправильно в Docker
```

**Решение:**
```yaml
# docker-compose.yml
BLOOD_MICROSERVICE_URL: http://blood-microservice:8081  # ✅ Правильно
```

## 🔐 Security Best Practices

### 1. Не expose порты без необходимости

```yaml
# ❌ Плохо - микросервис доступен извне
blood-microservice:
  ports:
    - "8081:8081"

# ✅ Хорошо для production - только внутри Docker сети
blood-microservice:
  expose:
    - 8081
  # Доступен только для других контейнеров в той же сети
```

### 2. Используйте отдельные сети для изоляции

```yaml
networks:
  frontend_net:  # Frontend <-> Backend
  backend_net:   # Backend <-> Microservices
  db_net:        # Microservices <-> Database

backend:
  networks:
    - frontend_net
    - backend_net

blood-microservice:
  networks:
    - backend_net
```

## 📈 Production Considerations

### 1. Health Checks

```yaml
blood-microservice:
  healthcheck:
    test: ["CMD", "wget", "-q", "--spider", "http://localhost:8081/health"]
    interval: 30s
    timeout: 10s
    retries: 3
    start_period: 40s
```

### 2. Dependencies

```yaml
backend:
  depends_on:
    blood-microservice:
      condition: service_healthy  # Ждать пока микросервис будет ready
```

### 3. Resource Limits

```yaml
blood-microservice:
  deploy:
    resources:
      limits:
        cpus: '0.5'
        memory: 512M
      reservations:
        cpus: '0.25'
        memory: 256M
```

## 🔄 Restart Policies

```yaml
# Автоматический перезапуск при сбое
restart: unless-stopped

# Или с задержкой
restart: on-failure
deploy:
  restart_policy:
    condition: on-failure
    delay: 5s
    max_attempts: 3
```

## 📝 Полезные команды

```bash
# Проверить статус всех сервисов
docker compose ps

# Перезапустить конкретный сервис
docker compose restart blood-microservice

# Пересобрать и перезапустить
docker compose up -d --build blood-microservice

# Посмотреть IP адреса всех контейнеров
docker compose exec backend cat /etc/hosts

# Проверить переменные окружения
docker compose exec backend env | grep BLOOD

# Тест связности между сервисами
docker compose exec backend wget -O- http://blood-microservice:8081/health
```

## 🎓 Дополнительные ресурсы

- [Docker Networking Overview](https://docs.docker.com/network/)
- [Docker Compose Networking](https://docs.docker.com/compose/networking/)
- [Service Discovery in Docker](https://docs.docker.com/network/#service-discovery)

---

**Важно:** В Docker сервисы общаются через DNS имена (имена сервисов), а не через `localhost`!