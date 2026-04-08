import { bot, pinologger } from "../../instances";

import { generateMessage } from "./helpers";

// Отклик донора на рецепиента.
interface RecipientApplyEvent {
  DonorName: string;
  DonorBloodGroup: string;
  RecipientProviderMaxID: string;
  RecipientPetName: string;
  RecipientPetSearchingBloodGroup: string[];
  RecipientPetNeededVolume: number;
  CreatedAt: string;
}

export const handleRecipientApply = async (event: RecipientApplyEvent) => {
  const {
    DonorName,
    DonorBloodGroup,
    RecipientProviderMaxID,
    RecipientPetName,
    RecipientPetSearchingBloodGroup,
    RecipientPetNeededVolume,
  } = event;

  if (!RecipientProviderMaxID || RecipientProviderMaxID.trim() === "") {
    pinologger.warn(
      { donorName: DonorName },
      "RecipientProviderMaxID is empty, skipping notification",
    );
    return;
  }

  try {
    const message = generateMessage({
      recipientPetName: RecipientPetName,
      recipientPetNeededVolume: RecipientPetNeededVolume,
      recipientPetSearchingBloodGroup: RecipientPetSearchingBloodGroup,
      donorName: DonorName,
      donorBloodGroup: DonorBloodGroup,
    });

    await bot.api.sendMessageToUser(Number(RecipientProviderMaxID), message);

    pinologger.info(
      {
        recipientId: RecipientProviderMaxID,
        donorName: DonorName,
      },
      "Sent recipient apply notification",
    );
  } catch (err) {
    pinologger.error(
      { error: err },
      "Failed to send recipient apply notification",
    );
  }
};
