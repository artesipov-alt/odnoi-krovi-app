import { pinologger } from "../../instances";
import { sendMessageToUser } from "../../max";

interface DonorCompletedEvent {
  DonorPetName: string;
  DonorBloodGroup: string;
  RecipientProviderMaxID: string;
  RecipientPetName: string;
  Amount: number;
  CreatedAt: string;
}

export const handleDonorCompleted = async (event: DonorCompletedEvent) => {
  const { DonorPetName, DonorBloodGroup, RecipientProviderMaxID } = event;

  if (!RecipientProviderMaxID || RecipientProviderMaxID.trim() === "") {
    pinologger.warn(
      { donorPetName: DonorPetName },
      "RecipientProviderMaxID is empty, skipping notification",
    );
    return;
  }

  try {
    const message = `Донор (${DonorPetName}, группа ${DonorBloodGroup}) сообщил, что Вы уже провели донацию.
Подтвердите донацию на Портале, чтобы донор получил бонусы за помощь.
Через 3 дня донация будет подтверждена автоматически.
Если донация еще не состоялась, можете отказаться и связаться с донором для уточнения деталей.`;

    await sendMessageToUser(RecipientProviderMaxID, message);

    pinologger.info(
      {
        recipientId: RecipientProviderMaxID,
        donorPetName: DonorPetName,
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
