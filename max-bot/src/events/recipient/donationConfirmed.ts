import { bot, pinologger } from "../../instances";

import { generateDonationMessage } from "./helpers";

// Уведомление о подтвержденной донации (от реципиента донору).
interface DonationConfirmedEvent {
  DonorData: {
    UserName: string;
    PetName: string;
    ProviderMaxID: string;
    ProviderTelegram: string;
    Phone: string;
    BloodGroup: string;
  };
  Volume: number; // Объем донации в мл
}

export const handleDonationConfirmed = async (
  event: DonationConfirmedEvent,
) => {
  const { DonorData, Volume } = event;

  let targetId = DonorData.ProviderMaxID;
  if (!targetId || targetId.trim() === "") {
    targetId = DonorData.ProviderTelegram;
  }

  if (!targetId || targetId.trim() === "") {
    pinologger.warn(
      { donorUserName: DonorData.UserName },
      "Donor ProviderMaxID and ProviderTelegram are empty, skipping notification",
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
