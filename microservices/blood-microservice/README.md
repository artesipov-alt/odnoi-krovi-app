# 🩸 Blood Search Microservice

Микросервис для управления поиском крови для животных. Простой сервис для добавления питомцев в базу и поиска подходящих доноров по нужным критериям.

## 🎯 Назначение

Микросервис помогает:
- Сохранять информацию о питомцах, которые нуждаются в крови
- Искать доноров по типу животного, группе крови и региону
- Отслеживать статус заявок и объемы крови

## 📋 API

### Проверка состояния сервиса

```bash
curl http://localhost:8081/health
```

Ответ содержит статус сервиса.

### Добавить питомца

```bash
curl -X POST http://localhost:8081/bloodsearch.v1.BloodSearchPool/AddPet \
  -H "Content-Type: application/json" \
  -d '{
    "pet_id": "pet-123",
    "pet_type": "dog",
    "blood_group": "DEA 1.1+",
    "blood_volume_needed": 500,
    "regions": [1, 2, 3],
    "status": "active"
  }'
```

Питомец добавляется в базу.

### Поиск питомцев

```bash
curl -X POST http://localhost:8081/bloodsearch.v1.BloodSearchPool/GetPets \
  -H "Content-Type: application/json" \
  -d '{
    "pet_type": "dog",
    "blood_group": "DEA 1.1+",
    "regions": [1, 2]
  }'
```

Возвращается список подходящих питомцев.

## ⚙️ Переменные окружения

| Переменная  | Описание           | Значение по умолчанию                              |
|-------------|--------------------|---------------------------------------------------|
| DATABASE_URL| Подключение к базе | `host=localhost user=postgres password=postgres dbname=bloodsearch port=5432 sslmode=disable` |
| SERVER_PORT | Порт сервера       | `8081`                                            |
| APP_ENV     | Окружение          | `dev`                                             |

Пример файла `.env`:

```env
DATABASE_URL=host=localhost user=postgres password=postgres dbname=bloodsearch port=5432 sslmode=disable
SERVER_PORT=8081
APP_ENV=development
```

## 🚀 Запуск

### Локально

1. Установите зависимости:
```bash
go mod download
```

2. Создайте базу данных:
```bash
psql -U postgres -c "CREATE DATABASE bloodsearch;"
```

3. Настройте переменные окружения:
```bash
cp .env.example .env
# Отредактируйте .env при необходимости
```

4. Запустите сервис:
```bash
go run cmd/microservice/main.go
```

Сервис будет доступен по адресу `http://localhost:8081`.

### Docker

1. Соберите Docker-образ из корня проекта:
```bash
docker build -f microservices/blood-microservice/Dockerfile -t blood-microservice .
```

2. Запустите контейнер:
```bash
docker run -p 8081:8081 \
  -e DATABASE_URL="host=host.docker.internal user=postgres password=postgres dbname=bloodsearch port=5432 sslmode=disable" \
  -e APP_ENV=production \
  blood-microservice
```

### Docker Compose

```bash
docker-compose up blood-microservice
```

## 🧪 Тестирование

### Проверка состояния

```bash
curl http://localhost:8081/health
```

### Добавление питомца

```bash
curl -X POST http://localhost:8081/bloodsearch.v1.BloodSearchPool/AddPet \
  -H "Content-Type: application/json" \
  -d '{
    "pet_id": "test-pet-1",
    "pet_type": "dog",
    "blood_group": "DEA 1.1+",
    "blood_volume_needed": 500,
    "regions": [77],
    "status": "active"
  }'
```

### Поиск питомцев

```bash
curl -X POST http://localhost:8081/bloodsearch.v1.BloodSearchPool/GetPets \
  -H "Content-Type: application/json" \
  -d '{
    "pet_type": "dog",
    "blood_group": "DEA 1.1+",
    "regions": [77]
  }'
```

## 📝 Changelog

### v1.0.0 (2025-12-02)
- Начальная версия с добавлением и поиском питомцев
- Поддержка PostgreSQL
- Health check
- Запуск в Docker
- Структурированное логирование
- Корректное завершение работы сервиса

### Планы (v1.1.0)
- Добавить кеширование через Redis
- Написать unit-тесты
- Добавить поддержку очередей для уведомлений
- Управление запасами крови

## 📄 Лицензия

Часть проекта "Одной Крови" — платформы донорства крови для животных.

---

**Сделано с ❤️ для наших четвероногих друзей** 🐾
