import { InlineKeyboard } from "grammy";

/**
 * Создает кнопку для перехода в чат с пользователем по его Telegram ID
 * (Не используется в текущей реализации, так как контакты отправляются через sendContact)
 */
export const createUserChatButton = (telegramId: string) => {
  return new InlineKeyboard().url(
    "Написать пользователю",
    `tg://user?id=${telegramId}`,
  );
};
