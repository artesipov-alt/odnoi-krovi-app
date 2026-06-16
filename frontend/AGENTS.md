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

## Загрузка фото (известные проблемы и решения)

### Механизм загрузки (3 шага)
1. `POST /v1/uploads/presign/{id}` — получение presigned PUT URL для S3 (VK Cloud)
2. `PUT {presigned_url}` — загрузка файла напрямую в S3
3. `POST /v1/uploads/confirm` — подтверждение

### Проблема Telegram WebView на Android

Telegram Android WebView имеет два бага, отсутствующих в Max WebView и iOS:

**1. Перехват файлового инпута:** `accept='image/*'` заставляет Telegram показать свою галерею вместо системного диалога. Камера и файлы недоступны.
- **Симптом:** системный диалог выбора (камера / галерея / файлы) не отображается, сразу открывается галерея Telegram.
- **Решение:** `accept='image/*,application/pdf'` — `application/pdf` вынуждает Telegram отдать системный пикер. При этом PDF отфильтровываются в коде (принимаются только изображения).
- Файлы: `ImgEditor/index.tsx`, `profile/index.tsx`

**2. Протухающий `content://` URI:** Android выдаёт временный URI, разрешение отзывается после закрытия пикера.
- **Симптом (вебхук):** `"The requested file could not be read, typically due to permission problems that have occurred after a reference to a file was acquired."`
- **Решение:** `arrayBuffer()` сразу в `onChange`, создание in-memory `File` до закрытия пикера.
- Файлы: `ImgEditor/index.tsx:68-69`, `profile/index.tsx:416-418`

**3. CORS preflight к S3:** `PUT` с `Content-Type` требует preflight, который Telegram WebView может блокировать.
- **Симптом (вебхук):** `"Failed to fetch"` — ошибка из `fetch(url, { method: 'PUT', body: file })` при загрузке на S3. При этом `arrayBuffer()` в `onChange` отрабатывает успешно, `getPhotoLink` возвращает presigned URL успешно, падает именно PUT к S3.
- **Характер:** интермиттентный (из 5 попыток 4 fail, 1 success — паттерн CORS preflight-кеширования).
- **Решение:** PUT без `Content-Type` заголовка, тело — `ArrayBuffer`. Presigned URL на бэкенде не проверяет ContentType.
- Файл: `api/apiServices/addPhoto.ts` (`putFile`)
- **Резерв:** `putFileViaXHR` — XHR без кастомных заголовков, сохранён в коде, но не вызывается. Если `fetch` PUT снова начнёт падать — раскомментировать fallback в `addPhoto`.
- **Инфраструктура:** S3 CORS должен быть настроен на бакете (AllowedOrigins, AllowedMethods: GET, PUT, HEAD).

## Команды

| Команда | Описание |
|---|---|
| `npm run dev` | Dev-сервер (Vite, порт 5173) |
| `npm run build` | Production сборка |
| `npm run typecheck` | Проверка типов TS (tsc --noEmit) |
| `npm run lint:ts` | ESLint проверка |
| `npm run prettier:fix` | Prettier форматирование |
| `npm run pre-commit` | Полная проверка перед коммитом (typecheck → lint → test) |