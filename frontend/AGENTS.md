# Frontend — React SPA

## Общий обзор

Single Page Application для поиска доноров крови среди животных. Фреймворк — **React 18** + **TypeScript**. Сборка — **Vite 5**. Стейт-менеджмент — **Redux** + **TanStack React Query**. UI Kit — **MUI v7** (Material UI) + **Emotion** (styled). Роутинг — **React Router v7**. Интеграция с Telegram Mini App.

## Организация пакетов

```
frontend/
├── src/
│   ├── api/                # API-клиенты (Axios), запросы к бэкенду
│   │   ├── index.ts        # Экспорт всех API-модулей
│   │   ├── instance.ts     # Axios instance (base URL, interceptors)
│   │   ├── queryClient.ts  # TanStack Query client + кеш
│   │   ├── types.ts        # Общие типы (PetType, Analyzes и т.д.)
│   │   ├── auth.ts
│   │   ├── user.ts
│   │   ├── pets.ts
│   │   ├── bloodRequest.ts
│   │   ├── donor.ts
│   │   ├── reference.ts
│   │   ├── photo.ts
│   │   └── apiServices/    # OpenAPI-сгенерированные сервисы
│   ├── components/         # Переиспользуемые UI-компоненты
│   │   ├── Alert/
│   │   ├── Bonuses/
│   │   ├── CircularProgress/
│   │   ├── CompletedDonation/
│   │   ├── Curtain/
│   │   ├── DatePicker/
│   │   ├── ErrorBoundary/
│   │   ├── ImgEditor/
│   │   ├── Layout/
│   │   ├── Loading/
│   │   ├── Multiselect/
│   │   ├── Profiles/
│   │   ├── PromoSlider/
│   │   ├── RadioButton/
│   │   ├── SmsInput/
│   │   ├── Switch/
│   │   ├── TextField/
│   │   └── Timer/
│   ├── pages/              # Страницы (каждая в своей папке)
│   │   ├── about/
│   │   ├── adding/
│   │   ├── bonuses/
│   │   ├── donationsHistory/
│   │   ├── owner/
│   │   ├── profile/
│   │   ├── recipientsList/
│   │   ├── registration/
│   │   └── search/
│   ├── hooks/              # Кастомные React-хуки
│   │   ├── useAuth.ts
│   │   ├── useBodyScrollLock.ts
│   │   ├── useDicts.ts
│   │   ├── useGetPoolRequest.ts
│   │   ├── useGetRecipientsList.ts
│   │   ├── useGetUserById.ts
│   │   ├── usePetsQuery.ts
│   │   └── usePlannedDonations.ts
│   ├── types/              # TypeScript-типы и декларации
│   │   └── index.ts
│   ├── utils/              # Утилиты
│   │   ├── utils.ts
│   │   └── regexps.ts
│   ├── styles/             # Стили (Less, CSS, fonts)
│   │   ├── colors.less
│   │   ├── normalize.css
│   │   └── *.ttf (шрифты)
│   ├── imgs/               # Статические изображения
│   ├── App.tsx             # Корневой компонент
│   └── main.tsx            # Точка входа
├── build/                  # Сборка
├── mock_dist/              # Мок-данные для разработки
├── typings/                # Декларации типов
├── tsconfig.json           # Конфиг TypeScript
├── vite.config.js          # Конфиг Vite (алиасы, плагины)
└── package.json
```

## Алиасы (Vite resolve)

| Алиас | Путь |
|---|---|
| `api` | `src/api` |
| `components` | `src/components` |
| `pages` | `src/pages` |
| `hooks` | `src/hooks` |
| `utils` | `src/utils` |
| `styles` | `src/styles` |
| `imgs` | `src/imgs` |
| `services` | `src/services` |
| `context` | `src/context` |

Импорт всегда через алиасы: `import { useAuth } from 'hooks/useAuth'`

## Ключевые технологии

| Технология | Применение |
|---|---|
| **React 18** | UI-фреймворк |
| **TypeScript** | Типизация |
| **Vite 5** | Сборка, dev-сервер (порт 5173) |
| **MUI v7** | UI Kit (Material Design) |
| **Emotion** | CSS-in-JS (styled components) |
| **TanStack React Query** | Серверный стейт, кеширование, инвалидация |
| **Redux** | Клиентский стейт (auth, routing) |
| **React Router v7** | Клиентский роутинг |
| **Axios** | HTTP-клиент к бэкенду |
| **Less** | Препроцессор для CSS modules |
| **@twa-dev/sdk** | Telegram Mini App SDK |

## API-слой

- Базовый URL — через Axios instance (`src/api/instance.ts`).
- Типы API — частично ручные (`src/api/types.ts`), частично сгенерированные OpenAPI (`shared/ts/`).
- Для запросов данных — **TanStack Query** хуки в `src/hooks/`.
- Для мутаций — TanStack Query `useMutation`.
- OpenAPI-сгенерированный клиент лежит в `shared/ts/` (axios-based).

## Стилизация

- **CSS modules** с Less (файлы `*.module.less`).
- Глобальные переменные — `src/styles/colors.less`.
- CSS modules дают сгенерированные имена: `{fileName}_{name}_{hash}`.
- MUI компоненты стилизуются через **Emotion styled** или `sx` prop.

## Telegram Mini App

- Интеграция через `@twa-dev/sdk`.
- Аутентификация через Telegram Init Data, которая передаётся на бэкенд.
- Адаптация интерфейса под Telegram WebView.

## Команды

| Команда | Описание |
|---|---|
| `npm run dev` | Dev-сервер (Vite, порт 5173) |
| `npm run build` | Production сборка |
| `npm run typecheck` | Проверка типов TS (tsc --noEmit) |
| `npm run lint:ts` | ESLint проверка |
| `npm run prettier:fix` | Prettier форматирование |
| `npm run pre-commit` | Полная проверка перед коммитом (typecheck → lint → test) |