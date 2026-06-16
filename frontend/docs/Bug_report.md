
## Bug Report: Загрузка фото в Telegram Mini App (Android WebView)

### Окружение

- **Приложение**: React SPA (Vite, TypeScript) внутри Telegram WebView / Max WebView
- **Механизм загрузки (3 шага)**:
  1. `POST /v1/uploads/presign/{id}` — получение presigned PUT URL для S3 (VK Cloud)
  2. `PUT {presigned_url}` — загрузка файла напрямую в S3
  3. `POST /v1/uploads/confirm` — подтверждение
- **Бэкенд**: Go-монолит (Huma v2), presigned URL через AWS SDK v2, ContentType не проверяется (закомментирован)
- **Файловый инпут**: `<input type='file' accept='image/*' />` в компоненте `ImgEditor`

### Симптоматика

| Платформа | Результат |
|---|---|
| iPhone, Telegram Mini App | ✅ Работает |
| Android, Max Mini App | ✅ Камера, файлы, галерея — всё работает |
| Android, Telegram Mini App | ❌ Системный диалог не показывается → сразу галерея Telegram. После выбора — ошибка |

### Два типа ошибок

**Тип 1 — "протухший content:// URI":**
```
"The requested file could not be read, typically due to permission problems
that have occurred after a reference to a file was acquired."
```
Android выдаёт временный `content://` URI — разрешение отзывается после закрытия пикера.

**Тип 2 — "Failed to fetch" (PUT к S3):**
```
"Failed to fetch"
```
`fetch(url, { method: 'PUT', body: file })` к S3 presigned URL падает.
При этом `getPhotoLink` отрабатывает ✅, файл читается ✅, проблема именно на шаге PUT.

Ошибка **интермиттентная**: из 5 попыток 4 fail, 1 success — паттерн CORS preflight-кеширования.

### Хронология попыток

| # | Подход | Результат |
|---|---|---|
| 1 | `fetch(url, { method: 'PUT', body: File })` | ❌ |
| 2 | `XHR` + `FileReader` → `ArrayBuffer` | ❌ |
| 3 | Вернулись к `fetch(body: File)` | ❌ |
| 4 | `XHR` + `ArrayBuffer` (повторно) | ❌ |
| 5 | Убирали/возвращали `accept='image/*'` | ❌ Telegram всё равно перехватывает пикер |
| 6 | Fallback `Content-Type` по расширению | ❌ |
| 7 | Вебхук: `fetch` → `sendBeacon` → `Image()` для отладки | `Image()` работает, остальные нет |
| 8 | `arrayBuffer()` сразу в `onChange` + in-memory `File` | ✅ в Max, ❌ в Telegram |
| 9 | `accept='image/*,application/pdf'` для системного пикера | ❌ не сработал из-за кеша |
| 10 | Повторный деплой → Telegram подтянул свежий билд | ✅ системный пикер появился! |
| 11 | На свежем билде: `arrayBuffer()` + системный пикер | ✅ загрузка работает! |

### Корневая причина

**Две проблемы, обе специфичны для Telegram Android WebView:**

1. **`accept='image/*'`** — Telegram перехватывает файловый инпут, показывает свою галерею вместо системного диалога. Это не даёт выбрать камеру/файлы, а также галерея Telegram отдаёт `content://` URI, который может протухнуть.

2. **CORS preflight (`OPTIONS`) к S3** — `PUT` с `Content-Type: image/jpeg` требует preflight. Telegram WebView имеет более строгую CORS-политику, чем Max WebView или обычный Chrome. Интермиттентность (4 fail / 1 success) — preflight иногда кешируется и проходит, иногда нет. Вероятно, S3 CORS на бакете не настроен.

### Финальное решение (что в коде сейчас)

**1. Системный пикер — `accept='image/*,application/pdf'`**

Файлы: `ImgEditor/index.tsx`, `profile/index.tsx`
```html
<input type='file' accept='image/*,application/pdf' />
```
`application/pdf` заставляет Telegram WebView отдать системный пикер вместо своей галереи. При этом PDF отфильтровываются — принимаются только изображения:
```ts
if (!newFile.type.startsWith('image/')) return;
```

**2. Защита от протухающего URI — `arrayBuffer()` сразу в `onChange`**

Файлы: `ImgEditor/index.tsx`, `profile/index.tsx`
```ts
const buffer = await newFile.arrayBuffer();
const safeFile = new File([buffer], newFile.name, { type: newFile.type || 'image/jpeg' });
```
Чтение файла в память немедленно, пока пикер открыт и `content://` URI ещё валиден.

**3. PUT без Content-Type — упрощение CORS preflight**

Файл: `api/apiServices/addPhoto.ts`
```ts
const buffer = await file.arrayBuffer();
const res = await fetch(url, { method: 'PUT', body: buffer });
// БЕЗ Content-Type заголовка
```
Presigned URL генерируется без проверки ContentType (закомментирован в бэкенде). Без кастомных заголовков preflight проще — не нужен `Access-Control-Allow-Headers`.

**4. XHR-fallback — оставлен в коде, не вызывается**

Функция `putFileViaXHR` сохранена на случай, если `fetch` PUT снова начнёт падать в Telegram WebView. XHR без `setRequestHeader` — минимальные требования к CORS. Если понадобится — достаточно раскомментировать fallback в `addPhoto`:
```ts
try { await putFile(...); } catch { await putFileViaXHR(...); }
```

### Что нужно проверить на стороне инфраструктуры

**S3 CORS на бакете** — если ошибки вернутся, настроить CORS на бакете VK Cloud:

```xml
<CORSConfiguration>
  <CORSRule>
    <AllowedOrigin>https://dev.1krovi.app</AllowedOrigin>
    <AllowedOrigin>https://1krovi.app</AllowedOrigin>
    <AllowedMethod>GET</AllowedMethod>
    <AllowedMethod>PUT</AllowedMethod>
    <AllowedMethod>HEAD</AllowedMethod>
    <MaxAgeSeconds>3600</MaxAgeSeconds>
  </CORSRule>
</CORSConfiguration>
```

### Файлы, затронутые изменениями

- `frontend/src/components/ImgEditor/index.tsx` — `accept`, `arrayBuffer()`, фильтр image
- `frontend/src/pages/profile/index.tsx` — `accept`, `arrayBuffer()`
- `frontend/src/api/apiServices/addPhoto.ts` — `putFile` без Content-Type, XHR-fallback в резерве
