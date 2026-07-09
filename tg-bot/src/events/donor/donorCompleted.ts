import { pinologger } from "../../instances";
import { sendTelegramMessage } from "../../telegram";
import { createOpenAppKeyboard } from "../../telegramButtons";

interface DonorCompletedEvent {
  donorPetName: string;
  donorBloodGroup: string;
  recipientProviderTelegramId: string;
  recipientPetName: string;
  amount: number;
  createdAt: string;
}

export const handleDonorCompleted = async (event: DonorCompletedEvent) => {
  const { donorPetName, donorBloodGroup, recipientProviderTelegramId } = event;

  if (
    !recipientProviderTelegramId ||
    recipientProviderTelegramId.trim() === ""
  ) {
    pinologger.warn(
      { donorPetName },
      "recipientProviderTelegramId is empty, skipping notification",
    );
    return;
  }

  try {
    const message = `Донор (${donorPetName}, группа ${donorBloodGroup}) сообщил, что Вы уже провели донацию.
Подтвердите донацию на Портале, чтобы донор получил бонусы за помощь.
Через 3 дня донация будет подтверждена автоматически.
Если донация еще не состоялась, можете отказаться и связаться с донором для уточнения деталей.`;

    await sendTelegramMessage(recipientProviderTelegramId, message, {
      reply_markup: createOpenAppKeyboard(),
    });

    pinologger.info(
      {
        recipientId: recipientProviderTelegramId,
        donorPetName,
      },
      "Sent donor completed notification",
    );
  } catch (err) {
    pinologger.error(
      { error: err, recipientId: recipientProviderTelegramId },
      "Failed to send donor completed notification",
    );
  }
};
