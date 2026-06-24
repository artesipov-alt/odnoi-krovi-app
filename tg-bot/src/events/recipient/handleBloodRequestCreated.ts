import { pinologger } from "../../instances";
import { sendTelegramMessage } from "../../telegram";
import { createOpenAppKeyboard } from "../../telegramButtons";

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

  const keyboard = createOpenAppKeyboard();

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
