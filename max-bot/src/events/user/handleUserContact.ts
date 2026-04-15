import { bot, pinologger } from "../../instances";

interface UserContactEvent {
  NotifyProvider: string;
  SendTo: string;
  UserData: {
    Name: string;
    ProviderMaxID: string;
    ProviderTelegram: string;
    Phone: string;
  };
  CreatedAt: string;
  Recipient: any; // Assuming it's the user object, but not used
}

export const handleUserContact = async (event: UserContactEvent) => {
  const { NotifyProvider, SendTo, UserData } = event;

  if (
    (NotifyProvider !== "telegram_bot" && NotifyProvider !== "max_bot") ||
    !SendTo ||
    SendTo.trim() === ""
  ) {
    pinologger.warn(
      { notifyProvider: NotifyProvider, sendTo: SendTo },
      "Invalid provider or SendTo is empty, skipping notification",
    );
    return;
  }

  try {
    const message = `Контакт пользователя: ${UserData.Name}\nТелефон: ${UserData.Phone}`;

    await bot.api.sendMessageToUser(Number(SendTo), message);

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
