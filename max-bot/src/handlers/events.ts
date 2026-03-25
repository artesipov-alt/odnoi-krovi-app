import { bot, pinologger } from "../instances";

interface DonorData {
  ProviderMaxID: string;
  ProviderTelegram: string;
  Name: string;
  Phone: string;
}

interface RecipientData {
  ProviderMaxID: string;
  ProviderTelegram: string;
  Name: string;
  Phone: string;
}

interface ApplyDonorEvent {
  DonorData: DonorData;
  RecipientData: RecipientData;
}

interface RecipientApplyEvent {
  DonorName: string;
  DonorBloodGroup: string;
  RecipientProviderMaxID: string;
  CreatedAt: string;
}

const RECIPIENT_MESSAGE = `\nПожалуйста, перейдите в чат с донором по указанным ниже контактам. Будьте вежливы и доброжелательны в общении. Помните, что ваша доброта и уважение помогут сделать процесс максимально комфортным для обеих сторон.\n`;
const DONOR_MESSAGE =
  "Реципиент принял Ваше предложение. В ближайшее время с вами свяжутся.";
const RECIPIENT_NOTIFICATION = `Новый донор ${"{donorName}"} с группой крови ${"{bloodGroup}"} откликнулся на вашу просьбу о крови.`;

export const handleDonorApply = async (event: ApplyDonorEvent) => {
  const { DonorData, RecipientData } = event;

  try {
    await bot.api.sendMessageToUser(
      Number(RecipientData.ProviderMaxID),
      RECIPIENT_MESSAGE,
      {
        attachments: [
          {
            type: "contact",
            payload: {
              name: DonorData.Name,
              contact_id: Number(DonorData.ProviderMaxID),
              vcf_phone: DonorData.Phone,
              vcf_info: `BEGIN:VCARD\r\nVERSION:3.0\r\nFN:${DonorData.Name}\r\nTEL:${DonorData.Phone}\r\nEND:VCARD`,
            },
          },
        ],
      },
    );

    await bot.api.sendMessageToUser(
      Number(DonorData.ProviderMaxID),
      DONOR_MESSAGE,
      {
        attachments: [
          {
            type: "contact",
            payload: {
              name: RecipientData.Name,
              contact_id: Number(RecipientData.ProviderMaxID),
              vcf_phone: RecipientData.Phone,
              vcf_info: `BEGIN:VCARD\r\nVERSION:3.0\r\nFN:${RecipientData.Name}\r\nTEL:${RecipientData.Phone}\r\nEND:VCARD`,
            },
          },
        ],
      },
    );

    pinologger.info(
      {
        recipientId: RecipientData.ProviderMaxID,
        donorId: DonorData.ProviderMaxID,
      },
      "Sent donor apply notification and recipient contact",
    );
  } catch (err) {
    pinologger.error({ error: err }, "Failed to send message");
  }
};

export const handleRecipientApply = async (event: RecipientApplyEvent) => {
  const { DonorName, DonorBloodGroup, RecipientProviderMaxID } = event;

  try {
    const message = RECIPIENT_NOTIFICATION.replace(
      "{donorName}",
      DonorName,
    ).replace("{bloodGroup}", DonorBloodGroup);

    await bot.api.sendMessageToUser(Number(RecipientProviderMaxID), message);

    pinologger.info(
      {
        recipientId: RecipientProviderMaxID,
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
