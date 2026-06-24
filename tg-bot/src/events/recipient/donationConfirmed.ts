import { pinologger } from "../../instances";
import { sendTelegramMessage } from "../../telegram";
import { createOpenAppKeyboard } from "../../telegramButtons";

// Уведомление о подтвержденной донации (от реципиента донору).
interface DonationConfirmedEvent {
  DonorData: {
    UserName: string;
    PetName: string;
    ProviderTelegram: string;
    Phone: string;
    BloodGroup: string;
  };
  RecipientData: {
    PetName: string;
    BloodGroup: string;
  };
  Volume: number; // Объем донации в мл
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
  const { DonorData, RecipientData, Volume } = event;

  let targetId = DonorData.ProviderTelegram;

  if (!targetId || targetId.trim() === "") {
    pinologger.warn(
      { donorUserName: DonorData.UserName },
      "Donor ProviderTelegram is empty, skipping notification",
    );
    return;
  }

  try {
    const message = generateDonationMessage({
      volume: Volume,
      recipientPetName: RecipientData.PetName,
      recipientBloodGroup: RecipientData.BloodGroup,
    });

    await sendTelegramMessage(targetId, message, {
      reply_markup: createOpenAppKeyboard(),
    });

    pinologger.info(
      {
        targetId,
        volume: Volume,
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
