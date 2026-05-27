import { pinologger } from "../../instances";
import { sendMessageToUser } from "../../max";
import { getAppOpenKeyboard } from "../../keyboards";

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
  RecipientData: {
    PetName: string;
    BloodGroup: string;
  };
  Volume: number; // Объем донации в мл
}

export const handleDonationConfirmed = async (
  event: DonationConfirmedEvent,
) => {
  const { DonorData, RecipientData, Volume } = event;

  const targetId = DonorData.ProviderMaxID;

  if (!targetId || targetId.trim() === "") {
    pinologger.warn(
      { donorUserName: DonorData.UserName },
      "Donor ProviderMaxID is empty, skipping notification",
    );
    return;
  }

  try {
    const message = generateDonationMessage({
      volume: Volume,
      recipientPetName: RecipientData.PetName,
      recipientBloodGroup: RecipientData.BloodGroup,
    });

    await sendMessageToUser(targetId, message, {
      attachments: [getAppOpenKeyboard()],
    });

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
