import type { Context } from "@maxhub/max-bot-api";
import { Keyboard } from "@maxhub/max-bot-api";
import { Templates } from "../config/templates";
import { usersApi, pinologger } from "../instances";
import { getMainKeyboard } from "../keyboards";

// ============ Keyboard Builders ============

const getBackKeyboard = () => {
  return Keyboard.inlineKeyboard([
    [Keyboard.button.callback("⬅️ Назад", "back")],
  ]);
};

// ============ User Service ============

const parsePayload = (
  payload: string | undefined,
): {
  utm_campaign?: string;
  utm_source?: string;
  utm_medium?: string;
  utm_content?: string;
  utm_term?: string;
} => {
  if (!payload) return { utm_campaign: "organic" };

  if (payload.includes("=")) {
    // New referral link format: parse as query string
    const params = new URLSearchParams(payload);
    return {
      utm_campaign: params.get("utm_campaign") || undefined,
      utm_source: params.get("utm_source") || undefined,
      utm_medium: params.get("utm_medium") || undefined,
      utm_content: params.get("utm_content") || undefined,
      utm_term: params.get("utm_term") || undefined,
    };
  } else {
    // Old company link format: payload is utm_campaign
    return { utm_campaign: payload };
  }
};

const authUser = async (
  maxId: number,
  fullName: string,
  utmData: {
    utm_campaign?: string;
    utm_source?: string;
    utm_medium?: string;
    utm_content?: string;
    utm_term?: string;
  },
): Promise<void> => {
  await usersApi.authUserViaService({
    xInternalKey: process.env.INTERNAL_MAX_BOT_SECRET,
    serviceSignInBody: {
      providerName: "max_bot",
      providerId: String(maxId),
      fullName,
      metaData: {
        utm_campaign: utmData.utm_campaign || "organic",
        utm_source: utmData.utm_source || "max_bot",
        utm_medium: utmData.utm_medium || undefined,
        utm_content: utmData.utm_content || undefined,
        utm_term: utmData.utm_term || undefined,
      },
    },
  });
  pinologger.info({ maxId, fullName }, "User authenticated");
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
 * Extracts full name from Max user data with fallback logic
 */
const getFullName = (user: Context["user"]): string => {
  //@ts-ignore
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
  const isCommand =
    ctx.message?.body.text?.startsWith("/start") ||
    ctx.updateType === "bot_started";

  if (isCommand) {
    await ctx.reply(messageText, {
      format: "markdown",
      attachments: [keyboard],
    });
  } else {
    await ctx.editMessage({
      text: messageText,
      format: "markdown",
      attachments: [keyboard],
    });
  }
};

// ============ Handlers ============

export const startHandler = async (ctx: Context) => {
  // Extract user data based on event type
  let maxId: number;
  let fullName: string;
  let payload: string | undefined;

  if (ctx.updateType === "bot_started") {
    maxId = (ctx.update as any).user_id;
    fullName = getFullName((ctx.update as any).user);
    payload = (ctx.update as any).payload;
  } else {
    // Early validation for command messages
    if (!ctx.message?.sender?.user_id) {
      pinologger.error("User ID is not available");
      throw new Error("User ID is not available");
    }
    maxId = Number(ctx.message.sender.user_id);
    fullName = getFullName(ctx.message.sender);
    payload = extractStartPayload(ctx.message.body.text ?? undefined);
  }

  // Prevent registering the bot itself
  if (maxId === ctx.botInfo?.user_id) {
    pinologger.warn({ maxId }, "Attempted to register bot as user, skipping");
    return;
  }

  pinologger.info(
    { maxId, fullName, payload, updateType: ctx.updateType },
    "Start handler data",
  );
  if (ctx.message?.body?.attachments) {
    pinologger.info(
      { attachments: JSON.stringify(ctx.message.body.attachments, null, 2) },
      "Detailed attachments",
    );
  }

  const utmData = parsePayload(payload);

  // Make authentication non-blocking with timeout to prevent delays
  const authPromise = authUser(maxId, fullName, utmData);
  const timeoutPromise = new Promise((_, reject) =>
    setTimeout(() => reject(new Error("Authentication timeout")), 5000),
  );

  try {
    await Promise.race([authPromise, timeoutPromise]);
    pinologger.info({ maxId }, "User authentication successful");
  } catch (error: any) {
    pinologger.error(
      { maxId, error: error.message },
      "Error in user authentication or timeout",
    );
    // Don't throw - still show welcome message
  }

  // Always show welcome message
  await sendStartResponse(ctx);
};

export const helpHandler = async (ctx: Context) => {
  const keyboard = getBackKeyboard();
  const isCommand = ctx.message?.body.text === "/help";

  if (isCommand) {
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
    [Keyboard.button.callback("⬅️ Назад", "back")],
  ]);

  const isCommand = ctx.message?.body.text === "/profile";

  if (isCommand) {
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
  await ctx.reply("API test handler not implemented for AuthV1Api");
};

export const errCommandTest = async (ctx: Context) => {
  await ctx.reply("Произошла тестовая ошибка, не переживайте так задумано");
  throw new Error("Тестовая ошибка для проверки");
};
