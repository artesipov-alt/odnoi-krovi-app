import { pinologger } from "../../instances";
import { sendMessageToUser } from "../../max";
import { getAppOpenKeyboard } from "../../keyboards";

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
  const { bloodTypes, regions, maxId } = event;

  if (!maxId || maxId.trim() === "") {
    pinologger.warn({ event }, "Donor maxId is empty, skipping notification");
    return;
  }

  try {
    const message = `Питомцам нужна ваша помощь!\n\nНажмите "Стать донором" в приложении, чтобы узнать детали.`;

    const keyboard = getAppOpenKeyboard();

    await sendMessageToUser(maxId, message, {
      attachments: [keyboard],
    });

    pinologger.info(
      {
        targetId: maxId,
        bloodTypes,
        regions,
      },
      "Sent new blood request notification to donor",
    );
  } catch (err) {
    pinologger.error(
      { error: err, targetId: maxId },
      "Failed to send new blood request notification",
    );
  }
};
