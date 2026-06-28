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
import { analyticHandler } from "./src/handlers/commands";

import { bot, pinologger, redis } from "./src/instances";
import { logger } from "./src/middleware/logger";
import { run } from "@grammyjs/runner";
import { limitter } from "./src/middleware/ratelimitter";
import { errorHandler } from "./src/handlers/errors";

async function main() {
  //Плагины бота
  bot.use(logger, limitter);

  // Установка команд бота
  await bot.api.setMyCommands([
    { command: "start", description: "Запустить бота" },
    // { command: "profile", description: "Профиль пользователя" },
    // { command: "help", description: "Помощь" },
  ]);

  bot.api.config.use();

  //Команды бота
  bot.command("start", startHandler);
  bot.command("stats", analyticHandler);
  // bot.command("help", helpHandler);
  // bot.command("profile", profileHandler);
  // bot.command("err", errCommandTest);
  // bot.command("api", apiTestHandler);

  //Колбэки (нажатия на кнопки)
  // bot.callbackQuery("profile", profileHandler);
  // bot.callbackQuery("help", helpHandler);
  // bot.callbackQuery("back", startHandler);

  const { first_name, last_name, id } = await bot.api.getMe();

  pinologger.info(
    `Бот ${first_name || id}${last_name ? ` ${last_name}` : ""} запущен`,
  );

  const env = Bun.env.ENV || "production";
  const channelPrefix =
    env === "development" || env === "dev" ? "dev:" : "prod:";
  const channel = (name: string) => `${channelPrefix}${name}`;

  // Redis event handlers
  const eventHandlers: Record<string, (event: any) => Promise<void>> = {
    [channel("donor_response_apply")]: handleDonorApply,
    [channel("recipient_response_apply")]: handleRecipientApply,
    [channel("donor_cancel")]: handleDonorCancel,
    [channel("donor_reject")]: handleDonorReject,
    [channel("donor_not_confirmed")]: handleDonorNotConfirmed,
    [channel("donor_completed")]: handleDonorCompleted,
    [channel("blood_request_created")]: handleBloodRequestCreated,
    [channel("donation_confirmed")]: handleDonationConfirmed,
    [channel("user_contact")]: handleUserContact,
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

  // Redis subscriptions for events
  Object.keys(eventHandlers).forEach(subscribeToChannel);

  redis.on("message", async (channel, message) => {
    const handler = eventHandlers[channel];
    if (handler) {
      try {
        await handler(JSON.parse(message));
      } catch (err) {
        pinologger.error(
          { error: err, channel },
          `Failed to handle event for ${channel}`,
        );
      }
    } else {
      pinologger.warn({ channel }, "Unknown Redis channel");
    }
  });

  //Бот работает через раннер для паралельности задач
  const runner = run(bot);

  //Обработка ошибок
  bot.catch(errorHandler);

  // Graceful shutdown
  const shutdown = async (signal: string) => {
    pinologger.info({ signal }, "Shutting down...");
    redis.disconnect();
    await runner.stop();
    process.exit(0);
  };

  process.on("SIGTERM", () => shutdown("SIGTERM"));
  process.on("SIGINT", () => shutdown("SIGINT"));
}

main();
