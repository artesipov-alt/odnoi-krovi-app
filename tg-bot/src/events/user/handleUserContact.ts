import { pinologger } from "../../instances";
import { sendTelegramMessage, sendTelegramContact } from "../../telegram";

export interface UserContactEvent {
  NotifyProvider: string;
  SendTo: string;
  UserData: {
    Name: string;
    ProviderMaxID: string;
    ProviderTelegram: string;
    Phone: string;
  };
  CreatedAt: string;
  Recipient: any;
}

export const handleUserContact = async (event: UserContactEvent) => {
  const { NotifyProvider, SendTo, UserData } = event;

  // Обрабатываем только события, предназначенные для Telegram Bot
  if (NotifyProvider !== "telegram_bot") {
    pinologger.warn(
      { notifyProvider: NotifyProvider },
      "NotifyProvider is not telegram_bot, skipping",
    );
    return;
  }

  if (!SendTo || SendTo.trim() === "") {
    pinologger.warn(
      { notifyProvider: NotifyProvider },
      "SendTo is empty, skipping notification",
    );
    return;
  }

  try {
    // Отправляем контакт (если есть телефон)
    if (UserData.Phone) {
      const nameParts = UserData.Name.split(" ");
      await sendTelegramContact(
        SendTo,
        UserData.Phone,
        nameParts[0] || UserData.Name,
        { last_name: nameParts.slice(1).join(" ") || undefined },
      );
    } else {
      // Если нет телефона, отправляем просто текст
      const message = `Контакт пользователя: ${UserData.Name}`;
      await sendTelegramMessage(SendTo, message);
    }

    pinologger.info(
      { sendTo: SendTo, userName: UserData.Name },
      "Sent user contact notification",
    );
  } catch (err) {
    pinologger.error(
      { error: err },
      "Failed to send user contact notification",
    );
  }
};
