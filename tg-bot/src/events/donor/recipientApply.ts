import { pinologger } from "../../instances";
import { sendTelegramMessage } from "../../telegram";
import { createOpenAppKeyboard } from "../../telegramButtons";

// Отклик донора на реципиента
interface RecipientApplyEvent {
  donorName: string;
  donorBloodGroup: string;
  recipientProviderTelegramId: string;
  recipientPetName: string;
  recipientPetSearchingBloodGroup: string[];
  recipientPetNeededVolume: number;
  createdAt: string;
}

const generateMessage = (params: {
  recipientPetName: string;
  recipientPetNeededVolume: number;
  recipientPetSearchingBloodGroup: string[];
  donorName: string;
  donorBloodGroup: string;
}): string => {
  const donorName = params.donorName?.trim() || "Анонимный донор";
  const donorBloodGroup =
    params.donorBloodGroup === "UNKNOWN"
      ? "не определена"
      : params.donorBloodGroup;

  return `На Ваш поиск (${params.recipientPetName}, ${params.recipientPetNeededVolume} мл, группа ${params.recipientPetSearchingBloodGroup.join(", ")}) откликнулся донор ${donorName} (группа ${donorBloodGroup}).\nОзнакомьтесь с информацией о доноре в приложении.`;
};

export const handleRecipientApply = async (event: RecipientApplyEvent) => {
  const {
    donorName,
    donorBloodGroup,
    recipientProviderTelegramId,
    recipientPetName,
    recipientPetSearchingBloodGroup,
    recipientPetNeededVolume,
  } = event;

  if (
    !recipientProviderTelegramId ||
    recipientProviderTelegramId.trim() === ""
  ) {
    pinologger.warn(
      { donorName },
      "recipientProviderTelegramId is empty, skipping notification",
    );
    return;
  }

  try {
    const message = generateMessage({
      recipientPetName,
      recipientPetNeededVolume,
      recipientPetSearchingBloodGroup,
      donorName,
      donorBloodGroup,
    });

    await sendTelegramMessage(recipientProviderTelegramId, message, {
      reply_markup: createOpenAppKeyboard(),
    });

    pinologger.info(
      {
        recipientId: recipientProviderTelegramId,
        donorName,
      },
      "Sent recipient apply notification",
    );
  } catch (err) {
    pinologger.error(
      { error: err },
      "Failed to send recipient apply notification",
    );
  }
};
