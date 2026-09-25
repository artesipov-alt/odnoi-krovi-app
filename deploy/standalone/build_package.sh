#!/usr/bin/env bash
# Сборка standalone-пакета экземпляра ПО «Одной Крови» для передачи
# на экспертизу (реестр российского ПО, п. 11 «а» и «д»).
#
# Результат: deploy/standalone/dist/odnoi-krovi-standalone-<version>.tar.gz
#   - docker-образы (docker save, работают без доступа в интернет);
#   - docker-compose.standalone.yml + Caddyfile + демо-данные;
#   - INSTALL.md — инструкция по установке и эксплуатации.
#
# Запуск из корня монорепозитория: deploy/standalone/build_package.sh

set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
cd "$REPO_ROOT"

DIST_DIR="deploy/standalone/dist"
STAGE_DIR="$DIST_DIR/odnoi-krovi-standalone"
BACKEND_IMAGE="odnoi-krovi-app/backend:standalone"
FRONTEND_IMAGE="odnoi-krovi-app/frontend:standalone"

echo "==> 1/6 Сборка frontend (bun run build)"
(cd frontend && bun run build)

echo "==> 2/6 Сборка docker-образов"
docker build -f backend/Dockerfile -t "$BACKEND_IMAGE" .
docker build -f frontend/Dockerfile -t "$FRONTEND_IMAGE" .

echo "==> 3/6 Подготовка содержимого пакета"
rm -rf "$STAGE_DIR"
mkdir -p "$STAGE_DIR/images"

cp deploy/standalone/Caddyfile "$STAGE_DIR/"
cp deploy/standalone/demo_data.sql "$STAGE_DIR/"
cp -r deploy/standalone/photos "$STAGE_DIR/photos"
cp deploy/standalone/INSTALL.md "$STAGE_DIR/"
cp -r backend/migrations/bootstrap "$STAGE_DIR/bootstrap"

# В пакете справочники лежат рядом с compose, а не в дереве монорепозитория
sed 's|\.\./\.\./backend/migrations/bootstrap|./bootstrap|' \
    deploy/standalone/docker-compose.standalone.yml > "$STAGE_DIR/docker-compose.standalone.yml"

echo "==> 4/6 Экспорт образов (docker save)"
docker save "$BACKEND_IMAGE" -o "$STAGE_DIR/images/backend-standalone.tar"
docker save "$FRONTEND_IMAGE" -o "$STAGE_DIR/images/frontend-standalone.tar"
# postgres и redis — чтобы пакет работал без доступа в интернет
docker save postgres:16-alpine redis:7-alpine -o "$STAGE_DIR/images/infra-standalone.tar"

echo "==> 5/6 Проверка содержимого"
ls -la "$STAGE_DIR" "$STAGE_DIR/images"

echo "==> 6/6 Архивация"
cd "$DIST_DIR"
tar -czf odnoi-krovi-standalone.tar.gz odnoi-krovi-standalone

echo "Готово: $DIST_DIR/odnoi-krovi-standalone.tar.gz ($(du -h odnoi-krovi-standalone.tar.gz | cut -f1))"
