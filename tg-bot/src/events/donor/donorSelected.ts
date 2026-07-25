import { pinologger } from "../../instances";
import { sendTelegramMessage, sendTelegramContact } from "../../telegram";
import { createOpenAppKeyboard } from "../../telegramButtons";

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

export const handleDonorSelected = async (event: DonorSelectedEvent) => {
  const { donorData, recipientData } = event;

  const donorProviderTelegram = donorData.providerTelegram;
  const recipientProviderTelegram = recipientData.providerTelegram;

  // ── Донору: контакт реципиента + уведомление ──
  if (donorProviderTelegram && donorProviderTelegram.trim() !== "") {
    try {
      if (recipientData.providerTelegram && recipientData.phone) {
        const nameParts = recipientData.userName.split(" ");
        await sendTelegramContact(
          donorProviderTelegram,
          recipientData.phone,
          nameParts[0] || recipientData.userName,
          { last_name: nameParts.slice(1).join(" ") || undefined },
        );
      }

      const typeLabel = petTypeLabel(recipientData.petType);
      const bloodGroup =
        recipientData.bloodGroup === "UNKNOWN"
          ? "не определена"
          : recipientData.bloodGroup;

      const message = `На Портале ищут донора для питомца - ${typeLabel} ${recipientData.petName} (${recipientData.volume} мл, группа ${bloodGroup})\n\nРанее Вы разрешили связаться с Вами, если нужна помощь.\nХозяин реципиента получил Ваши контакты. Дождитесь, пока с Вами свяжутся, или напишите хозяину реципиента`;

      await sendTelegramMessage(donorProviderTelegram, message, {
        reply_markup: createOpenAppKeyboard(),
      });

      pinologger.info(
        { donorId: donorProviderTelegram, recipientPetName: recipientData.petName },
        "Sent donor selected notification to donor",
      );
    } catch (err) {
      pinologger.error(
        { error: err, donorId: donorProviderTelegram },
        "Failed to send donor selected notification to donor",
      );
    }
  } else {
    pinologger.warn("Donor providerTelegram is empty, skipping donor notification");
  }

  // ── Реципиенту: контакт донора + уведомление ──
  if (recipientProviderTelegram && recipientProviderTelegram.trim() !== "") {
    try {
      if (donorData.providerTelegram && donorData.phone) {
        const nameParts = donorData.userName.split(" ");
        await sendTelegramContact(
          recipientProviderTelegram,
          donorData.phone,
          nameParts[0] || donorData.userName,
          { last_name: nameParts.slice(1).join(" ") || undefined },
        );
      }

      const donorBloodGroup =
        donorData.bloodGroup === "UNKNOWN"
          ? "не определена"
          : donorData.bloodGroup;

      const message = `Вы выбрали донора — ${donorData.petName} (группа ${donorBloodGroup}).\n\nКонтакты хозяина донора направлены. Свяжитесь с ним для обсуждения донации.`;

      await sendTelegramMessage(recipientProviderTelegram, message, {
        reply_markup: createOpenAppKeyboard(),
      });

      pinologger.info(
        { recipientId: recipientProviderTelegram, donorPetName: donorData.petName },
        "Sent donor selected notification to recipient",
      );
    } catch (err) {
      pinologger.error(
        { error: err, recipientId: recipientProviderTelegram },
        "Failed to send donor selected notification to recipient",
      );
    }
  } else {
    pinologger.warn("Recipient providerTelegram is empty, skipping recipient notification");
  }
};