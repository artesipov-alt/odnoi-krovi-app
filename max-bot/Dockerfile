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

# Устанавливаем ca-certificates и копируем сертификаты Минцифры
# Копируем все файлы сертификатов (и .cer, и .crt)
# update-ca-certificates ищет файлы с расширением .crt
COPY max-bot/certs/ /tmp/certs/
RUN apk add --no-cache ca-certificates \
    && mkdir -p /usr/local/share/ca-certificates \
    && cp /tmp/certs/*.crt /usr/local/share/ca-certificates/ \
    && for f in /tmp/certs/*.cer; do \
         [ -f "$f" ] && cp "$f" "/usr/local/share/ca-certificates/$(basename $f .cer).crt"; \
       done \
    && update-ca-certificates \
    && echo "=== Certificates in bundle: $(grep -c 'BEGIN CERTIFICATE' /etc/ssl/certs/ca-certificates.crt) ===" \
    && echo "=== Added Russian Trusted certs ===" \
    && ls -la /usr/local/share/ca-certificates/

ENV NODE_EXTRA_CA_CERTS=/etc/ssl/certs/ca-certificates.crt

WORKDIR /app/max-bot

# Копируем только артефакты сборки и манифесты
COPY --from=builder /app/max-bot/dist ./dist
COPY --from=builder /app/max-bot/package.json ./package.json
COPY --from=builder /app/max-bot/package-lock.json ./package-lock.json

# Только prod-зависимости (без devDependencies)
RUN npm ci --omit=dev

# Запуск бота (убран --use-system-ca, это флаг Bun)
CMD ["node", "dist/bot.js"]
