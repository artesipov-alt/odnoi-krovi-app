# Use Bun alpine image for building
FROM oven/bun:1.3.13-alpine AS builder

# Use /app as project root in the container
WORKDIR /app

# Copy the bot folder as-is into /app/bot so relative imports like "../../shared/..." resolve
# (important: keep original repository layout)
COPY max-bot ./max-bot

# Copy shared folder at repository root (so /app/shared exists)
COPY shared ./shared

# Switch into the bot working directory to install and build
WORKDIR /app/max-bot

# Install dependencies (use frozen lockfile to ensure reproducible builds)
RUN bun install --no-cache

# Build the project (assumes `build` script in bot/package.json)
RUN bun run build

# Runtime image
FROM oven/bun:1.3.13-alpine

# Сертификаты Минцифры для TLS-валидации platform-api2.max.ru.
# Без них fetch к Max API падает с UNABLE_TO_GET_ISSUER_CERT_LOCALLY,
# и подписка вебхука (src/maxbot.ts) не проходит.
#
# ca-certificates нужен как базовый пакет, чтобы update-ca-certificates
# мог собрать итоговый bundle в /etc/ssl/certs/ca-certificates.crt.
COPY max-bot/certs/*.cer /usr/local/share/ca-certificates/
RUN apk add --no-cache ca-certificates \
    && update-ca-certificates \
    && ls -la /etc/ssl/certs/ca-certificates.crt

# Подстраховка: Bun может использовать свой CA-bundle, отличный от
# системного. Явно указываем путь, чтобы fetch точно доверял нашим CA.
ENV NODE_EXTRA_CA_CERTS=/etc/ssl/certs/ca-certificates.crt

# Keep same workdir layout as during build so runtime imports still resolve
WORKDIR /app/max-bot

# Copy built output and package files from builder
COPY --from=builder /app/max-bot/dist ./dist
COPY --from=builder /app/max-bot/package.json ./package.json
COPY --from=builder /app/max-bot/bun.lock ./bun.lock

# Install production dependencies only
RUN bun install --production --no-cache

# Start the bot (assumes `prod` script in package.json)
CMD ["bun", "run", "prod"]
