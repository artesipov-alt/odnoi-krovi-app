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
import { handleDonorReject } from "./src/events/donor/donorReject";
import { handleDonorNotConfirmed } from "./src/events/donor/donorNotConfirmed";
import { handleDonorCompleted } from "./src/events/donor/donorCompleted";
import { handleBloodRequestCreated } from "./src/events/recipient/handleBloodRequestCreated";
import { handleDonationConfirmed } from "./src/events/recipient/donationConfirmed";
import { handleUserContact } from "./src/events/user/handleUserContact";

import { bot, pinologger, redis } from "./src/instances";
import { logger } from "./src/middleware/logger";
import { errorHandler } from "./src/handlers/errors";
import { startServer } from "./src/server";

// Event handlers map для Redis
const eventHandlers: Record<string, (event: any) => Promise<void>> = {
  donor_response_apply: handleDonorApply,
  recipient_response_apply: handleRecipientApply,
  donor_cancel: handleDonorCancel,
  donor_reject: handleDonorReject,
  donor_not_confirmed: handleDonorNotConfirmed,
  donor_completed: handleDonorCompleted,
  blood_request_created: handleBloodRequestCreated,
  donation_confirmed: handleDonationConfirmed,
  user_contact: handleUserContact,
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
  // Плагины бота
  bot.use(logger);

  // Установка команд бота
  await bot.api.setMyCommands([
    { name: "start", description: "Запустить бота" },
  ]);

  // Команды и действия бота
  bot.on(`bot_started`, startHandler);
  bot.command("start", startHandler);
  bot.command("err", errCommandTest);
  bot.command("api", apiTestHandler);
  bot.action("back", startHandler);

  const { name, username, user_id } = await bot.api.getMyInfo();
  pinologger.info(`Бот ${name || username} ${user_id} инициализирован`);

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

  // Запуск HTTP сервера для вебхуков
  const server = startServer(bot);

  // Обработка ошибок
  bot.catch(errorHandler);

  // Graceful shutdown
  const shutdown = () => {
    pinologger.info("Shutting down...");
    server.stop();
    redis.quit();
    process.exit(0);
  };

  process.on("SIGINT", shutdown);
  process.on("SIGTERM", shutdown);
}

main();
