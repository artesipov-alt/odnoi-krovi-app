FROM oven/bun:1-alpine AS builder

WORKDIR /app

COPY package.json bun.lock ./

# Очищаем кэш и устанавливаем зависимости
RUN bun install --frozen-lockfile --no-cache

COPY . .
RUN bun run build

FROM oven/bun:1-alpine

WORKDIR /app

COPY --from=builder /app/dist ./dist
COPY --from=builder /app/package.json ./

CMD ["bun", "run", "prod"]
