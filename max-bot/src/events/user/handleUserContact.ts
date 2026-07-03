import { pinologger } from "../../instances";
import { sendMessageToUser } from "../../max";
import { getAppOpenKeyboard } from "../../keyboards";
import { generateVCF } from "../recipient/helpers";

interface UserContactEvent {
  notifyProvider: string;
  sendTo: string;
  userData: {
    name: string;
    providerMaxId: string;
    providerTelegram: string;
    phone: string;
  };
  createdAt: string;
  recipient: any;
}

export const handleUserContact = async (event: UserContactEvent) => {
  const { notifyProvider, sendTo, userData } = event;

  // Обрабатываем только события, предназначенные для Max Bot
  if (notifyProvider !== "max_bot") {
    pinologger.warn(
      { notifyProvider },
      "notifyProvider is not max_bot, skipping",
    );
    return;
  }

  if (!sendTo || sendTo.trim() === "") {
    pinologger.warn(
      { notifyProvider },
      "sendTo is empty, skipping notification",
    );
    return;
  }

  try {
    const message = `Контакт пользователя`;

    await sendMessageToUser(sendTo, message, {
      attachments: [getAppOpenKeyboard()],
    });

    await sendMessageToUser(sendTo, "", {
      attachments: [
        {
          type: "contact",
          payload: {
            name: userData.name,
            contact_id: Number(userData.providerMaxId),
            vcf_phone: userData.phone,
            vcf_info: generateVCF(userData.name, userData.phone),
          },
        },
      ],
    });

    pinologger.info(
      { sendTo, userName: userData.name },
      "Sent user contact notification",
    );
  } catch (err) {
    pinologger.error(
      { error: err },
      "Failed to send user contact notification",
    );
  }
};
