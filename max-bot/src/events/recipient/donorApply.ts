import { bot, pinologger } from "../../instances";

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

  try {
    const recipientMessage = generateRecipientMessage({
      donorName: DonorData.PetName,
      donorBloodGroup: DonorData.BloodGroup,
    });

    const donorMessage = generateDonorMessage({
      recipientName: RecipientData.PetName,
      recipientBloodGroup: RecipientData.BloodGroup,
      recipientVolume: RecipientData.Volume,
    });

    await bot.api.sendMessageToUser(
      Number(recipientProviderMaxID),
      recipientMessage,
      {
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
      },
    );

    await bot.api.sendMessageToUser(Number(donorProviderMaxID), donorMessage, {
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
        recipientId: recipientProviderMaxID,
        donorId: donorProviderMaxID,
      },
      "Sent donor apply notification and recipient contact",
    );
  } catch (err) {
    pinologger.error({ error: err }, "Failed to send message");
  }
};
