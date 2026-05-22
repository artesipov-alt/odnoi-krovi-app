import { bot, pinologger } from "./instances";

/**
 * Безопасно отправляет сообщение в Telegram с валидацией chat_id.
 * Возвращает `null`, если ID невалидный или пустой.
 */
export const sendTelegramMessage = async (
  telegramId: string | undefined | null,
  text: string,
  extra?: Parameters<typeof bot.api.sendMessage>[2],
) => {
  if (!telegramId?.trim()) return null;

  const chatId = Number(telegramId);
  if (isNaN(chatId)) {
    pinologger.warn({ telegramId }, "Invalid numeric Telegram ID, skipping");
    return null;
  }

  return bot.api.sendMessage(chatId, text, extra);
};

/**
 * Безопасно отправляет контакт в Telegram с валидацией chat_id.
 * Возвращает `null`, если ID невалидный или пустой.
 */
export const sendTelegramContact = async (
  telegramId: string | undefined | null,
  phone: string,
  firstName: string,
  extra?: { last_name?: string },
) => {
  if (!telegramId?.trim()) return null;

  const chatId = Number(telegramId);
  if (isNaN(chatId)) {
    pinologger.warn({ telegramId }, "Invalid numeric Telegram ID, skipping");
    return null;
  }

  return bot.api.sendContact(chatId, phone, firstName, extra);
};
