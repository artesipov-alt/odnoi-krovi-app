import { pinologger } from "../../instances";
import { sendMessageToUser } from "../../max";
import { getAppOpenKeyboard } from "../../keyboards";

interface DonorCancelEvent {
  donorName: string;
  donorBloodGroup: string;
  recipientProviderMaxId: string;
  recipientPetName: string;
  createdAt: string;
}

export const handleDonorCancel = async (event: DonorCancelEvent) => {
  const {
    donorName,
    donorBloodGroup,
    recipientProviderMaxId,
    recipientPetName,
  } = event;

  const donorBloodGroupDisplay =
    donorBloodGroup === "UNKNOWN" ? "не определена" : donorBloodGroup;

  if (!recipientProviderMaxId || recipientProviderMaxId.trim() === "") {
    pinologger.warn(
      { donorName },
      "recipientProviderMaxId is empty, skipping notification",
    );
    return;
  }

  try {
    const message = `Донор (${donorName}, группа ${donorBloodGroupDisplay}) отказался от донации. Можете найти нового донора на Портале.`;

    await sendMessageToUser(recipientProviderMaxId, message, {
      attachments: [getAppOpenKeyboard()],
    });

    pinologger.info(
      {
        recipientId: recipientProviderMaxId,
        donorName,
      },
      "Sent donor cancel notification",
    );
  } catch (err) {
    pinologger.error(
      { error: err },
      "Failed to send donor cancel notification",
    );
  }
};
