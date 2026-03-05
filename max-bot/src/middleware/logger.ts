import type { Context, NextFn } from "@maxhub/max-bot-api";
import { pinologger } from "../instances";

export async function logger(ctx: Context, next: NextFn) {
  const startTime = Date.now();

  try {
    await next();
    const processingTime = Date.now() - startTime;

    pinologger.info({
      updateId: ctx.update?.timestamp || "unknown",
      messageType: ctx.updateType,
      chatId: ctx.chat?.chat_id,
      user: getUserDisplayName(ctx.message),
      processingTime: `${processingTime}ms`,
      timestamp: new Date().toISOString(),
    });
  } catch (error) {
    throw error;
  }
}

function getUserDisplayName(from: Context["message"]): string {
  if (!from) return "unknown";

  if (from.sender?.name) {
    return from.sender.name;
  }

  return from.sender?.username || "unknown";
}
