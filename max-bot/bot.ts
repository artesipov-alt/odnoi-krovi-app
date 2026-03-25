import {
  helpHandler,
  startHandler,
  profileHandler,
  errCommandTest,
  apiTestHandler,
} from "./src/handlers/commands";
import { handleDonorApply } from "./src/handlers/events";
import { bot, pinologger, redis } from "./src/instances";
import { logger } from "./src/middleware/logger";
import { errorHandler } from "./src/handlers/errors";

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

  bot.start();

  // Redis subscription for events
  redis.subscribe("donor_response.apply", (err, count) => {
    if (err) {
      pinologger.error({ error: err }, "Failed to subscribe to Redis channel");
    } else {
      pinologger.info(`Subscribed to ${count} channel(s)`);
    }
  });

  redis.on("message", (channel, message) => {
    if (channel === "donor_response.apply") {
      try {
        const event = JSON.parse(message);
        handleDonorApply(event);
      } catch (err) {
        pinologger.error({ error: err }, "Failed to parse event");
      }
    }
  });

  //Обработка ошибок
  bot.catch(errorHandler);
}

main();
