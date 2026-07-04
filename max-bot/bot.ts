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
import { handleNotificationRespond } from "./src/events/notification/handleNotificationRespond";

import { bot, pinologger, redis } from "./src/instances";
import { logger } from "./src/middleware/logger";
import { errorHandler } from "./src/handlers/errors";
import { startServer } from "./src/server";
import { registerWebhook } from "./src/maxbot";
import {
  dispatchEvent,
  EVENT_TYPES,
  type EventEnvelope,
  type EventHandler,
  type EventHandlerMap,
} from "../shared/ts/events";

const env = process.env.ENV || "production";
const channelPrefix = env === "development" || env === "dev" ? "dev:" : "prod:";
const channel = (name: string) => `${channelPrefix}${name}`;

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

  // Обработка ответов на уведомления (п.5)
  bot.action(/^notification_(yes|no)_(.+)$/, async (ctx) => {
    if (!ctx.match || ctx.match.length < 3) {
      pinologger.warn(
        { match: ctx.match },
        "Invalid notification callback match",
      );
      await ctx.answerOnCallback({});
      return;
    }
    const [, action, requestId] = ctx.match;
    await handleNotificationRespond(ctx, action!, requestId!);
  });

  const { name, username, user_id } = await bot.api.getMyInfo();
  pinologger.info(`Бот ${name || username} ${user_id} инициализирован`);

  // Регистрация вебхука: без `message_callback` в update_types Max не
  // доставляет события о нажатиях на inline-кнопки (см. src/maxbot.ts).
  await registerWebhook(
    process.env.MAX_BOT_TOKEN,
    process.env.MAX_BOT_WEBHOOK_URL,
    process.env.WEBHOOK_SECRET,
  );

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

  // Запуск HTTP сервера для вебхуков
  const server = startServer(bot);

  // Обработка ошибок
  bot.catch(errorHandler);

  // Graceful shutdown
  const shutdown = async () => {
    pinologger.info("Shutting down...");
    server.stop();
    await redis.quit();
    process.exit(0);
  };

  process.on("SIGINT", shutdown);
  process.on("SIGTERM", shutdown);
}

main();
