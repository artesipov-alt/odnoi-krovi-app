import { pinologger } from "../../instances";
import { sendMessageToUser } from "../../max";
import { getAppOpenKeyboard } from "../../keyboards";

// DonorSelected — реципиент выбрал донора из списка потенциальных.
// Уведомления отправляются и донору, и реципиенту.
interface DonorSelectedEvent {
  donorData: {
    userName: string;
    petName: string;
    phone: string;
    bloodGroup: string;
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
  const recipientProviderMaxID = recipientData.providerMaxId;

  // ── Донору: контакт реципиента + уведомление ──
  if (donorProviderMaxID && donorProviderMaxID.trim() !== "") {
    try {
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
        { donorId: donorProviderMaxID, recipientPetName: recipientData.petName },
        "Sent donor selected notification to donor",
      );
    } catch (err) {
      pinologger.error(
        { error: err, donorId: donorProviderMaxID },
        "Failed to send donor selected notification to donor",
      );
    }
  } else {
    pinologger.warn("Donor providerMaxId is empty, skipping donor notification");
  }

  // ── Реципиенту: контакт донора + уведомление ──
  if (recipientProviderMaxID && recipientProviderMaxID.trim() !== "") {
    try {
      await sendMessageToUser(recipientProviderMaxID, "", {
        attachments: [
          {
            type: "contact",
            payload: {
              name: donorData.userName,
              contact_id: Number(donorData.providerMaxId),
              vcf_phone: donorData.phone,
              vcf_info: generateVCF(donorData.userName, donorData.phone),
            },
          },
        ],
      });

      const donorBloodGroup =
        donorData.bloodGroup === "UNKNOWN"
          ? "не определена"
          : donorData.bloodGroup;

      const message = `Вы выбрали донора — ${donorData.petName} (группа ${donorBloodGroup}).\n\nКонтакты хозяина донора направлены. Свяжитесь с ним для обсуждения донации.`;

      await sendMessageToUser(recipientProviderMaxID, message, {
        attachments: [getAppOpenKeyboard()],
      });

      pinologger.info(
        { recipientId: recipientProviderMaxID, donorPetName: donorData.petName },
        "Sent donor selected notification to recipient",
      );
    } catch (err) {
      pinologger.error(
        { error: err, recipientId: recipientProviderMaxID },
        "Failed to send donor selected notification to recipient",
      );
    }
  } else {
    pinologger.warn("Recipient providerMaxId is empty, skipping recipient notification");
  }
};