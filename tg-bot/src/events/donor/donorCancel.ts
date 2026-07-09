import { pinologger } from "../../instances";
import { sendTelegramMessage } from "../../telegram";
import { createOpenAppKeyboard } from "../../telegramButtons";

interface DonorCancelEvent {
  donorName: string;
  donorBloodGroup: string;
  recipientProviderTelegramId: string;
  recipientPetName: string;
  createdAt: string;
}

export const handleDonorCancel = async (event: DonorCancelEvent) => {
  const {
    donorName,
    donorBloodGroup,
    recipientProviderTelegramId,
    recipientPetName,
  } = event;

  const donorBloodGroupDisplay =
    donorBloodGroup === "UNKNOWN" ? "не определена" : donorBloodGroup;

  if (
    !recipientProviderTelegramId ||
    recipientProviderTelegramId.trim() === ""
  ) {
    pinologger.warn(
      { donorName },
      "recipientProviderTelegramId is empty, skipping notification",
    );
    return;
  }

  try {
    const message = `Донор (${donorName}, группа ${donorBloodGroupDisplay}) отказался от донации. Можете найти нового донора на Портале.`;

    await sendTelegramMessage(recipientProviderTelegramId, message, {
      reply_markup: createOpenAppKeyboard(),
    });

    pinologger.info(
      {
        recipientId: recipientProviderTelegramId,
        donorName,
      },
      "Sent donor cancel notification",
    );
  } catch (err) {
    pinologger.error(
      { error: err, recipientId: recipientProviderTelegramId },
      "Failed to send donor cancel notification",
    );
  }
};
