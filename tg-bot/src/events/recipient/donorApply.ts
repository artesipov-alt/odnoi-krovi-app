import { pinologger } from "../../instances";
import { sendTelegramMessage, sendTelegramContact } from "../../telegram";
import { createOpenAppKeyboard } from "../../telegramButtons";

// Отклик реципиента на донора (Принятие заявки)
interface ApplyDonorEvent {
  donorData: DonorData;
  recipientData: RecipientData;
}

interface DonorData {
  providerTelegram: string;
  userName: string;
  petName: string;
  phone: string;
  bloodGroup: string;
}

interface RecipientData {
  providerTelegram: string;
  userName: string;
  petName: string;
  phone: string;
  bloodGroup: string;
  volume: number;
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
  const { donorData, recipientData } = event;

  const donorProviderTelegram = donorData.providerTelegram;
  const recipientProviderTelegram = recipientData.providerTelegram;

  if (!recipientProviderTelegram || recipientProviderTelegram.trim() === "") {
    pinologger.warn(
      { donorId: donorProviderTelegram },
      "Recipient providerTelegram is empty, skipping notification",
    );
    return;
  }

  if (!donorProviderTelegram || donorProviderTelegram.trim() === "") {
    pinologger.warn(
      { recipientId: recipientProviderTelegram },
      "Donor providerTelegram is empty, skipping notification",
    );
    return;
  }

  // Отправляем уведомление реципиенту
  try {
    // Сначала отправляем контакт донора
    if (donorData.providerTelegram && donorData.phone) {
      const nameParts = donorData.userName.split(" ");
      await sendTelegramContact(
        recipientProviderTelegram,
        donorData.phone,
        nameParts[0] || donorData.userName,
        { last_name: nameParts.slice(1).join(" ") || undefined },
      );
    }

    const recipientMessage = generateRecipientMessage({
      donorName: donorData.petName,
      donorBloodGroup: donorData.bloodGroup,
    });

    await sendTelegramMessage(recipientProviderTelegram, recipientMessage, {
      reply_markup: createOpenAppKeyboard(),
    });

    pinologger.info(
      {
        recipientId: recipientProviderTelegram,
        donorName: donorData.petName,
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
    // Сначала отправляем контакт реципиента
    if (recipientData.providerTelegram && recipientData.phone) {
      const nameParts = recipientData.userName.split(" ");
      await sendTelegramContact(
        donorProviderTelegram,
        recipientData.phone,
        nameParts[0] || recipientData.userName,
        { last_name: nameParts.slice(1).join(" ") || undefined },
      );
    }

    const donorMessage = generateDonorMessage({
      recipientName: recipientData.petName,
      recipientBloodGroup: recipientData.bloodGroup,
      recipientVolume: recipientData.volume,
    });

    await sendTelegramMessage(donorProviderTelegram, donorMessage, {
      reply_markup: createOpenAppKeyboard(),
    });

    pinologger.info(
      {
        donorId: donorProviderTelegram,
        recipientPetName: recipientData.petName,
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
