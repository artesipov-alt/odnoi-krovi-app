import { bot, pinologger } from "../../instances";

import { generateDonationMessage } from "./helpers";

// Уведомление о подтвержденной донации (от реципиента донору).
interface DonationConfirmedEvent {
  DonorData: {
    ProviderMaxID: string;
    Name: string;
  };
  Volume: number; // Объем донации в мл
}

export const handleDonationConfirmed = async (
  event: DonationConfirmedEvent,
) => {
  const { DonorData, Volume } = event;

  const targetId = DonorData.ProviderMaxID;

  if (!targetId || targetId.trim() === "") {
    pinologger.warn(
      { donorId: DonorData.ProviderMaxID },
      "Donor ID is empty, skipping notification",
    );
    return;
  }

  try {
    const message = generateDonationMessage({
      volume: Volume,
    });

    await bot.api.sendMessageToUser(Number(targetId), message);

    pinologger.info(
      {
        targetId,
        volume: Volume,
      },
      "Sent donation confirmed notification to donor",
    );
  } catch (err) {
    pinologger.error(
      { error: err, targetId },
      "Failed to send donation confirmed notification",
    );
  }
};
