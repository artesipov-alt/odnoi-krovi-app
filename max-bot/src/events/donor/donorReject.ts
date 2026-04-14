odnoi-krovi-app/max-bot/src/events/donor/donorReject.ts
import { bot, pinologger } from "../../instances";

interface DonorRejectEvent {
  RecipientPetName: string;
  RecipientBloodGroup: string;
  DonorProviderMaxID: string;
  DonorPetName: string;
  RejectedReason: string;
  CreatedAt: string;
}

export const handleDonorReject = async (event: DonorRejectEvent) => {
  const {
    RecipientPetName,
    RecipientBloodGroup,
    DonorProviderMaxID,
  } = event;

  if (!DonorProviderMaxID || DonorProviderMaxID.trim() === "") {
    pinologger.warn(
      { recipientPetName: RecipientPetName },
      "DonorProviderMaxID is empty, skipping notification",
    );
    return;
  }

  try {
    const message = `Реципиент (${RecipientPetName}, группа ${RecipientBloodGroup}) сообщил об отмене донации. Можете помочь другим реципиентам на Портале.`;

    await bot.api.sendMessageToUser(Number(DonorProviderMaxID), message);

    pinologger.info(
      {
        donorId: DonorProviderMaxID,
        recipientPetName: RecipientPetName,
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
