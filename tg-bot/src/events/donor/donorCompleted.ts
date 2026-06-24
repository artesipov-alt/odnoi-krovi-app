import { pinologger } from "../../instances";
import { sendTelegramMessage } from "../../telegram";
import { createOpenAppKeyboard } from "../../telegramButtons";

interface DonorCompletedEvent {
  DonorPetName: string;
  DonorBloodGroup: string;
  RecipientProviderTelegramID: string;
  RecipientPetName: string;
  Amount: number;
  CreatedAt: string;
}

export const handleDonorCompleted = async (event: DonorCompletedEvent) => {
  const { DonorPetName, DonorBloodGroup, RecipientProviderTelegramID } = event;

  if (
    !RecipientProviderTelegramID ||
    RecipientProviderTelegramID.trim() === ""
  ) {
    pinologger.warn(
      { donorPetName: DonorPetName },
      "RecipientProviderTelegramID is empty, skipping notification",
    );
    return;
  }

  try {
    const message = `Донор (${DonorPetName}, группа ${DonorBloodGroup}) сообщил, что Вы уже провели донацию.
Подтвердите донацию на Портале, чтобы донор получил бонусы за помощь.
Через 3 дня донация будет подтверждена автоматически.
Если донация еще не состоялась, можете отказаться и связаться с донором для уточнения деталей.`;

    await sendTelegramMessage(RecipientProviderTelegramID, message, {
      reply_markup: createOpenAppKeyboard(),
    });

    pinologger.info(
      {
        recipientId: RecipientProviderTelegramID,
        donorPetName: DonorPetName,
      },
      "Sent donor completed notification",
    );
  } catch (err) {
    pinologger.error(
      { error: err, recipientId: RecipientProviderTelegramID },
      "Failed to send donor completed notification",
    );
  }
};
