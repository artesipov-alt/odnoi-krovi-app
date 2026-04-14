import { bot, pinologger } from "../../instances";

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

  if (!RecipientProviderMaxID || RecipientProviderMaxID.trim() === "") {
    pinologger.warn(
      { donorName: DonorName },
      "RecipientProviderMaxID is empty, skipping notification",
    );
    return;
  }

  try {
    const message = `Донор (${DonorName}, группа ${DonorBloodGroup}) отказался от донации. Можете найти нового донора на Портале.`;

    await bot.api.sendMessageToUser(Number(RecipientProviderMaxID), message);

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
