import { pinologger } from "../../instances";
import { sendMessageToUser } from "../../max";
import { getAppOpenKeyboard } from "../../keyboards";

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

  for (const donor of AvilableDonors) {
    const targetId = donor.MaxID;

    if (!targetId || targetId.trim() === "") {
      pinologger.warn({ donor }, "Donor MaxID is empty, skipping notification");
      continue;
    }

    try {
      const message = `Питомцам нужна ваша помощь!\n\nНажмите "Стать донором" в приложении, чтобы узнать детали.`;

      const keyboard = getAppOpenKeyboard();

      await sendMessageToUser(targetId, message, {
        attachments: [keyboard],
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
