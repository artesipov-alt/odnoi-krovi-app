import { InlineKeyboard } from "grammy";

/**
 * Создает кнопку для перехода в чат с пользователем по его Telegram ID
 * (Не используется в текущей реализации, так как контакты отправляются через sendContact)
 */
/**
 * Создает кнопку "🩸 Открыть приложение" с WebApp ссылкой.
 * URL выбирается в зависимости от окружения (dev/prod).
 */
export const createOpenAppKeyboard = () => {
  const env = Bun.env.ENV || "production";
  const isDev = env === "development" || env === "dev";
  const webAppUrl = isDev ? "https://dev.1krovi.app" : "https://1krovi.app";

  return new InlineKeyboard().webApp(
    "🩸 Открыть приложение",
    webAppUrl,
  );
};

export const createUserChatButton = (telegramId: string) => {
  return new InlineKeyboard().url(
    "Написать пользователю",
    `tg://user?id=${telegramId}`,
  );
};
