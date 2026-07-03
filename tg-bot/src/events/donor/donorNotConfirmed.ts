import { pinologger } from "../../instances";
import { sendTelegramMessage, sendTelegramContact } from "../../telegram";
import { createOpenAppKeyboard } from "../../telegramButtons";

interface DonorNotConfirmedEvent {
  donorPetName: string;
  donorBloodGroup: string;
  recipientPetName: string;
  recipientBloodGroup: string;
  recipientUserData: {
    name: string;
    providerTelegram: string;
    phone: string;
  };
  donorUserData: {
    name: string;
    providerTelegram: string;
    phone: string;
  };
  createdAt: string;
}

export const handleDonorNotConfirmed = async (
  event: DonorNotConfirmedEvent,
) => {
  const {
    donorPetName,
    donorBloodGroup,
    recipientPetName,
    recipientBloodGroup,
    recipientUserData,
    donorUserData,
  } = event;

  // Notify donor
  if (
    donorUserData.providerTelegram &&
    donorUserData.providerTelegram.trim() !== ""
  ) {
    try {
      const donorMessage = `Хозяин реципиента (${recipientPetName}, группа ${recipientBloodGroup}) не подтвердил донацию. Можете связаться с ним для уточнения ситуации.`;

      await sendTelegramMessage(donorUserData.providerTelegram, donorMessage, {
        reply_markup: createOpenAppKeyboard(),
      });

      // Отправляем контакт реципиента
      if (recipientUserData.providerTelegram && recipientUserData.phone) {
        const nameParts = recipientUserData.name.split(" ");
        await sendTelegramContact(
          donorUserData.providerTelegram,
          recipientUserData.phone,
          nameParts[0] || recipientUserData.name,
          { last_name: nameParts.slice(1).join(" ") || undefined },
        );
      }

      pinologger.info(
        {
          donorId: donorUserData.providerTelegram,
          recipientPetName,
        },
        "Sent donor not confirmed notification to donor",
      );
    } catch (err) {
      pinologger.error(
        { error: err, donorId: donorUserData.providerTelegram },
        "Failed to send donor not confirmed notification to donor",
      );
    }
  }

  // Notify recipient
  if (
    recipientUserData.providerTelegram &&
    recipientUserData.providerTelegram.trim() !== ""
  ) {
    try {
      const recipientMessage = `Вы не подтвердили донацию (${donorPetName}, группа ${donorBloodGroup}). Можете связаться с хозяином донора для уточнения ситуации.`;

      await sendTelegramMessage(
        recipientUserData.providerTelegram,
        recipientMessage,
        { reply_markup: createOpenAppKeyboard() },
      );

      // Отправляем контакт донора
      if (donorUserData.providerTelegram && donorUserData.phone) {
        const nameParts = donorUserData.name.split(" ");
        await sendTelegramContact(
          recipientUserData.providerTelegram,
          donorUserData.phone,
          nameParts[0] || donorUserData.name,
          { last_name: nameParts.slice(1).join(" ") || undefined },
        );
      }

      pinologger.info(
        {
          recipientId: recipientUserData.providerTelegram,
          donorPetName,
        },
        "Sent donor not confirmed notification to recipient",
      );
    } catch (err) {
      pinologger.error(
        { error: err, recipientId: recipientUserData.providerTelegram },
        "Failed to send donor not confirmed notification to recipient",
      );
    }
  }
};
