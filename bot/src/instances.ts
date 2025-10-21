import { Bot } from "grammy";

import type { Context } from "grammy";
import pino from "pino";

export const bot = new Bot<Context>(Bun.env.DEV_BOT_API_KEY!);
export const pinologger = pino({
  transport: {
    target: "pino-pretty",
    options: {
      colorize: true,
    },
  },
});
