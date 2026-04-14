import {
  helpHandler,
  startHandler,
  profileHandler,
  errCommandTest,
  apiTestHandler,
} from "./src/handlers/commands";
import { handleDonorApply } from "./src/events/recipient/donorApply";
import { handleRecipientApply } from "./src/events/donor/recipientApply";
import { handleDonorCancel } from "./src/events/donor/donorCancel";
import { handleNewRecipients } from "./src/events/recipient/newRecipients";
import { handleDonationConfirmed } from "./src/events/recipient/donationConfirmed";

import { bot, pinologger, redis } from "./src/instances";
import { logger } from "./src/middleware/logger";
import { errorHandler } from "./src/handlers/errors";

// Event handlers map
const eventHandlers: Record<string, (event: any) => Promise<void>> = {
  donor_response_apply: handleDonorApply,
  recipient_response_apply: handleRecipientApply,
  donor_cancel: handleDonorCancel,
  new_recipients: handleNewRecipients,
  donation_confirmed: handleDonationConfirmed,
};

// Helper function for subscribing to channels
const subscribeToChannel = (channel: string) => {
  redis.subscribe(channel, (err, count) => {
    if (err) {
      pinologger.error({ error: err }, `Failed to subscribe to ${channel}`);
    } else {
      pinologger.info(`Subscribed to ${channel}`);
    }
  });
};

// Redis error handling
redis.on("error", (err) => {
  pinologger.error({ error: err }, "Redis connection error");
});

redis.on("connect", () => {
  pinologger.info("Connected to Redis");
});

async function main() {
  //Плагины бота
  bot.use(logger);

  // Установка команд бота
  await bot.api.setMyCommands([
    { name: "start", description: "Запустить бота" },
    { name: "profile", description: "Профиль пользователя" },
    { name: "help", description: "Помощь" },
  ]);

  bot.on(`bot_started`, startHandler);

  //Команды бота
  bot.command("start", startHandler);
  bot.command("help", helpHandler);
  bot.command("profile", profileHandler);
  bot.command("err", errCommandTest);
  bot.command("api", apiTestHandler);

  //Колбэки (нажатия на кнопки)
  bot.action("profile", profileHandler);
  bot.action("help", helpHandler);
  bot.action("back", startHandler);

  const { name, username, user_id } = await bot.api.getMyInfo();

  pinologger.info(`Бот ${name || username} ${user_id} запущен`);

  // Redis subscriptions for events
  Object.keys(eventHandlers).forEach(subscribeToChannel);

  redis.on("message", (channel, message) => {
    const handler = eventHandlers[channel];
    if (handler) {
      try {
        handler(JSON.parse(message));
      } catch (err) {
        pinologger.error(
          { error: err },
          `Failed to parse event for ${channel}`,
        );
      }
    } else {
      pinologger.warn({ channel }, "Unknown Redis channel");
    }
  });

  pinologger.info("Initiating bot polling");
  bot.start();
  pinologger.info("Bot polling initiated");

  //Обработка ошибок
  bot.catch(errorHandler);
}

main();
