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
import {
  dispatchEvent,
  EVENT_TYPES,
  type EventEnvelope,
  type EventHandler,
  type EventHandlerMap,
} from "../shared/ts/events";

// Типизированная dispatch-таблица. Record<EventType, ...> гарантирует, что
// при добавлении нового EventType компилятор потребует handler.
// Параметр лямбды типизирован как `any` — иначе TS не сводит `unknown` из
// `EventHandler` к узким event-интерфейсам. Runtime-валидация payload остаётся
// на стороне каждого handler'а.
const eventHandlers: EventHandlerMap = {
  [EVENT_TYPES.BLOOD_REQUEST_CREATED]: ((payload: any) =>
    handleBloodRequestCreated(payload)) as EventHandler,
  [EVENT_TYPES.DONOR_APPLY]: ((payload: any) =>
    handleDonorApply(payload)) as EventHandler,
  [EVENT_TYPES.DONATION_CONFIRMED]: ((payload: any) =>
    handleDonationConfirmed(payload)) as EventHandler,
  [EVENT_TYPES.RECIPIENT_APPLY]: ((payload: any) =>
    handleRecipientApply(payload)) as EventHandler,
  [EVENT_TYPES.DONOR_CANCEL]: ((payload: any) =>
    handleDonorCancel(payload)) as EventHandler,
  [EVENT_TYPES.DONOR_REJECT]: ((payload: any) =>
    handleDonorReject(payload)) as EventHandler,
  [EVENT_TYPES.DONOR_NOT_CONFIRMED]: ((payload: any) =>
    handleDonorNotConfirmed(payload)) as EventHandler,
  [EVENT_TYPES.DONOR_COMPLETED]: ((payload: any) =>
    handleDonorCompleted(payload)) as EventHandler,
  [EVENT_TYPES.USER_CONTACT]: ((payload: any) =>
    handleUserContact(payload)) as EventHandler,
};

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

  // Redis event handlers — маппинг по каналу
  const channelHandlers: Record<string, (message: string) => Promise<void>> = {
    [channel("events")]: async (message: string) => {
      const envelope: EventEnvelope = JSON.parse(message);
      pinologger.info({ type: envelope.type }, "Received event");
      await dispatchEvent(envelope, eventHandlers, (type) =>
        pinologger.warn({ type }, "Unknown event type"),
      );
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
    const handler = channelHandlers[channel];
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
