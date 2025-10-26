FROM golang:1.25-alpine AS builder

WORKDIR /app

# Установка swag CLI для генерации Swagger документации
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Генерация Swagger документации
RUN swag init -g cmd/server/main.go -o docs

# Сборка приложения
RUN go build -o main ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/.env .
COPY --from=builder /app/docs ./docs

EXPOSE 3000

CMD ["./main"]
