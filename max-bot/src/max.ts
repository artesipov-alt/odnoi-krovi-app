import { bot, pinologger } from "./instances";

/**
 * Безопасно отправляет сообщение пользователю через Max API с валидацией MaxID.
 * Возвращает `null`, если ID невалидный или пустой.
 */
export const sendMessageToUser = async (
  maxId: string | undefined | null,
  text: string,
  extra?: Parameters<typeof bot.api.sendMessageToUser>[2],
) => {
  if (!maxId?.trim()) return null;

  const id = Number(maxId);
  if (isNaN(id)) {
    pinologger.warn({ maxId }, "Invalid numeric MaxID, skipping");
    return null;
  }

  return bot.api.sendMessageToUser(id, text, extra);
};
