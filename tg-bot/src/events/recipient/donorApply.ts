import { pinologger } from "../../instances";
import { sendTelegramMessage, sendTelegramContact } from "../../telegram";
import { createOpenAppKeyboard } from "../../telegramButtons";

// Отклик реципиента на донора.(Принятие заявки)
interface ApplyDonorEvent {
  DonorData: DonorData;
  RecipientData: RecipientData;
}

interface DonorData {
  ProviderTelegram: string;
  UserName: string;
  PetName: string;
  Phone: string;
  BloodGroup: string;
}

interface RecipientData {
  ProviderTelegram: string;
  UserName: string;
  PetName: string;
  Phone: string;
  BloodGroup: string;
  Volume: number;
}

const generateRecipientMessage = (params: {
  donorName: string;
  donorBloodGroup: string;
}): string => {
  const donorName = params.donorName?.trim() || "Анонимный донор";
  const donorBloodGroup =
    params.donorBloodGroup === "UNKNOWN"
      ? "не определена"
      : params.donorBloodGroup;

  return `\nКонтакты хозяина донора - ${donorName} (группа ${donorBloodGroup})\n\nНаправляем контакты хозяина донора - обсудите возможность донации.\nБудьте вежливы и доброжелательны в общении!\nЕсли не получится договориться о донации, можете продолжить поиск в приложении.\n`;
};

const generateDonorMessage = (params: {
  recipientName: string;
  recipientBloodGroup: string;
  recipientVolume: number;
}): string => {
  const recipientBloodGroup =
    params.recipientBloodGroup === "UNKNOWN"
      ? "не определена"
      : params.recipientBloodGroup;

  return `На ваше предложение откликнулся реципиент - ${params.recipientName} (${params.recipientVolume} мл, группа ${recipientBloodGroup})\n\nХозяин реципиента получил Ваши контакты. Дождитесь, пока с Вами свяжутся, или напишите хозяину реципиента`;
};

export const handleDonorApply = async (event: ApplyDonorEvent) => {
  const { DonorData, RecipientData } = event;

  const donorProviderTelegram = DonorData.ProviderTelegram;
  const recipientProviderTelegram = RecipientData.ProviderTelegram;

  if (!recipientProviderTelegram || recipientProviderTelegram.trim() === "") {
    pinologger.warn(
      { donorId: donorProviderTelegram },
      "Recipient ProviderTelegram is empty, skipping notification",
    );
    return;
  }

  if (!donorProviderTelegram || donorProviderTelegram.trim() === "") {
    pinologger.warn(
      { recipientId: recipientProviderTelegram },
      "Donor ProviderTelegram is empty, skipping notification",
    );
    return;
  }

  // Отправляем уведомление реципиенту
  try {
    const recipientMessage = generateRecipientMessage({
      donorName: DonorData.PetName,
      donorBloodGroup: DonorData.BloodGroup,
    });

    await sendTelegramMessage(recipientProviderTelegram, recipientMessage, {
      reply_markup: createOpenAppKeyboard(),
    });

    // Отправляем контакт донора реципиенту
    if (DonorData.ProviderTelegram && DonorData.Phone) {
      const nameParts = DonorData.UserName.split(" ");
      await sendTelegramContact(
        recipientProviderTelegram,
        DonorData.Phone,
        nameParts[0] || DonorData.UserName,
        { last_name: nameParts.slice(1).join(" ") || undefined },
      );
    }

    pinologger.info(
      {
        recipientId: recipientProviderTelegram,
        donorName: DonorData.PetName,
      },
      "Sent donor apply notification to recipient",
    );
  } catch (err) {
    pinologger.error(
      { error: err, recipientId: recipientProviderTelegram },
      "Failed to send donor apply notification to recipient",
    );
  }

  // Отправляем уведомление донору (независимо от отправки реципиенту)
  try {
    const donorMessage = generateDonorMessage({
      recipientName: RecipientData.PetName,
      recipientBloodGroup: RecipientData.BloodGroup,
      recipientVolume: RecipientData.Volume,
    });

    await sendTelegramMessage(donorProviderTelegram, donorMessage, {
      reply_markup: createOpenAppKeyboard(),
    });

    // Отправляем контакт реципиента донору
    if (RecipientData.ProviderTelegram && RecipientData.Phone) {
      const nameParts = RecipientData.UserName.split(" ");
      await sendTelegramContact(
        donorProviderTelegram,
        RecipientData.Phone,
        nameParts[0] || RecipientData.UserName,
        { last_name: nameParts.slice(1).join(" ") || undefined },
      );
    }

    pinologger.info(
      {
        donorId: donorProviderTelegram,
        recipientPetName: RecipientData.PetName,
      },
      "Sent donor apply notification to donor",
    );
  } catch (err) {
    pinologger.error(
      { error: err, donorId: donorProviderTelegram },
      "Failed to send donor apply notification to donor",
    );
  }
};
