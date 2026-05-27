import { pinologger } from "../../instances";
import { sendMessageToUser } from "../../max";

interface DonorCancelEvent {
  DonorName: string;
  DonorBloodGroup: string;
  RecipientProviderMaxID: string;
  RecipientPetName: string;
  CreatedAt: string;
}

export const handleDonorCancel = async (event: DonorCancelEvent) => {
  const {
    DonorName,
    DonorBloodGroup,
    RecipientProviderMaxID,
    RecipientPetName,
  } = event;

  const donorBloodGroup =
    DonorBloodGroup === "UNKNOWN" ? "не определена" : DonorBloodGroup;

  if (!RecipientProviderMaxID || RecipientProviderMaxID.trim() === "") {
    pinologger.warn(
      { donorName: DonorName },
      "RecipientProviderMaxID is empty, skipping notification",
    );
    return;
  }

  try {
    const message = `Донор (${DonorName}, группа ${donorBloodGroup}) отказался от донации. Можете найти нового донора на Портале.`;

    await sendMessageToUser(RecipientProviderMaxID, message);

    pinologger.info(
      {
        recipientId: RecipientProviderMaxID,
        donorName: DonorName,
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
