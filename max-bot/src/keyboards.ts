import { Keyboard } from "@maxhub/max-bot-api";

const BOT_ID = "id3200014662";

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
