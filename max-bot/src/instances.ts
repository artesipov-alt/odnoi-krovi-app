import { Bot } from "@maxhub/max-bot-api";
import { Configuration, AuthV1Api } from "../../shared/ts/index";

import type { Context } from "@maxhub/max-bot-api";
import pino from "pino";

export const bot = new Bot<Context>(Bun.env.MAX_BOT_TOKEN!);
export const pinologger = pino({
  level: "debug",
  transport: {
    target: "pino-pretty",
    options: {
      colorize: true,
    },
  },
});

// API Configuration
const apiConfig = new Configuration({
  basePath: Bun.env.API_BASE_URL || "http://localhost:8080/api/v1",
  headers: {
    "Content-Type": "application/json",
  },
  // Добавьте middleware для логирования, если нужно
  middleware: [
    {
      pre: async (context: any) => {
        pinologger.debug(
          { url: context.url, method: context.init.method },
          "API Request",
        );
        return context;
      },
      post: async (context: any) => {
        pinologger.debug(
          { url: context.url, status: context.response.status },
          "API Response",
        );
        return context.response;
      },
    },
  ],
});

// API Client Instances
export const usersApi = new AuthV1Api(apiConfig);
