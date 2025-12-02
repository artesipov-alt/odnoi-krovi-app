# ---------- STAGE 1: build ----------
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Копируем модули отдельно, чтобы кешировать зависимости
COPY microservices/blood-microservice/go.mod microservices/blood-microservice/go.sum ./
RUN go mod download

# Копируем только исходный код для сборки
COPY microservices/blood-microservice/cmd/ ./cmd/
COPY microservices/blood-microservice/internal/ ./internal/
COPY microservices/blood-microservice/pkg/ ./pkg/
COPY microservices/blood-microservice/gen/ ./gen/

# Сборка бинарника
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/microservice


# ---------- STAGE 2: runtime ----------
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Копируем бинарник
COPY --from=builder /app/main .

EXPOSE 8081

CMD ["./main"]
