import { pinologger } from "../../instances";
import { sendTelegramMessage } from "../../telegram";
import { createOpenAppKeyboard } from "../../telegramButtons";

interface DonorRejectEvent {
  RecipientPetName: string;
  RecipientBloodGroup: string;
  DonorProviderTelegramID: string;
  DonorPetName: string;
  RejectedReason: string;
  CreatedAt: string;
}

export const handleDonorReject = async (event: DonorRejectEvent) => {
  const { RecipientPetName, RecipientBloodGroup, DonorProviderTelegramID } =
    event;

  if (!DonorProviderTelegramID || DonorProviderTelegramID.trim() === "") {
    pinologger.warn(
      { recipientPetName: RecipientPetName },
      "DonorProviderTelegramID is empty, skipping notification",
    );
    return;
  }

  try {
    const message = `Реципиент (${RecipientPetName}, группа ${RecipientBloodGroup}) сообщил об отмене донации. Можете помочь другим реципиентам на Портале.`;

    await sendTelegramMessage(DonorProviderTelegramID, message, {
      reply_markup: createOpenAppKeyboard(),
    });

    pinologger.info(
      {
        donorId: DonorProviderTelegramID,
        recipientPetName: RecipientPetName,
      },
      "Sent donor reject notification",
    );
  } catch (err) {
    pinologger.error(
      { error: err, donorId: DonorProviderTelegramID },
      "Failed to send donor reject notification",
    );
  }
};
