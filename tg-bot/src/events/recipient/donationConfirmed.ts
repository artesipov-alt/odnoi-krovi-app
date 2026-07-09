import { pinologger } from "../../instances";
import { sendTelegramMessage } from "../../telegram";
import { createOpenAppKeyboard } from "../../telegramButtons";

// Уведомление о подтвержденной донации (от реципиента донору).
interface DonationConfirmedEvent {
  donorData: {
    userName: string;
    petName: string;
    providerTelegram: string;
    phone: string;
    bloodGroup: string;
  };
  recipientData: {
    petName: string;
    bloodGroup: string;
  };
  volume: number;
}

const generateDonationMessage = (params: {
  volume: number;
  recipientPetName: string;
  recipientBloodGroup: string;
}): string => {
  const recipientBloodGroup =
    params.recipientBloodGroup === "UNKNOWN"
      ? "не определена"
      : params.recipientBloodGroup;

  return `Донация подтверждена (реципиент ${params.recipientPetName}, группа ${recipientBloodGroup}). Спасибо за Вашу помощь! Вам начислены бонусы – посмотрите их на Портале.`;
};

export const handleDonationConfirmed = async (
  event: DonationConfirmedEvent,
) => {
  const { donorData, recipientData, volume } = event;

  const targetId = donorData.providerTelegram;

  if (!targetId || targetId.trim() === "") {
    pinologger.warn(
      { donorUserName: donorData.userName },
      "Donor providerTelegram is empty, skipping notification",
    );
    return;
  }

  try {
    const message = generateDonationMessage({
      volume,
      recipientPetName: recipientData.petName,
      recipientBloodGroup: recipientData.bloodGroup,
    });

    await sendTelegramMessage(targetId, message, {
      reply_markup: createOpenAppKeyboard(),
    });

    pinologger.info(
      {
        targetId,
        volume,
      },
      "Sent donation confirmed notification to donor",
    );
  } catch (err) {
    pinologger.error(
      { error: err, targetId },
      "Failed to send donation confirmed notification",
    );
  }
};
