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

# Node.js корректно читает системный CA-store и переменную
# NODE_EXTRA_CA_CERTS, поэтому сертификаты Минцифры для platform-api2.max.ru
# подхватываются автоматически.
#
# ca-certificates нужен как базовый пакет, чтобы update-ca-certificates
# мог собрать итоговый bundle в /etc/ssl/certs/ca-certificates.crt.
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
CMD ["node", "dist/bot.js"]
