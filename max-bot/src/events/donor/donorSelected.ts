import { pinologger } from "../../instances";
import { sendMessageToUser } from "../../max";
import { getAppOpenKeyboard } from "../../keyboards";

// Реципиент выбрал донора из списка потенциальных.
// Уведомление отправляется только донору.
interface DonorSelectedEvent {
  donorData: {
    providerMaxId: string;
    providerTelegram: string;
  };
  recipientData: {
    petName: string;
    petType: string; // "cat" или "dog"
    volume: number; // мл
    bloodGroup: string;
  };
}

const petTypeLabel = (petType: string): string => {
  switch (petType) {
    case "cat":
      return "кошка";
    case "dog":
      return "собака";
    default:
      return "питомец";
  }
};

export const handleDonorSelected = async (event: DonorSelectedEvent) => {
  const { donorData, recipientData } = event;

  const donorProviderMaxID = donorData.providerMaxId;

  if (!donorProviderMaxID || donorProviderMaxID.trim() === "") {
    pinologger.warn("Donor providerMaxId is empty, skipping notification");
    return;
  }

  try {
    const typeLabel = petTypeLabel(recipientData.petType);
    const bloodGroup =
      recipientData.bloodGroup === "UNKNOWN"
        ? "не определена"
        : recipientData.bloodGroup;

    const message = `На Портале ищут донора для питомца - ${typeLabel} ${recipientData.petName} (${recipientData.volume} мл, группа ${bloodGroup})\n\nРанее Вы разрешили связаться с Вами, если нужна помощь.\nХозяин реципиента получил Ваши контакты. Дождитесь, пока с Вами свяжутся, или напишите хозяину реципиента`;

    await sendMessageToUser(donorProviderMaxID, message, {
      attachments: [getAppOpenKeyboard()],
    });

    pinologger.info(
      {
        donorId: donorProviderMaxID,
        recipientPetName: recipientData.petName,
      },
      "Sent donor selected notification",
    );
  } catch (err) {
    pinologger.error(
      { error: err, donorId: donorProviderMaxID },
      "Failed to send donor selected notification",
    );
  }
};