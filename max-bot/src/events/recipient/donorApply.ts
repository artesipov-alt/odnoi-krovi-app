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
  DonorData: DonorData;
  RecipientData: RecipientData;
}

interface DonorData {
  ProviderMaxID: string;
  ProviderTelegram: string;
  UserName: string;
  PetName: string;
  Phone: string;
  BloodGroup: string;
}

interface RecipientData {
  ProviderMaxID: string;
  ProviderTelegram: string;
  UserName: string;
  PetName: string;
  Phone: string;
  BloodGroup: string;
  Volume: number;
}

export const handleDonorApply = async (event: ApplyDonorEvent) => {
  const { DonorData, RecipientData } = event;

  const donorProviderMaxID = DonorData.ProviderMaxID;
  const recipientProviderMaxID = RecipientData.ProviderMaxID;

  if (!recipientProviderMaxID || recipientProviderMaxID.trim() === "") {
    pinologger.warn(
      { donorId: donorProviderMaxID },
      "RecipientProviderMaxID is empty, skipping notification",
    );
    return;
  }

  if (!donorProviderMaxID || donorProviderMaxID.trim() === "") {
    pinologger.warn(
      { recipientId: recipientProviderMaxID },
      "DonorProviderMaxID is empty, skipping notification",
    );
    return;
  }

  // Отправляем уведомление реципиенту
  try {
    const recipientMessage = generateRecipientMessage({
      donorName: DonorData.PetName,
      donorBloodGroup: DonorData.BloodGroup,
    });

    await sendMessageToUser(recipientProviderMaxID, recipientMessage, {
      attachments: [getAppOpenKeyboard()],
    });

    await sendMessageToUser(recipientProviderMaxID, "", {
      attachments: [
        {
          type: "contact",
          payload: {
            name: DonorData.UserName,
            contact_id: Number(donorProviderMaxID),
            vcf_phone: DonorData.Phone,
            vcf_info: generateVCF(DonorData.UserName, DonorData.Phone),
          },
        },
      ],
    });

    pinologger.info(
      {
        recipientId: recipientProviderMaxID,
        donorName: DonorData.PetName,
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
    const donorMessage = generateDonorMessage({
      recipientName: RecipientData.PetName,
      recipientBloodGroup: RecipientData.BloodGroup,
      recipientVolume: RecipientData.Volume,
    });

    await sendMessageToUser(donorProviderMaxID, donorMessage, {
      attachments: [getAppOpenKeyboard()],
    });

    await sendMessageToUser(donorProviderMaxID, "", {
      attachments: [
        {
          type: "contact",
          payload: {
            name: RecipientData.UserName,
            contact_id: Number(recipientProviderMaxID),
            vcf_phone: RecipientData.Phone,
            vcf_info: generateVCF(RecipientData.UserName, RecipientData.Phone),
          },
        },
      ],
    });

    pinologger.info(
      {
        donorId: donorProviderMaxID,
        recipientPetName: RecipientData.PetName,
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
