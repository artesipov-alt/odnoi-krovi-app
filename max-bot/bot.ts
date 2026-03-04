import {
  helpHandler,
  startHandler,
  profileHandler,
  errCommandTest,
  apiTestHandler,
} from "./src/handlers/commands";
import { bot, pinologger } from "./src/instances";
import { logger } from "./src/middleware/logger";
import { run } from "@grammyjs/runner";
import { limitter } from "./src/middleware/ratelimitter";
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

  bot.on(`bot_started`, async (ctx) => {
    pinologger.info(ctx);
  });

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

  //Обработка ошибок
  bot.catch(errorHandler);
}

main();
