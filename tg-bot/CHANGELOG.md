# Changelog

Все заметные изменения в **tg-bot** фиксируются в этом файле.

Формат основан на [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
и проект следует [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.1.0] — 2026-06-24

### Added

- **Кнопка «🩸 Открыть приложение»** добавлена на все уведомления (8 событий):
  - Отклик донора на реципиента (`donorApply`)
  - Отклик реципиента на донора (`recipientApply`)
  - Отмена донации донором (`donorCancel`)
  - Отклонение донации реципиентом (`donorReject`)
  - Неподтверждение донации (`donorNotConfirmed`)
  - Завершение донации донором (`donorCompleted`)
  - Подтверждение донации (`donationConfirmed`)
  - Новый запрос крови (`bloodRequestCreated`)

### Changed

- **Рефакторинг клавиатуры** — повторяющийся код создания inline-клавиатуры вынесен в единую функцию `createOpenAppKeyboard()` в `src/telegramButtons.ts`. URL выбирается автоматически по `Bun.env.ENV`:
  - dev → `https://dev.1krovi.app`
  - prod → `https://1krovi.app`
- **Унифицирован URL** — в `donationConfirmed.ts` был hardcoded `Bun.env.WEB_APP_URL || "https://app.1krovi.app"`, теперь используется общий механизм.

### Fixed

- **Интерполяция строк** — в `donorCancel.ts` и `donorReject.ts` строки сообщений были на двойных кавычках вместо обратных, из-за чего переменные не подставлялись. Исправлено.
