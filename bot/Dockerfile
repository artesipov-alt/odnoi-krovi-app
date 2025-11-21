FROM oven/bun:1-alpine AS builder

WORKDIR /app

# Copy package files first for better caching
COPY bot/package.json bot/bun.lock ./

# Clean cache and install dependencies
RUN bun install --frozen-lockfile --no-cache

# Copy bot source code preserving directory structure
COPY bot/src ./src
COPY bot/bot.ts ./
COPY bot/build.ts ./
COPY bot/tsconfig.json ./

# Copy shared directory from project root
COPY shared ./shared

RUN bun run build

FROM oven/bun:1-alpine

WORKDIR /app

COPY --from=builder /app/dist ./dist
COPY --from=builder /app/package.json ./

# Install production dependencies only
RUN bun install --production --no-cache

CMD ["bun", "run", "prod"]
