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

export const handleDonorApply = async (event: ApplyDonorEvent) => {
  const { DonorData, RecipientData } = event;

  try {
    await bot.api.sendMessageToChat(
      Number(RecipientData.ProviderMaxID),
      "Контакт донора:",
      {
        attachments: [
          {
            type: "contact",
            payload: {
              vcf_info: `BEGIN:VCARD\r\nVERSION:3.0\r\nFN:${DonorData.Name}\r\nTEL:${DonorData.Phone}\r\nEND:VCARD`,
              //@ts-ignore
              max_info: {
                user_id: Number(DonorData.ProviderMaxID),
                first_name: DonorData.Name,
                is_bot: false,
                last_activity_time: Date.now(),
                name: DonorData.Name,
              },
            },
          },
        ],
      },
    );
    pinologger.info(
      { recipientId: RecipientData.ProviderMaxID },
      "Sent donor apply notification",
    );
  } catch (err) {
    pinologger.error({ error: err }, "Failed to send message");
  }
};
