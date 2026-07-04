import { Keyboard } from "@maxhub/max-bot-api";

const PROD_BOT_ID = "id3200014662";
const DEV_BOT_ID = "id3200014662_2";

const BOT_ID = process.env.ENV === "development" ? DEV_BOT_ID : PROD_BOT_ID;

export const getMainKeyboard = () => {
  return Keyboard.inlineKeyboard([
    [
      Keyboard.button.link(
        "🩸 Открыть приложение",
        `https://max.ru/${BOT_ID}_bot?startapp`,
      ),
    ],
    // [
    //   Keyboard.button.callback("❓ Помощь", "help"),
    //   Keyboard.button.callback("👤 Профиль", "profile"),
    // ],
  ]);
};

export const getAppOpenKeyboard = () => {
  return Keyboard.inlineKeyboard([
    [
      Keyboard.button.link(
        "🩸 Открыть приложение",
        `https://max.ru/${BOT_ID}_bot?startapp`,
      ),
    ],
  ]);
};

/**
 * Клавиатура для уведомления recipient_empty_showcase.
 * Две кнопки: "Да" и "Нет". requestId передаётся через callbackData.
 */
export const getNotificationKeyboard = (requestId: string) => {
  return Keyboard.inlineKeyboard([
    [
      Keyboard.button.callback("✅ Да", `notification_yes_${requestId}`),
      Keyboard.button.callback("❌ Нет", `notification_no_${requestId}`),
    ],
  ]);
};
