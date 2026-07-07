import { pinologger } from "../../instances";
import { sendMessageToUser } from "../../max";
import { getAppOpenKeyboard } from "../../keyboards";

import { generateVCF } from "../recipient/helpers";

interface DonorNotConfirmedEvent {
  donorPetName: string;
  donorBloodGroup: string;
  donorProviderMaxId: string;
  recipientPetName: string;
  recipientBloodGroup: string;
  recipientUserData: {
    name: string;
    providerMaxId: string;
    providerTelegram: string;
    phone: string;
  };
  donorUserData: {
    name: string;
    providerMaxId: string;
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
    donorProviderMaxId,
    recipientPetName,
    recipientBloodGroup,
    recipientUserData,
    donorUserData,
  } = event;

  // Notify donor
  if (donorProviderMaxId && donorProviderMaxId.trim() !== "") {
    try {
      // Сначала отправляем контакт реципиента
      await sendMessageToUser(donorProviderMaxId, "", {
        attachments: [
          {
            type: "contact",
            payload: {
              name: recipientUserData.name,
              contact_id: Number(recipientUserData.providerMaxId),
              vcf_phone: recipientUserData.phone,
              vcf_info: generateVCF(
                recipientUserData.name,
                recipientUserData.phone,
              ),
            },
          },
        ],
      });

      const donorMessage = `Хозяин реципиента (${recipientPetName}, группа ${recipientBloodGroup}) не подтвердил донацию. Можете связаться с ним для уточнения ситуации.`;

      await sendMessageToUser(donorProviderMaxId, donorMessage, {
        attachments: [getAppOpenKeyboard()],
      });

      pinologger.info(
        {
          donorId: donorProviderMaxId,
          recipientPetName,
        },
        "Sent donor not confirmed notification to donor",
      );
    } catch (err) {
      pinologger.error(
        { error: err, donorId: donorProviderMaxId },
        "Failed to send donor not confirmed notification to donor",
      );
    }
  }

  // Notify recipient
  if (
    recipientUserData.providerMaxId &&
    recipientUserData.providerMaxId.trim() !== ""
  ) {
    try {
      // Сначала отправляем контакт донора
      await sendMessageToUser(recipientUserData.providerMaxId, "", {
        attachments: [
          {
            type: "contact",
            payload: {
              name: donorUserData.name,
              contact_id: Number(donorProviderMaxId),
              vcf_phone: donorUserData.phone,
              vcf_info: generateVCF(donorUserData.name, donorUserData.phone),
            },
          },
        ],
      });

      const recipientMessage = `Вы не подтвердили донацию (${donorPetName}, группа ${donorBloodGroup}). Можете связаться с хозяином донора для уточнения ситуации.`;

      await sendMessageToUser(
        recipientUserData.providerMaxId,
        recipientMessage,
        {
          attachments: [getAppOpenKeyboard()],
        },
      );

      pinologger.info(
        {
          recipientId: recipientUserData.providerMaxId,
          donorPetName,
        },
        "Sent donor not confirmed notification to recipient",
      );
    } catch (err) {
      pinologger.error(
        { error: err, recipientId: recipientUserData.providerMaxId },
        "Failed to send donor not confirmed notification to recipient",
      );
    }
  }
};
