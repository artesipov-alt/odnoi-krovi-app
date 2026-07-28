# Odnoy Krovi — Монорепозиторий

Проект по поиску доноров крови среди животных. Монорепозиторий, содержащий несколько подпроектов.

## Структура

| Подпроект | Описание | AGENTS.md |
|---|---|---|
| `backend/` | Go-монолит (DDD + CQRS, Huma v2, Ent) | [backend/AGENTS.md](backend/AGENTS.md) |
| `frontend/` | React SPA (Vite, MUI, TanStack Query) | [frontend/AGENTS.md](frontend/AGENTS.md) |
| `max-bot/` | max-bot | — |
| `tg-bot/` | Telegram бот | — |
| `shared/` | Общие пакеты (OpenAPI-сгенерированный TS-клиент) | — |

## Глобальные конвенции

- CI/CD через Taskfile.yaml и docker-compose
- OpenAPI спецификация генерируется из бэкенда (backend/docs/openapi.json)
- Сгенерированный TS-клиент лежит в shared/ts/
- Миграции БД: schema (DDL) — Ent auto-migrate при старте; data (справочники) — SQL-файлы в backend/migrations/, применяются через `task db:migrate` (см. backend/migrations/README.md)

## Вложенные AGENTS.md

Каждый подпроект имеет свой `AGENTS.md` в корневой директории. Агенты автоматически читают ближайший файл в дереве директорий, поэтому при работе с файлами внутри `backend/` будет использован `backend/AGENTS.md`, внутри `frontend/` — `frontend/AGENTS.md` и т.д.