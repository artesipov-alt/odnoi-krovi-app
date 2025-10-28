import type { Context } from "grammy";
import { BotError, InlineKeyboard } from "grammy";
import { Templates } from "../config/templates";
import { usersApi, pinologger } from "../instances";

const getMainKeyboard = () => {
  return new InlineKeyboard()
    .webApp("🩸 Открыть приложение", Bun.env.MINIAPP_DOMAIN!)
    .row()
    .text("❓ Помощь", "help")
    .text("👤 Профиль", "profile");
};

const getBackKeyboard = () => {
  return new InlineKeyboard().text("⬅️ Назад", "back");
};

export const startHandler = async (ctx: Context) => {
  const keyboard = getMainKeyboard();
  const isCommand = ctx.message?.text === "/start";

  // Early validation
  if (!ctx.from?.id) {
    throw new BotError("User ID is not available", ctx);
  }

  const telegramId = ctx.from.id;

  if (isCommand) {
    try {
      // Проверяем существование пользователя
      let isUserExist = false;

      try {
        await usersApi.userTelegramGet({ telegramId });
        isUserExist = true;
        pinologger.info({ telegramId }, "User exists");
      } catch (error: any) {
        // Если пользователь не найден (404 или 500), регистрируем его
        pinologger.warn(
          { telegramId, error: error.message },
          "User not found, will register",
        );
        isUserExist = false;
      }

      if (!isUserExist) {
        try {
          const fullName = getFullName(ctx.from);

          await usersApi.userRegisterSimplePost({
            request: {
              telegramId,
              fullName,
            },
          });
          pinologger.info(
            { telegramId, fullName },
            "User registered successfully",
          );
        } catch (registerError: any) {
          pinologger.error(
            { telegramId, error: registerError.message },
            "Failed to register user",
          );
          // Продолжаем выполнение, даже если регистрация не удалась
        }
      }
    } catch (error: any) {
      pinologger.error(
        { telegramId, error: error.message },
        "Error in user check/registration",
      );
      // Не бросаем ошибку, показываем пользователю стартовое сообщение
    }

    await ctx.reply(Templates.START.MESSAGE, {
      parse_mode: "Markdown",
      reply_markup: keyboard,
    });
  } else {
    await ctx.editMessageText(Templates.START.MESSAGE, {
      parse_mode: "Markdown",
      reply_markup: keyboard,
    });
  }
};

/**
 * Extracts full name from Telegram user data with fallback logic
 */
const getFullName = (user: NonNullable<Context["from"]>): string => {
  const { first_name = "", last_name = "", username = "" } = user;

  if (first_name || last_name) {
    return `${first_name} ${last_name}`.trim();
  }

  return username || "Unknown";
};

export const helpHandler = async (ctx: Context) => {
  const keyboard = getBackKeyboard();

  if (ctx.message?.text === "/help") {
    await ctx.reply(Templates.HELP.MESSAGE, {
      parse_mode: "Markdown",
      reply_markup: keyboard,
    });
  } else {
    await ctx.editMessageText(Templates.HELP.MESSAGE, {
      parse_mode: "Markdown",
      reply_markup: keyboard,
    });
  }
};

export const profileHandler = async (ctx: Context) => {
  const keyboard = new InlineKeyboard()
    .webApp("✏️ Редактировать профиль", Bun.env.MINIAPP_DOMAIN!)
    .row()
    .text("⬅️ Назад", "back");

  if (ctx.message?.text === "/profile") {
    await ctx.reply(Templates.PROFILE.MESSAGE, {
      parse_mode: "Markdown",
      reply_markup: keyboard,
    });
  } else {
    await ctx.editMessageText(Templates.PROFILE.MESSAGE, {
      parse_mode: "Markdown",
      reply_markup: keyboard,
    });
  }
};

export const apiTestHandler = async (ctx: Context) => {
  const data = await usersApi.userIdGet({
    id: 7,
  });
  await ctx.reply(data.fullName!);
};

export const errCommandTest = async (ctx: Context) => {
  await ctx.reply("Произошла тестовая ошибка, не переживайте так задумано");
  throw new Error("Тестовая ошибка для проверки");
};
