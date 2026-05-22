import { pinologger } from "../../instances";
import { sendTelegramMessage } from "../../telegram";

// Отклик донора на рецепиента.
interface RecipientApplyEvent {
  DonorName: string;
  DonorBloodGroup: string;
  RecipientProviderTelegramID: string;
  RecipientPetName: string;
  RecipientPetSearchingBloodGroup: string[];
  RecipientPetNeededVolume: number;
  CreatedAt: string;
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
    DonorName,
    DonorBloodGroup,
    RecipientProviderTelegramID,
    RecipientPetName,
    RecipientPetSearchingBloodGroup,
    RecipientPetNeededVolume,
  } = event;

  if (
    !RecipientProviderTelegramID ||
    RecipientProviderTelegramID.trim() === ""
  ) {
    pinologger.warn(
      { donorName: DonorName },
      "RecipientProviderTelegramID is empty, skipping notification",
    );
    return;
  }

  try {
    const message = generateMessage({
      recipientPetName: RecipientPetName,
      recipientPetNeededVolume: RecipientPetNeededVolume,
      recipientPetSearchingBloodGroup: RecipientPetSearchingBloodGroup,
      donorName: DonorName,
      donorBloodGroup: DonorBloodGroup,
    });

    await sendTelegramMessage(RecipientProviderTelegramID, message);

    pinologger.info(
      {
        recipientId: RecipientProviderTelegramID,
        donorName: DonorName,
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
