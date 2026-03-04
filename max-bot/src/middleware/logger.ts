import type { Context, NextFunction } from "grammy";
import { pinologger } from "../instances";

export async function logger(ctx: Context, next: NextFunction) {
  const startTime = Date.now();

  try {
    await next();
    const processingTime = Date.now() - startTime;

    pinologger.info({
      updateId: ctx.update.update_id,
      messageType: getMessageType(ctx),
      chatId: ctx.chat?.id,
      user: getUserDisplayName(ctx.from),
      processingTime: `${processingTime}ms`,
      timestamp: new Date().toISOString(),
    });
  } catch (error) {
    throw error;
  }
}

function getUserDisplayName(from: Context["from"]): string {
  if (!from) return "unknown";

  if (from.first_name || from.last_name) {
    return [from.first_name, from.last_name].filter(Boolean).join(" ");
  }

  return from.username || "unknown";
}

function getMessageType(ctx: Context): string {
  if (ctx.message?.text) {
    // Check if it's a command
    if (ctx.message.text.startsWith("/")) {
      return "command";
    }
    return "text";
  }
  if (ctx.message?.photo) return "photo";
  if (ctx.message?.document) return "document";
  if (ctx.callbackQuery) return "callback_query";
  if (ctx.inlineQuery) return "inline_query";
  return "unknown";
}
