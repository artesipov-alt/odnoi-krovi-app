# ---------- STAGE 1: build ----------
FROM golang:1.26-alpine AS builder

# Устанавливаем рабочую директорию на уровень монорепо
WORKDIR /build


# Копируем go.mod и go.sum из backend
COPY backend/go.mod backend/go.sum ./backend/

# Копируем локальные зависимости (для replace директив в go.mod)


# Скачиваем зависимости из директории backend
WORKDIR /build/backend
RUN go mod download

# Копируем исходный код backend
COPY backend/cmd/ ./cmd/
COPY backend/internal/ ./internal/
COPY backend/docs/ ./docs/
COPY backend/docsui/ ./docsui/
COPY backend/pkg/ ./pkg/

# Сборка бинарников: API-сервер и runner one-off миграций
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o migrate ./cmd/migrate


# ---------- STAGE 2: runtime ----------
FROM alpine:latest

RUN apk --no-cache add ca-certificates curl

WORKDIR /root/

# Копируем бинарники и документацию
COPY --from=builder /build/backend/main .
COPY --from=builder /build/backend/migrate .
COPY --from=builder /build/backend/docs ./docs
# SQL-миграции для runner-а (./migrate ищет migrations/one-off относительно WORKDIR)
COPY backend/migrations ./migrations

EXPOSE 3000
CMD ["./main"]
