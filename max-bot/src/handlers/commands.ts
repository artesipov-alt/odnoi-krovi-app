import type { Context } from "@maxhub/max-bot-api";
import { Keyboard } from "@maxhub/max-bot-api";
import { Templates } from "../config/templates";
import { usersApi, pinologger } from "../instances";

// ============ Keyboard Builders ============

const getMainKeyboard = () => {
  return Keyboard.inlineKeyboard([
    [
      Keyboard.button.link(
        "🩸 Открыть приложение",
        "https://max.ru/id3200014662_2_bot?startapp",
      ),
    ],
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

// ============ User Service ============

const checkUserExists = async (maxId: number): Promise<boolean> => {
  return usersApi
    .authUser({
      authUserBody: {
        providerId: maxId,
        providerName: "max_bot",
      },
    })
    .then(() => {
      pinologger.info({ maxId }, "User exists");
      return true;
    })
    .catch((error) => {
      pinologger.warn(
        { maxId, error: error.message },
        "User not found, will register",
      );
      return false;
    });
};

const registerUser = (
  maxId: number,
  fullName: string,
  utmCampaign?: string,
): Promise<unknown> => {
  return usersApi.registerUserSimple({
    createUserBody: {
      fullName,
      providerId: maxId,
      providerName: "max_bot",
      metaData: {
        utm_campaign: utmCampaign || "organic",
        utm_source: "max_bot",
      },
    },
  });
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
  pinologger.info({ ctx }, "Full ctx");

  try {
    await checkUserExists(maxId).then((userExists) => {
      if (userExists) return;

      return registerUser(maxId, fullName, payload).then(() => {
        pinologger.info(
          { maxId, fullName, utmCampaign: payload },
          "User registered successfully",
        );
      });
    });
  } catch (error: any) {
    pinologger.error(
      { maxId, error: error.message },
      "Error in user check/registration",
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
    [Keyboard.button.link("✏️ Редактировать профиль", Bun.env.MINIAPP_DOMAIN!)],
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
  const data = await usersApi.getUserById({
    id: "",
  });
  await ctx.reply(data.fullName!);
};

export const errCommandTest = async (ctx: Context) => {
  await ctx.reply("Произошла тестовая ошибка, не переживайте так задумано");
  throw new Error("Тестовая ошибка для проверки");
};
