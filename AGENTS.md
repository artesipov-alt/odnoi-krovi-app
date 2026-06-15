# Odnoy Krovi — Монорепозиторий

Проект по поиску доноров крови среди животных. Монорепозиторий, содержащий несколько подпроектов.

## Структура

| Подпроект | Описание | AGENTS.md |
|---|---|---|
| `backend/` | Go-монолит (DDD + CQRS, Huma v2, Ent) | [.agents/backend/AGENTS.md](.agents/backend/AGENTS.md) |
| `frontend/` | React SPA (Vite, MUI, TanStack Query) | [.agents/frontend/AGENTS.md](.agents/frontend/AGENTS.md) |
| `max-bot/` | max-bot | — |
| `tg-bot/` | Telegram бот | — |
| `shared/` | Общие пакеты (OpenAPI-сгенерированный TS-клиент) | — |

## Глобальные конвенции

- CI/CD через Taskfile.yaml и docker-compose
- OpenAPI спецификация генерируется из бэкенда (backend/docs/openapi.json)
- Сгенерированный TS-клиент лежит в shared/ts/