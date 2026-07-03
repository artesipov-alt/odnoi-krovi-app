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
import { handleNotification } from "./src/events/notification/handleNotification";
import { analyticHandler } from "./src/handlers/commands";
import { handleNotificationRespond } from "./src/events/notification/handleNotificationRespond";

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

  // Обработка ответов на уведомления (п.5)
  bot.callbackQuery(/^notification_(yes|no)_(.+)$/, async (ctx) => {
    if (!ctx.match || ctx.match.length < 3) {
      pinologger.warn(
        { match: ctx.match },
        "Invalid notification callback match",
      );
      await ctx.answerCallbackQuery();
      return;
    }
    const [, action, requestId] = ctx.match;
    await handleNotificationRespond(ctx, action!, requestId!);
  });

  const { first_name, last_name, id } = await bot.api.getMe();

  pinologger.info(
    `Бот ${first_name || id}${last_name ? ` ${last_name}` : ""} запущен`,
  );

  const env = Bun.env.ENV || "production";
  const channelPrefix =
    env === "development" || env === "dev" ? "dev:" : "prod:";
  const channel = (name: string) => `${channelPrefix}${name}`;

  // Event envelope — парсим тип и диспатчим
  interface EventEnvelope {
    type: string;
    payload: any;
    createdAt: string;
  }

  // Redis event handlers — маппинг по типу события
  const eventHandlers: Record<string, (event: any) => Promise<void>> = {
    [channel("events")]: async (message: string) => {
      const envelope: EventEnvelope = JSON.parse(message);
      const { type, payload } = envelope;

      pinologger.info({ type }, "Received event");

      switch (type) {
        case "blood_request_created":
          await handleBloodRequestCreated(payload);
          break;
        case "donor_response_apply":
          await handleDonorApply(payload);
          break;
        case "donation_confirmed":
          await handleDonationConfirmed(payload);
          break;
        case "recipient_response_apply":
          await handleRecipientApply(payload);
          break;
        case "donor_cancel":
          await handleDonorCancel(payload);
          break;
        case "donor_reject":
          await handleDonorReject(payload);
          break;
        case "donor_not_confirmed":
          await handleDonorNotConfirmed(payload);
          break;
        case "donor_completed":
          await handleDonorCompleted(payload);
          break;
        case "user_contact":
          await handleUserContact(payload);
          break;
        default:
          pinologger.warn({ type }, "Unknown event type");
      }
    },
    [channel("notifications")]: async (message: string) => {
      const event = JSON.parse(message);
      await handleNotification(event);
    },
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

  // Redis subscriptions — только events и notifications
  subscribeToChannel(channel("events"));
  subscribeToChannel(channel("notifications"));

  redis.on("message", async (channel, message) => {
    const handler = eventHandlers[channel];
    if (handler) {
      try {
        await handler(message);
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
