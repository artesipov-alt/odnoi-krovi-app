# ---------- STAGE 1: build ----------
FROM golang:1.25-alpine AS builder

# Устанавливаем рабочую директорию на уровень монорепо
WORKDIR /build

# Установка swag
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Копируем go.mod и go.sum из backend
COPY backend/go.mod backend/go.sum ./backend/

# Копируем локальные зависимости (для replace директив в go.mod)
COPY microservices/blood-microservice/go.mod microservices/blood-microservice/go.sum ./microservices/blood-microservice/
COPY microservices/blood-microservice/ ./microservices/blood-microservice/

# Скачиваем зависимости из директории backend
WORKDIR /build/backend
RUN go mod download

# Копируем исходный код backend
COPY backend/cmd/ ./cmd/
COPY backend/internal/ ./internal/
COPY backend/pkg/ ./pkg/
COPY backend/app/ ./app/
COPY backend/docs/ ./docs/

# Регенерация Swagger (обновляем существующие docs)
RUN swag init -g cmd/server/main.go -o docs || echo "Swagger generation skipped"

# Сборка бинарника
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server


# ---------- STAGE 2: runtime ----------
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарник и документацию
COPY --from=builder /build/backend/main .
COPY --from=builder /build/backend/docs ./docs

EXPOSE 3000
CMD ["./main"]
