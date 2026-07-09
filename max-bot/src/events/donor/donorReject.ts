import { pinologger } from "../../instances";
import { sendMessageToUser } from "../../max";
import { getAppOpenKeyboard } from "../../keyboards";

interface DonorRejectEvent {
  recipientPetName: string;
  recipientBloodGroup: string;
  donorProviderMaxId: string;
  donorPetName: string;
  rejectedReason: string;
  createdAt: string;
}

export const handleDonorReject = async (event: DonorRejectEvent) => {
  const { recipientPetName, recipientBloodGroup, donorProviderMaxId } = event;

  if (!donorProviderMaxId || donorProviderMaxId.trim() === "") {
    pinologger.warn(
      { recipientPetName },
      "donorProviderMaxId is empty, skipping notification",
    );
    return;
  }

  try {
    const message = `Реципиент (${recipientPetName}, группа ${recipientBloodGroup}) сообщил об отмене донации. Можете помочь другим реципиентам на Портале.`;

    await sendMessageToUser(donorProviderMaxId, message, {
      attachments: [getAppOpenKeyboard()],
    });

    pinologger.info(
      {
        donorId: donorProviderMaxId,
        recipientPetName,
      },
      "Sent donor reject notification",
    );
  } catch (err) {
    pinologger.error(
      { error: err },
      "Failed to send donor reject notification",
    );
  }
};
