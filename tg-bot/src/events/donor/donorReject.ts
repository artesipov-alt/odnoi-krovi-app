import { pinologger } from "../../instances";
import { sendTelegramMessage } from "../../telegram";
import { createOpenAppKeyboard } from "../../telegramButtons";

interface DonorRejectEvent {
  recipientPetName: string;
  recipientBloodGroup: string;
  donorProviderTelegramId: string;
  donorPetName: string;
  rejectedReason: string;
  createdAt: string;
}

export const handleDonorReject = async (event: DonorRejectEvent) => {
  const { recipientPetName, recipientBloodGroup, donorProviderTelegramId } =
    event;

  if (!donorProviderTelegramId || donorProviderTelegramId.trim() === "") {
    pinologger.warn(
      { recipientPetName },
      "donorProviderTelegramId is empty, skipping notification",
    );
    return;
  }

  try {
    const message = `Реципиент (${recipientPetName}, группа ${recipientBloodGroup}) сообщил об отмене донации. Можете помочь другим реципиентам на Портале.`;

    await sendTelegramMessage(donorProviderTelegramId, message, {
      reply_markup: createOpenAppKeyboard(),
    });

    pinologger.info(
      {
        donorId: donorProviderTelegramId,
        recipientPetName,
      },
      "Sent donor reject notification",
    );
  } catch (err) {
    pinologger.error(
      { error: err, donorId: donorProviderTelegramId },
      "Failed to send donor reject notification",
    );
  }
};
