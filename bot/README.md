# 🤖 Telegram Bot - Одной Крови

<div align="center">

![Bun](https://img.shields.io/badge/Bun-1.2+-000000?style=for-the-badge&logo=bun&logoColor=white) ![TypeScript](https://img.shields.io/badge/TypeScript-5.0+-3178C6?style=for-the-badge&logo=typescript&logoColor=white) ![Grammy](https://img.shields.io/badge/Grammy-1.38+-00ADD8?style=for-the-badge&logo=telegram&logoColor=white)

</div>

**Telegram Bot для платформы "Одной Крови"** — помощник для поиска донорской крови для животных и управления профилем пользователя через Telegram.

## 🎯 Функциональность

### Основные возможности:
- **📱 Интеграция с Mini App** - Быстрый доступ к основному приложению
- **👤 Управление профилем** - Просмотр и редактирование данных пользователя
- **🔍 Поиск доноров** - Быстрый поиск через inline режим (в разработке)
- **🔔 Уведомления** - Оповещения о новых запросах крови
- **📊 Статистика** - Отслеживание активности пользователей

### Команды бота:
- `/start` - Запуск бота и основное меню
- `/help` - Справка по использованию
- `/profile` - Просмотр профиля пользователя

## 🏗️ Архитектура

```
bot/
├── src/
│   ├── handlers/              # Обработчики сообщений и команд
│   │   ├── commands.ts        # Основные команды бота
│   │   └── errors.ts          # Обработка ошибок
│   ├── middleware/            # Промежуточное ПО
│   │   ├── logger.ts          # Логирование запросов
│   │   └── ratelimitter.ts    # Ограничение запросов
│   ├── config/                # Конфигурация
│   │   └── templates.ts       # Шаблоны сообщений
│   ├── utils/                 # Вспомогательные функции
│   └── instances.ts           # Инициализация бота и логгера
├── bot.ts                     # Основной файл бота
├── build.ts                   # Скрипт сборки
├── pm2.config.cjs             # Конфигурация для PM2
└── package.json
```

## 🚀 Быстрый старт

### Предварительные требования

- **Bun 1.3+** - JavaScript runtime
- **Telegram Bot Token** от @BotFather
- **Node.js 18+** (опционально, для PM2)

### Установка и запуск

1. **Клонирование и настройка**
```bash
cd bot
```

2. **Переменные окружения**

Создайте файл `.env` со следующими переменными:

```env
DEV_BOT_API_KEY=your_telegram_bot_token
PROD_BOT_API_KEY=your_telegram_bot_token
MINIAPP_DOMAIN=https://your-miniapp-domain.com
```


3. **Установка зависимостей**
```bash
bun install
```

4. **Запуск в режиме разработки**
```bash
bun run dev
```

4. **Сборка для production**
```bash
bun run build
```

5. **Запуск собранной версии**
```bash
bun run preview
```

## 🛠️ Технологический стек

### Основные технологии:
- **Bun** - Быстрый JavaScript runtime
- **TypeScript** - Статическая типизация
- **Grammy** - Фреймворк для Telegram Bot API
- **Pino** - Структурированное логирование

### Дополнительные пакеты:
- `@grammyjs/ratelimiter` - Ограничение запросов
- `@grammyjs/runner` - Параллельное выполнение
- `@grammyjs/transformer-throttler` - Ограничение API вызовов
- `pino-pretty` - Красивое форматирование логов


## 🚀 Production развертывание

### Использование PM2

```bash
# Установка PM2
npm install -g pm2

# Запуск бота
pm2 start pm2.config.cjs

# Мониторинг
pm2 dash

# Логи
pm2 logs odnoi-krovi-bot
```

## 🔒 Безопасность

### Rate Limiting
- Ограничение: 2 запроса в 1.5 секунды
- Хранилище: In-memory (планируется Redis)
- Сообщение при превышении лимита

### Валидация данных
- Проверка всех входящих сообщений
- Валидация callback данных
- Обработка ошибок на всех уровнях
