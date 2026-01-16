# ---------- STAGE 1: build ----------
FROM golang:1.25-alpine AS builder

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
COPY backend/app/ ./app/
COPY backend/docs/ ./docs/
COPY backend/docsui/ ./docsui/
COPY backend/ent/ ./ent/
COPY backend/pkg/ ./pkg/

# Генерируем свежую доку и делаем миграцию
RUN go run cmd/api/main.go -d

# Сборка бинарника
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/api


# ---------- STAGE 2: runtime ----------
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарник и документацию
COPY --from=builder /build/backend/main .
COPY --from=builder /build/backend/docs ./docs

EXPOSE 3000
CMD ["./main"]
