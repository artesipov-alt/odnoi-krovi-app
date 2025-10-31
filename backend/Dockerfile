# ---------- STAGE 1: build ----------
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Установка swag (один раз)
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Копируем модули отдельно, чтобы кешировать зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальной код
COPY . .

# Генерация Swagger (без фатального падения, если swag не нужен)
RUN swag init -g cmd/server/main.go -o docs || echo "Swagger generation skipped"

# Сборка бинарника
RUN go build -v -o main ./cmd/server


# ---------- STAGE 2: runtime ----------
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарник и нужные файлы
COPY --from=builder /app/main .
COPY --from=builder /app/.env .
COPY --from=builder /app/docs ./docs

EXPOSE 3000
CMD ["./main"]
