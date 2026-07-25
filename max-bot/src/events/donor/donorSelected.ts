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
    userName: string;
    petName: string;
    petType: string; // "cat" или "dog"
    volume: number; // мл
    bloodGroup: string;
    phone: string;
    providerMaxId: string;
    providerTelegram: string;
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

const generateVCF = (name: string, phone: string): string => {
  return `BEGIN:VCARD\r\nVERSION:3.0\r\nFN:${name}\r\nTEL:${phone}\r\nEND:VCARD`;
};

export const handleDonorSelected = async (event: DonorSelectedEvent) => {
  const { donorData, recipientData } = event;

  const donorProviderMaxID = donorData.providerMaxId;

  if (!donorProviderMaxID || donorProviderMaxID.trim() === "") {
    pinologger.warn("Donor providerMaxId is empty, skipping notification");
    return;
  }

  try {
    // Сначала отправляем контакт реципиента
    await sendMessageToUser(donorProviderMaxID, "", {
      attachments: [
        {
          type: "contact",
          payload: {
            name: recipientData.userName,
            contact_id: Number(recipientData.providerMaxId),
            vcf_phone: recipientData.phone,
            vcf_info: generateVCF(recipientData.userName, recipientData.phone),
          },
        },
      ],
    });

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