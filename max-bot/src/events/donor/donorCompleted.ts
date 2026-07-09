import { pinologger } from "../../instances";
import { sendMessageToUser } from "../../max";
import { getAppOpenKeyboard } from "../../keyboards";

interface DonorCompletedEvent {
  donorPetName: string;
  donorBloodGroup: string;
  recipientProviderMaxId: string;
  recipientPetName: string;
  amount: number;
  createdAt: string;
}

export const handleDonorCompleted = async (event: DonorCompletedEvent) => {
  const { donorPetName, donorBloodGroup, recipientProviderMaxId } = event;

  if (!recipientProviderMaxId || recipientProviderMaxId.trim() === "") {
    pinologger.warn(
      { donorPetName },
      "recipientProviderMaxId is empty, skipping notification",
    );
    return;
  }

  try {
    const message = `Донор (${donorPetName}, группа ${donorBloodGroup}) сообщил, что Вы уже провели донацию.
Подтвердите донацию на Портале, чтобы донор получил бонусы за помощь.
Через 3 дня донация будет подтверждена автоматически.
Если донация еще не состоялась, можете отказаться и связаться с донором для уточнения деталей.`;

    await sendMessageToUser(recipientProviderMaxId, message, {
      attachments: [getAppOpenKeyboard()],
    });

    pinologger.info(
      {
        recipientId: recipientProviderMaxId,
        donorPetName,
      },
      "Sent donor completed notification",
    );
  } catch (err) {
    pinologger.error(
      { error: err },
      "Failed to send donor completed notification",
    );
  }
};
