import { bot, pinologger } from "../../instances";

import { generateDonationMessage } from "./helpers";

// Уведомление о состоявшейся донации.
interface DonationCompletedEvent {
  Initiator: "donor" | "recipient"; // Кто нажал "Донация состоялась"
  DonorData: {
    ProviderMaxID: string;
    Name: string;
  };
  RecipientData: {
    ProviderMaxID: string;
    Name: string;
  };
  Volume: number; // Объем донации в мл
}

export const handleDonationCompleted = async (
  event: DonationCompletedEvent,
) => {
  const { Initiator, DonorData, RecipientData, Volume } = event;

  const recipientId = RecipientData.ProviderMaxID;
  const donorId = DonorData.ProviderMaxID;

  // Определяем, кому отправить уведомление: если инициатор донор, то реципиенту, и наоборот
  const targetId = Initiator === "donor" ? recipientId : donorId;
  const targetName = Initiator === "donor" ? "recipient" : "donor";

  if (!targetId || targetId.trim() === "") {
    pinologger.warn(
      { initiator: Initiator },
      `${targetName} ID is empty, skipping notification`,
    );
    return;
  }

  try {
    const message = generateDonationMessage({
      initiator: Initiator,
      volume: Volume,
    });

    await bot.api.sendMessageToUser(Number(targetId), message);

    pinologger.info(
      {
        targetId,
        initiator: Initiator,
        volume: Volume,
      },
      "Sent donation completed notification",
    );
  } catch (err) {
    pinologger.error(
      { error: err, targetId },
      "Failed to send donation completed notification",
    );
  }
};
