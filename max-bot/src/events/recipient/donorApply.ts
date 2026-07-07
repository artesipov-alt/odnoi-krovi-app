import { pinologger } from "../../instances";
import { sendMessageToUser } from "../../max";
import { getAppOpenKeyboard } from "../../keyboards";

import {
  generateRecipientMessage,
  generateDonorMessage,
  generateVCF,
} from "./helpers";

// Отклик реципиента на донора.(Принятие заявки)
interface ApplyDonorEvent {
  donorData: {
    providerMaxId: string;
    providerTelegram: string;
    userName: string;
    petName: string;
    phone: string;
    bloodGroup: string;
  };
  recipientData: {
    providerMaxId: string;
    providerTelegram: string;
    userName: string;
    petName: string;
    phone: string;
    bloodGroup: string;
    volume: number;
  };
}

export const handleDonorApply = async (event: ApplyDonorEvent) => {
  const { donorData, recipientData } = event;

  const donorProviderMaxID = donorData.providerMaxId;
  const recipientProviderMaxID = recipientData.providerMaxId;

  if (!recipientProviderMaxID || recipientProviderMaxID.trim() === "") {
    pinologger.warn(
      { donorId: donorProviderMaxID },
      "Recipient providerMaxId is empty, skipping notification",
    );
    return;
  }

  if (!donorProviderMaxID || donorProviderMaxID.trim() === "") {
    pinologger.warn(
      { recipientId: recipientProviderMaxID },
      "Donor providerMaxId is empty, skipping notification",
    );
    return;
  }

  // Отправляем уведомление реципиенту
  try {
    // Сначала отправляем контакт донора
    await sendMessageToUser(recipientProviderMaxID, "", {
      attachments: [
        {
          type: "contact",
          payload: {
            name: donorData.userName,
            contact_id: Number(donorProviderMaxID),
            vcf_phone: donorData.phone,
            vcf_info: generateVCF(donorData.userName, donorData.phone),
          },
        },
      ],
    });

    const recipientMessage = generateRecipientMessage({
      donorName: donorData.petName,
      donorBloodGroup: donorData.bloodGroup,
    });

    await sendMessageToUser(recipientProviderMaxID, recipientMessage, {
      attachments: [getAppOpenKeyboard()],
    });

    pinologger.info(
      {
        recipientId: recipientProviderMaxID,
        donorName: donorData.petName,
      },
      "Sent donor apply notification to recipient",
    );
  } catch (err) {
    pinologger.error(
      { error: err, recipientId: recipientProviderMaxID },
      "Failed to send donor apply notification to recipient",
    );
  }

  // Отправляем уведомление донору (независимо от отправки реципиенту)
  try {
    // Сначала отправляем контакт реципиента
    await sendMessageToUser(donorProviderMaxID, "", {
      attachments: [
        {
          type: "contact",
          payload: {
            name: recipientData.userName,
            contact_id: Number(recipientProviderMaxID),
            vcf_phone: recipientData.phone,
            vcf_info: generateVCF(recipientData.userName, recipientData.phone),
          },
        },
      ],
    });

    const donorMessage = generateDonorMessage({
      recipientName: recipientData.petName,
      recipientBloodGroup: recipientData.bloodGroup,
      recipientVolume: recipientData.volume,
    });

    await sendMessageToUser(donorProviderMaxID, donorMessage, {
      attachments: [getAppOpenKeyboard()],
    });

    pinologger.info(
      {
        donorId: donorProviderMaxID,
        recipientPetName: recipientData.petName,
      },
      "Sent donor apply notification to donor",
    );
  } catch (err) {
    pinologger.error(
      { error: err, donorId: donorProviderMaxID },
      "Failed to send donor apply notification to donor",
    );
  }
};
