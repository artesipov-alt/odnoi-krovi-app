# Changelog

Все заметные изменения в **max-bot** фиксируются в этом файле.

Формат основан на [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
и проект следует [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.6.0] — 2026-06-24

### Added

- **Кнопка «🩸 Открыть приложение»** добавлена на все уведомления (7 событий):
  - Отклик донора на реципиента (`donorApply`)
  - Отклик реципиента на донора (`recipientApply`)
  - Отмена донации донором (`donorCancel`)
  - Отклонение донации реципиентом (`donorReject`)
  - Неподтверждение донации (`donorNotConfirmed`)
  - Завершение донации донором (`donorCompleted`)
  - Передача контакта пользователя (`handleUserContact`)
- Используется существующая функция `getAppOpenKeyboard()` из `src/keyboards.ts`.

### Changed

- **Разделение контакта и кнопки** — в сценариях, где уведомление содержит и contact-attachment, и кнопку, отправка разбита на два последовательных вызова `sendMessageToUser`:
  1. Текст-уведомление + кнопка открытия приложения
  2. Пустое сообщение + contact (VCF)
  Затронутые файлы: `donorApply.ts`, `donorNotConfirmed.ts`, `handleUserContact.ts`.
