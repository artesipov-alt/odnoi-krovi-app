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
  bot.use(logger, limitter);

  // Установка команд бота
  await bot.api.setMyCommands([
    { command: "start", description: "Запустить бота" },
    { command: "profile", description: "Профиль пользователя" },
    { command: "help", description: "Помощь" },
  ]);

  bot.api.config.use();

  //Команды бота
  bot.command("start", startHandler);
  bot.command("help", helpHandler);
  bot.command("profile", profileHandler);
  bot.command("err", errCommandTest);
  bot.command("api", apiTestHandler);

  //Колбэки (нажатия на кнопки)
  bot.callbackQuery("profile", profileHandler);
  bot.callbackQuery("help", helpHandler);
  bot.callbackQuery("back", startHandler);

  const { first_name, last_name, id } = await bot.api.getMe();

  pinologger.info(
    `Бот ${first_name || id}${last_name ? ` ${last_name}` : ""} запущен`,
  );

  //Бот работает через раннер для паралельности задач
  run(bot);

  //Обработка ошибок
  bot.catch(errorHandler);
}

main();
