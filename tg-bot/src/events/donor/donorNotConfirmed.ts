import { pinologger } from "../../instances";
import { sendTelegramMessage, sendTelegramContact } from "../../telegram";

interface DonorNotConfirmedEvent {
  DonorPetName: string;
  DonorBloodGroup: string;
  RecipientPetName: string;
  RecipientBloodGroup: string;
  RecipientUserData: {
    Name: string;
    ProviderTelegram: string;
    Phone: string;
  };
  DonorUserData: {
    Name: string;
    ProviderTelegram: string;
    Phone: string;
  };
  CreatedAt: string;
}

export const handleDonorNotConfirmed = async (
  event: DonorNotConfirmedEvent,
) => {
  const {
    DonorPetName,
    DonorBloodGroup,
    RecipientPetName,
    RecipientBloodGroup,
    RecipientUserData,
    DonorUserData,
  } = event;

  // Notify donor
  if (
    DonorUserData.ProviderTelegram &&
    DonorUserData.ProviderTelegram.trim() !== ""
  ) {
    try {
      const donorMessage = `Хозяин реципиента (${RecipientPetName}, группа ${RecipientBloodGroup}) не подтвердил донацию. Можете связаться с ним для уточнения ситуации.`;

      await sendTelegramMessage(DonorUserData.ProviderTelegram, donorMessage);

      // Отправляем контакт реципиента
      if (RecipientUserData.ProviderTelegram && RecipientUserData.Phone) {
        const nameParts = RecipientUserData.Name.split(" ");
        await sendTelegramContact(
          DonorUserData.ProviderTelegram,
          RecipientUserData.Phone,
          nameParts[0] || RecipientUserData.Name,
          { last_name: nameParts.slice(1).join(" ") || undefined },
        );
      }

      pinologger.info(
        {
          donorId: DonorUserData.ProviderTelegram,
          recipientPetName: RecipientPetName,
        },
        "Sent donor not confirmed notification to donor",
      );
    } catch (err) {
      pinologger.error(
        { error: err, donorId: DonorUserData.ProviderTelegram },
        "Failed to send donor not confirmed notification to donor",
      );
    }
  }

  // Notify recipient
  if (
    RecipientUserData.ProviderTelegram &&
    RecipientUserData.ProviderTelegram.trim() !== ""
  ) {
    try {
      const recipientMessage = `Вы не подтвердили донацию (${DonorPetName}, группа ${DonorBloodGroup}). Можете связаться с хозяином донора для уточнения ситуации.`;

      await sendTelegramMessage(
        RecipientUserData.ProviderTelegram,
        recipientMessage,
      );

      // Отправляем контакт донора
      if (DonorUserData.ProviderTelegram && DonorUserData.Phone) {
        const nameParts = DonorUserData.Name.split(" ");
        await sendTelegramContact(
          RecipientUserData.ProviderTelegram,
          DonorUserData.Phone,
          nameParts[0] || DonorUserData.Name,
          { last_name: nameParts.slice(1).join(" ") || undefined },
        );
      }

      pinologger.info(
        {
          recipientId: RecipientUserData.ProviderTelegram,
          donorPetName: DonorPetName,
        },
        "Sent donor not confirmed notification to recipient",
      );
    } catch (err) {
      pinologger.error(
        { error: err, recipientId: RecipientUserData.ProviderTelegram },
        "Failed to send donor not confirmed notification to recipient",
      );
    }
  }
};
