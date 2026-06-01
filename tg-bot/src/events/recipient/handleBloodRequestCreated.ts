import { pinologger } from "../../instances";
import { sendTelegramMessage } from "../../telegram";
import { InlineKeyboard } from "grammy";

interface BloodRequestCreatedEvent {
  RequestID: string;
  BloodTypes: string[];
  Regions: string[];
  AvilableDonors: Array<{
    TelegramID: string;
    MaxID: string;
  }>;
  CreatedAt: string;
}

export const handleBloodRequestCreated = async (
  event: BloodRequestCreatedEvent,
) => {
  pinologger.info({ event }, "Received blood_request_created event");
  const { BloodTypes, Regions, AvilableDonors } = event;

  // Определяем URL приложения в зависимости от среды
  const getWebAppUrl = () => {
    // Определяем среду по NODE_ENV или BUN_ENV
    const env = Bun.env.ENV || "production";
    const isDev = env === "development" || env === "dev";
    return isDev ? "https://dev.1krovi.app" : "https://1krovi.app";
  };

  // Создаем клавиатуру для открытия приложения
  const keyboard = new InlineKeyboard().webApp(
    "🩸 Открыть приложение",
    getWebAppUrl(),
  );

  for (const donor of AvilableDonors) {
    const targetId = donor.TelegramID;

    if (!targetId || targetId.trim() === "") {
      pinologger.warn(
        { donor },
        "Donor TelegramID is empty, skipping notification",
      );
      continue;
    }

    try {
      const message = `Питомцам нужна ваша помощь!\n\nНажмите "Стать донором" в приложении, чтобы узнать детали.`;

      await sendTelegramMessage(targetId, message, {
        reply_markup: keyboard,
      });

      pinologger.info(
        {
          targetId,
          bloodTypes: BloodTypes,
          regions: Regions,
        },
        "Sent new blood request notification to donor",
      );
    } catch (err) {
      pinologger.error(
        { error: err, targetId },
        "Failed to send new blood request notification",
      );
    }
  }
};
