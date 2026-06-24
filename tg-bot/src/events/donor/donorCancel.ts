import { pinologger } from "../../instances";
import { sendTelegramMessage } from "../../telegram";
import { createOpenAppKeyboard } from "../../telegramButtons";

interface DonorCancelEvent {
  DonorName: string;
  DonorBloodGroup: string;
  RecipientProviderTelegramID: string;
  RecipientPetName: string;
  CreatedAt: string;
}

export const handleDonorCancel = async (event: DonorCancelEvent) => {
  const {
    DonorName,
    DonorBloodGroup,
    RecipientProviderTelegramID,
    RecipientPetName,
  } = event;

  const donorBloodGroup =
    DonorBloodGroup === "UNKNOWN" ? "не определена" : DonorBloodGroup;

  if (
    !RecipientProviderTelegramID ||
    RecipientProviderTelegramID.trim() === ""
  ) {
    pinologger.warn(
      { donorName: DonorName },
      "RecipientProviderTelegramID is empty, skipping notification",
    );
    return;
  }

  try {
    const message = `Донор (${DonorName}, группа ${donorBloodGroup}) отказался от донации. Можете найти нового донора на Портале.`;

    await sendTelegramMessage(RecipientProviderTelegramID, message, {
      reply_markup: createOpenAppKeyboard(),
    });

    pinologger.info(
      {
        recipientId: RecipientProviderTelegramID,
        donorName: DonorName,
      },
      "Sent donor cancel notification",
    );
  } catch (err) {
    pinologger.error(
      { error: err, recipientId: RecipientProviderTelegramID },
      "Failed to send donor cancel notification",
    );
  }
};
