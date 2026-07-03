import { pinologger } from "../../instances";
import { sendMessageToUser } from "../../max";
import { getAppOpenKeyboard } from "../../keyboards";

// Отклик донора на реципиента.
interface RecipientApplyEvent {
  donorName: string;
  donorBloodGroup: string;
  recipientProviderMaxId: string;
  recipientPetName: string;
  recipientPetSearchingBloodGroup: string[];
  recipientPetNeededVolume: number;
  createdAt: string;
}

export const handleRecipientApply = async (event: RecipientApplyEvent) => {
  const {
    donorName,
    donorBloodGroup,
    recipientProviderMaxId,
    recipientPetName,
    recipientPetSearchingBloodGroup,
    recipientPetNeededVolume,
  } = event;

  if (!recipientProviderMaxId || recipientProviderMaxId.trim() === "") {
    pinologger.warn(
      { donorName },
      "recipientProviderMaxId is empty, skipping notification",
    );
    return;
  }

  try {
    const message = `На Ваш поиск (${recipientPetName}, ${recipientPetNeededVolume} мл, группа ${recipientPetSearchingBloodGroup.join(", ")}) откликнулся донор ${donorName} (группа ${donorBloodGroup}).\nОзнакомьтесь с информацией о доноре в приложении.`;

    await sendMessageToUser(recipientProviderMaxId, message, {
      attachments: [getAppOpenKeyboard()],
    });

    pinologger.info(
      {
        recipientId: recipientProviderMaxId,
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
