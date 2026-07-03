import { pinologger } from "../../instances";
import { sendTelegramMessage, sendTelegramContact } from "../../telegram";

export interface UserContactEvent {
  notifyProvider: string;
  sendTo: string;
  userData: {
    name: string;
    providerMaxId: string;
    providerTelegram: string;
    phone: string;
  };
  createdAt: string;
  recipient: any;
}

export const handleUserContact = async (event: UserContactEvent) => {
  const { notifyProvider, sendTo, userData } = event;

  // Обрабатываем только события, предназначенные для Telegram Bot
  if (notifyProvider !== "telegram_bot") {
    pinologger.warn(
      { notifyProvider },
      "notifyProvider is not telegram_bot, skipping",
    );
    return;
  }

  if (!sendTo || sendTo.trim() === "") {
    pinologger.warn(
      { notifyProvider },
      "sendTo is empty, skipping notification",
    );
    return;
  }

  try {
    // Отправляем контакт (если есть телефон)
    if (userData.phone) {
      const nameParts = userData.name.split(" ");
      await sendTelegramContact(
        sendTo,
        userData.phone,
        nameParts[0] || userData.name,
        { last_name: nameParts.slice(1).join(" ") || undefined },
      );
    } else {
      // Если нет телефона, отправляем просто текст
      const message = `Контакт пользователя: ${userData.name}`;
      await sendTelegramMessage(sendTo, message);
    }

    pinologger.info(
      { sendTo, userName: userData.name },
      "Sent user contact notification",
    );
  } catch (err) {
    pinologger.error(
      { error: err },
      "Failed to send user contact notification",
    );
  }
};
