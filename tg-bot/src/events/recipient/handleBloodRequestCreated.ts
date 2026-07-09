import { pinologger } from "../../instances";
import { sendTelegramMessage } from "../../telegram";
import { createOpenAppKeyboard } from "../../telegramButtons";

interface BloodRequestCreatedEvent {
  requestId: string;
  bloodTypes: string[];
  regions: string[];
  telegramId: string;
  maxId: string;
  createdAt: string;
}

export const handleBloodRequestCreated = async (
  event: BloodRequestCreatedEvent,
) => {
  pinologger.info({ event }, "Received blood_request_created event");
  const { bloodTypes, regions, telegramId } = event;

  if (!telegramId || telegramId.trim() === "") {
    pinologger.warn(
      { event },
      "Donor telegramId is empty, skipping notification",
    );
    return;
  }

  try {
    const message = `Питомцам нужна ваша помощь!\n\nНажмите "Стать донором" в приложении, чтобы узнать детали.`;

    await sendTelegramMessage(telegramId, message, {
      reply_markup: createOpenAppKeyboard(),
    });

    pinologger.info(
      {
        targetId: telegramId,
        bloodTypes,
        regions,
      },
      "Sent new blood request notification to donor",
    );
  } catch (err) {
    pinologger.error(
      { error: err, targetId: telegramId },
      "Failed to send new blood request notification",
    );
  }
};
