import { pinologger } from "../../instances";
import { sendMessageToUser } from "../../max";
import { generateVCF } from "../recipient/helpers";
import { getAppOpenKeyboard } from "../../keyboards";

interface UserContactEvent {
  NotifyProvider: string;
  SendTo: string;
  UserData: {
    Name: string;
    ProviderMaxID: string;
    ProviderTelegram: string;
    Phone: string;
  };
  CreatedAt: string;
  Recipient: any; // Assuming it's the user object, but not used
}

export const handleUserContact = async (event: UserContactEvent) => {
  const { NotifyProvider, SendTo, UserData } = event;

  // Обрабатываем только события, предназначенные для Max Bot
  if (NotifyProvider !== "max_bot") {
    pinologger.warn(
      { notifyProvider: NotifyProvider },
      "NotifyProvider is not max_bot, skipping",
    );
    return;
  }

  if (!SendTo || SendTo.trim() === "") {
    pinologger.warn(
      { notifyProvider: NotifyProvider },
      "SendTo is empty, skipping notification",
    );
    return;
  }

  try {
    const message = `Контакт пользователя`;

    await sendMessageToUser(SendTo, message, {
      attachments: [getAppOpenKeyboard()],
    });

    await sendMessageToUser(SendTo, "", {
      attachments: [
        {
          type: "contact",
          payload: {
            name: UserData.Name,
            contact_id: Number(UserData.ProviderMaxID),
            vcf_phone: UserData.Phone,
            vcf_info: generateVCF(UserData.Name, UserData.Phone),
          },
        },
      ],
    });

    pinologger.info(
      { sendTo: SendTo, userName: UserData.Name },
      "Sent user contact notification",
    );
  } catch (err) {
    pinologger.error(
      { error: err },
      "Failed to send user contact notification",
    );
  }
};
