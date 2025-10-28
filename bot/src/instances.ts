import { Bot } from "grammy";
import { UsersApi, Configuration } from "./api/index";

import type { Context } from "grammy";
import pino from "pino";

export const bot = new Bot<Context>(Bun.env.BOT_TOKEN!);
export const pinologger = pino({
  transport: {
    target: "pino-pretty",
    options: {
      colorize: true,
    },
  },
});

const API_BASE_URL = "https://1krovi.app/api/v1";

const config = new Configuration({
  basePath: API_BASE_URL,
});

// Явно передаем basePath как второй и третий параметры
export const userApi = new UsersApi(config, API_BASE_URL, fetch);
