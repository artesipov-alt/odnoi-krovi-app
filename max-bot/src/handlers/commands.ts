import type { Context } from "@maxhub/max-bot-api";
import { MaxError, Keyboard } from "@maxhub/max-bot-api";
import { Templates } from "../config/templates";
import { usersApi, pinologger } from "../instances";

const getMainKeyboard = () => {
  return Keyboard.inlineKeyboard([
    [Keyboard.button.link("🩸 Открыть приложение", Bun.env.MINIAPP_DOMAIN!)],
    [
      Keyboard.button.callback("❓ Помощь", "help"),
      Keyboard.button.callback("👤 Профиль", "profile"),
    ],
  ]);
};

const getBackKeyboard = () => {
  return Keyboard.inlineKeyboard([
    [Keyboard.button.callback("⬅️ Назад", "back")],
  ]);
};

export const startHandler = async (ctx: Context) => {
  const keyboard = getMainKeyboard();
  const isCommand = ctx.message?.body.text === "/start";
  const maxId = ctx.message!.sender!.user_id!;

  if (isCommand) {
    try {
      // Проверяем существование пользователя
      let isUserExist = false;
      usersApi
        .getUserByTelegram({ id: maxId })
        .then(() => {
          isUserExist = true;
          pinologger.info({ maxId }, "User exists");
        })
        .catch((error: any) => {
          // Если пользователь не найден (404 или 500), регистрируем его
          pinologger.warn(
            { maxId, error: error.message },
            "User not found, will register",
          );
          isUserExist = false;
        });

      if (!isUserExist) {
        const fullName = getFullName(ctx.user);
        usersApi
          .registerUserSimple({
            createUserBody: {
              telegramId: maxId,
              fullName,
            },
          })
          .then(() => {
            pinologger.info(
              { maxId, fullName },
              "User registered successfully",
            );
          })
          .catch((registerError: any) => {
            pinologger.error(
              { maxId, error: registerError.message },
              "Failed to register user",
            );
            // Продолжаем выполнение, даже если регистрация не удалась
          });
      }
    } catch (error: any) {
      pinologger.error(
        { maxId, error: error.message },
        "Error in user check/registration",
      );
      // Не бросаем ошибку, показываем пользователю стартовое сообщение
    }

    await ctx.reply(Templates.START.MESSAGE, {
      format: "markdown",
      attachments: [keyboard],
    });
  } else {
    await ctx.editMessage({
      text: Templates.START.MESSAGE,
      format: "markdown",
      attachments: [keyboard],
    });
  }
};

// /**
//  * Extracts full name from user data with fallback logic
//  */
const getFullName = (user: Context["user"]): string => {
  //@ts-ignore
  const { first_name = "", last_name = "", username = "" } = user;

  if (first_name || last_name) {
    return `${first_name} ${last_name}`.trim();
  }

  return username || "Unknown";
};

export const helpHandler = async (ctx: Context) => {
  const keyboard = getBackKeyboard();

  if (ctx.message?.body.text === "/help") {
    await ctx.reply(Templates.HELP.MESSAGE, {
      format: "markdown",
      attachments: [keyboard],
    });
  } else {
    await ctx.editMessage({
      text: Templates.HELP.MESSAGE,
      format: "markdown",
      attachments: [keyboard],
    });
  }
};

export const profileHandler = async (ctx: Context) => {
  const keyboard = Keyboard.inlineKeyboard([
    [Keyboard.button.link("✏️ Редактировать профиль", Bun.env.MINIAPP_DOMAIN!)],
    [Keyboard.button.callback("⬅️ Назад", "back")],
  ]);

  if (ctx.message?.body.text === "/profile") {
    await ctx.reply(Templates.PROFILE.MESSAGE, {
      format: "markdown",
      attachments: [keyboard],
    });
  } else {
    await ctx.editMessage({
      text: Templates.PROFILE.MESSAGE,
      format: "markdown",
      attachments: [keyboard],
    });
  }
};

export const apiTestHandler = async (ctx: Context) => {
  const data = await usersApi.getUserById({
    id: "",
  });
  await ctx.reply(data.fullName!);
};

export const errCommandTest = async (ctx: Context) => {
  await ctx.reply("Произошла тестовая ошибка, не переживайте так задумано");
  throw new Error("Тестовая ошибка для проверки");
};
