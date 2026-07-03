import { pinologger } from "../../instances";
import { sendMessageToUser } from "../../max";
import { getAppOpenKeyboard } from "../../keyboards";

import { generateDonationMessage } from "./helpers";

// Уведомление о подтвержденной донации (от реципиента донору).
interface DonationConfirmedEvent {
  donorData: {
    userName: string;
    petName: string;
    providerMaxId: string;
    providerTelegram: string;
    phone: string;
    bloodGroup: string;
  };
  recipientData: {
    petName: string;
    bloodGroup: string;
  };
  volume: number;
}

export const handleDonationConfirmed = async (
  event: DonationConfirmedEvent,
) => {
  const { donorData, recipientData, volume } = event;

  const targetId = donorData.providerMaxId;

  if (!targetId || targetId.trim() === "") {
    pinologger.warn(
      { donorUserName: donorData.userName },
      "Donor providerMaxId is empty, skipping notification",
    );
    return;
  }

  try {
    const message = generateDonationMessage({
      volume,
      recipientPetName: recipientData.petName,
      recipientBloodGroup: recipientData.bloodGroup,
    });

    await sendMessageToUser(targetId, message, {
      attachments: [getAppOpenKeyboard()],
    });

    pinologger.info(
      {
        targetId,
        volume,
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
