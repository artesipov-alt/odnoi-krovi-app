import { pinologger } from "../../instances";
import { sendMessageToUser } from "../../max";
import { getAppOpenKeyboard } from "../../keyboards";

import { generateVCF } from "../recipient/helpers";

interface DonorNotConfirmedEvent {
  DonorPetName: string;
  DonorBloodGroup: string;
  DonorProviderMaxID: string;
  RecipientPetName: string;
  RecipientBloodGroup: string;
  RecipientUserData: {
    Name: string;
    ProviderMaxID: string;
    ProviderTelegram: string;
    Phone: string;
  };
  DonorUserData: {
    Name: string;
    ProviderMaxID: string;
    ProviderTelegram: string;
    Phone: string;
  };
  CreatedAt: string;
}

export const handleDonorNotConfirmed = async (
  event: DonorNotConfirmedEvent,
) => {
  const {
    DonorPetName,
    DonorBloodGroup,
    DonorProviderMaxID,
    RecipientPetName,
    RecipientBloodGroup,
    RecipientUserData,
    DonorUserData,
  } = event;

  // Notify donor
  if (DonorProviderMaxID && DonorProviderMaxID.trim() !== "") {
    try {
      const donorMessage = `Хозяин реципиента (${RecipientPetName}, группа ${RecipientBloodGroup}) не подтвердил донацию. Можете связаться с ним для уточнения ситуации.`;

      await sendMessageToUser(DonorProviderMaxID, donorMessage, {
        attachments: [getAppOpenKeyboard()],
      });

      await sendMessageToUser(DonorProviderMaxID, "", {
        attachments: [
          {
            type: "contact",
            payload: {
              name: RecipientUserData.Name,
              contact_id: Number(RecipientUserData.ProviderMaxID),
              vcf_phone: RecipientUserData.Phone,
              vcf_info: generateVCF(
                RecipientUserData.Name,
                RecipientUserData.Phone,
              ),
            },
          },
        ],
      });

      pinologger.info(
        {
          donorId: DonorProviderMaxID,
          recipientPetName: RecipientPetName,
        },
        "Sent donor not confirmed notification to donor",
      );
    } catch (err) {
      pinologger.error(
        { error: err, donorId: DonorProviderMaxID },
        "Failed to send donor not confirmed notification to donor",
      );
    }
  }

  // Notify recipient
  if (
    RecipientUserData.ProviderMaxID &&
    RecipientUserData.ProviderMaxID.trim() !== ""
  ) {
    try {
      const recipientMessage = `Вы не подтвердили донацию (${DonorPetName}, группа ${DonorBloodGroup}). Можете связаться с хозяином донора для уточнения ситуации.`;

      await sendMessageToUser(
        RecipientUserData.ProviderMaxID,
        recipientMessage,
        {
          attachments: [getAppOpenKeyboard()],
        },
      );

      await sendMessageToUser(
        RecipientUserData.ProviderMaxID,
        "",
        {
          attachments: [
            {
              type: "contact",
              payload: {
                name: DonorUserData.Name,
                contact_id: Number(DonorProviderMaxID),
                vcf_phone: DonorUserData.Phone,
                vcf_info: generateVCF(DonorUserData.Name, DonorUserData.Phone),
              },
            },
          ],
        },
      );

      pinologger.info(
        {
          recipientId: RecipientUserData.ProviderMaxID,
          donorPetName: DonorPetName,
        },
        "Sent donor not confirmed notification to recipient",
      );
    } catch (err) {
      pinologger.error(
        { error: err, recipientId: RecipientUserData.ProviderMaxID },
        "Failed to send donor not confirmed notification to recipient",
      );
    }
  }
};
