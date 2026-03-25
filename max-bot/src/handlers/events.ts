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
              name: DonorData.Name,
              contact_id: Number(DonorData.ProviderMaxID),
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
