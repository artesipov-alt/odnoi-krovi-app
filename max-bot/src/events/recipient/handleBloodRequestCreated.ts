import { bot, pinologger } from "../../instances";

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
    let targetId = donor.MaxID;
    if (!targetId || targetId.trim() === "") {
      targetId = donor.TelegramID;
    }

    if (!targetId || targetId.trim() === "") {
      pinologger.warn(
        { donor },
        "Donor MaxID and TelegramID are empty, skipping notification",
      );
      continue;
    }

    try {
      const message = `Появилась новая заявка на донацию крови (${BloodTypes.join(", ")}) в регионах: ${Regions.join(", ")}. Проверьте на Портале.`;

      await bot.api.sendMessageToUser(Number(targetId), message);

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
