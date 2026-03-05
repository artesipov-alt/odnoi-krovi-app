# Use Bun alpine image for building
FROM oven/bun:1-alpine AS builder

# Use /app as project root in the container
WORKDIR /app

# Copy the bot folder as-is into /app/bot so relative imports like "../../shared/..." resolve
# (important: keep original repository layout)
COPY tg-bot ./tg-bot

# Copy shared folder at repository root (so /app/shared exists)
COPY shared ./shared

# Switch into the bot working directory to install and build
WORKDIR /app/tg-bot

# Install dependencies (use frozen lockfile to ensure reproducible builds)
RUN bun install --frozen-lockfile --no-cache

# Build the project (assumes `build` script in bot/package.json)
RUN bun run build

# Runtime image
FROM oven/bun:1-alpine

# Keep same workdir layout as during build so runtime imports still resolve
WORKDIR /app/tg-bot

# Copy built output and package files from builder
COPY --from=builder /app/tg-bot/dist ./dist
COPY --from=builder /app/tg-bot/package.json ./package.json
COPY --from=builder /app/tg-bot/bun.lock ./bun.lock

# Install production dependencies only
RUN bun install --production --no-cache

# Start the bot (assumes `prod` script in package.json)
CMD ["bun", "run", "prod"]
