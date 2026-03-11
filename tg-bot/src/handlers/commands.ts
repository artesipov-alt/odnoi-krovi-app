import type { Context } from "grammy";
import { BotError, InlineKeyboard } from "grammy";
import { Templates } from "../config/templates";
import { usersApi, pinologger } from "../instances";

// ============ Keyboard Builders ============

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

// ============ User Service ============

const authUser = async (
  telegramId: number,
  fullName: string,
  utmCampaign?: string,
): Promise<void> => {
  await usersApi.authUserViaService({
    xInternalKey: Bun.env.INTERNAL_TG_BOT_SECRET,
    serviceSignInBody: {
      providerId: String(telegramId),
      fullName,
      metaData: {
        utm_campaign: utmCampaign || "organic",
        utm_source: "telegram_bot",
      },
    },
  });
  pinologger.info({ telegramId, fullName, utmCampaign }, "User authenticated");
};

// ============ Helpers ============

/**
 * Extracts payload from /start command
 * Example: "/start ddl042026" -> "ddl042026"
 */
const extractStartPayload = (text: string | undefined): string | undefined => {
  if (!text) return undefined;

  const parts = text.split(" ");
  if (parts.length < 2) return undefined;

  return parts[1];
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

/**
 * Sends start message (reply or edit based on context)
 */
const sendStartResponse = async (ctx: Context): Promise<void> => {
  const keyboard = getMainKeyboard();
  const messageText = Templates.START.MESSAGE;
  const isCommand = ctx.message?.text?.startsWith("/start");

  if (isCommand) {
    await ctx.reply(messageText, {
      parse_mode: "Markdown",
      reply_markup: keyboard,
    });
  } else {
    await ctx.editMessageText(messageText, {
      parse_mode: "Markdown",
      reply_markup: keyboard,
    });
  }
};

// ============ Handlers ============

export const startHandler = async (ctx: Context) => {
  // Early validation
  if (!ctx.from?.id) {
    throw new BotError("User ID is not available", ctx);
  }

  const telegramId = ctx.from.id;
  const payload = extractStartPayload(ctx.message?.text);

  const fullName = getFullName(ctx.from!);

  try {
    await authUser(telegramId, fullName, payload);
  } catch (error: any) {
    pinologger.error(
      { telegramId, error: error.message },
      "Error in user authentication",
    );
    // Don't throw - still show welcome message
  }

  // Always show welcome message
  await sendStartResponse(ctx);
};

export const helpHandler = async (ctx: Context) => {
  const keyboard = getBackKeyboard();
  const isCommand = ctx.message?.text === "/help";

  if (isCommand) {
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

  const isCommand = ctx.message?.text === "/profile";

  if (isCommand) {
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
  await ctx.reply("API test handler not implemented for AuthV1Api");
};

export const errCommandTest = async (ctx: Context) => {
  await ctx.reply("Произошла тестовая ошибка, не переживайте так задумано");
  throw new Error("Тестовая ошибка для проверки");
};
