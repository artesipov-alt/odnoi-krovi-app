import type { Context } from "grammy";
import { BotError, InlineKeyboard } from "grammy";
import { Templates } from "../config/templates";
import { userApi } from "../instances";

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

  if (ctx.message?.text === "/start") {
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
  try {
    const data = await userApi.userIdGet(7);
    console.log("API Response:", data);
    console.log("Full Name:", data?.fullName);

    if (data) {
      await ctx.reply(
        `Это ${data.fullName || "нет имени"}\nВесь объект: ${JSON.stringify(data, null, 2)}`,
      );
    } else {
      await ctx.reply("Данные не получены (undefined)");
    }
  } catch (error) {
    console.error("API Error:", error);
    if (error instanceof Response) {
      const text = await error.text();
      await ctx.reply(`Ошибка API: ${error.status} - ${text}`);
    } else {
      await ctx.reply(`Ошибка: ${error}`);
    }
  }
};

export const errCommandTest = async (ctx: Context) => {
  await ctx.reply("Произошла тестовая ошибка, не переживайте так задумано");
  throw new Error("Тестовая ошибка для проверки");
};
