import { bot, pinologger } from "../../instances";

// Уведомление донора о новых реципиентах.
interface NewRecipientsEvent {
  DonorProviderMaxIDs: string[];
}

const MESSAGE = `Вы можете спасти жизнь!

В приложении появились питомцы, которым нужна Ваша помощь.
Откликнитесь на запрос в приложении и договоритесь о донации.`;

export const handleNewRecipients = async (event: NewRecipientsEvent) => {
  const { DonorProviderMaxIDs } = event;

  if (!DonorProviderMaxIDs || DonorProviderMaxIDs.length === 0) {
    pinologger.warn(
      {},
      "DonorProviderMaxIDs array is empty, skipping notification",
    );
    return;
  }

  for (const id of DonorProviderMaxIDs) {
    if (!id || id.trim() === "") {
      pinologger.warn({ donorId: id }, "Invalid donor ID, skipping");
      continue;
    }

    try {
      await bot.api.sendMessageToUser(Number(id), MESSAGE);

      pinologger.info(
        {
          donorId: id,
        },
        "Sent new recipients notification",
      );
    } catch (err) {
      pinologger.error(
        { error: err, donorId: id },
        "Failed to send new recipients notification",
      );
    }
  }
};
