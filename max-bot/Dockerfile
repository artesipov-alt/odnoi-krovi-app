# ---- builder ----
FROM node:24-alpine AS builder

WORKDIR /app

# Копируем репозиторий целиком, чтобы сохранить относительные импорты
# типа `../../shared/...` из max-bot/src.
COPY max-bot ./max-bot
COPY shared ./shared

WORKDIR /app/max-bot

# CI-воспроизводимая установка
RUN npm ci

# Сборка через tsc → dist/
RUN npm run build

# ---- runtime ----
FROM node:24-alpine

# На Bun 1.3.13 бот получал UNABLE_TO_GET_ISSUER_CERT_LOCALLY при
# запросах к platform-api2.max.ru несмотря на то, что сертификаты Минцифры
# были в /usr/local/share/ca-certificates/, update-ca-certificates их
# корректно подхватывал, а NODE_EXTRA_CERTS указывал на собранный bundle.
# Точная причина не диагностирована; в Bun 1.3.x наблюдались регрессии
# с TLS-валидацией (https://github.com/oven-sh/bun/issues/31949).
# Node.js 24 читает и системный CA-store, и NODE_EXTRA_CA_CERTS штатно,
# поэтому собранный bundle с CA Минцифры подхватывается автоматически.
COPY max-bot/certs/*.cer /usr/local/share/ca-certificates/
RUN apk add --no-cache ca-certificates \
    && update-ca-certificates \
    && ls -la /etc/ssl/certs/ca-certificates.crt

ENV NODE_EXTRA_CA_CERTS=/etc/ssl/certs/ca-certificates.crt

WORKDIR /app/max-bot

# Копируем только артефакты сборки и манифесты
COPY --from=builder /app/max-bot/dist ./dist
COPY --from=builder /app/max-bot/package.json ./package.json
COPY --from=builder /app/max-bot/package-lock.json ./package-lock.json

# Только prod-зависимости (без devDependencies)
RUN npm ci --omit=dev

# Запуск бота
CMD ["node", "--use-system-ca", "dist/bot.js"]
