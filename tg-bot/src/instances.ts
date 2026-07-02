import { Bot } from "grammy";
import {
  Configuration,
  AuthV1Api,
  AdminV1Api,
  UsersV1Api,
  BloodRequestV1Api,
} from "../../shared/ts/index";

import type { Context } from "grammy";
import pino from "pino";
import Redis from "ioredis";

export const bot = new Bot<Context>(Bun.env.TG_BOT_TOKEN!, {
  client: {
    apiRoot: "https://bridge.1krovi.app",
    buildUrl: (root, token, method) => {
      // Собираем стандартный URL, но докидываем в конец наш секрет
      return `${root}/bot${token}/${method}?secret=${Bun.env.BRIDGE_TOKEN}`;
    },
  },
});

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
export const authApi = new AuthV1Api(apiConfig);
export const adminApi = new AdminV1Api(apiConfig);
export const userApi = new UsersV1Api(apiConfig);
export const bloodRequestApi = new BloodRequestV1Api(apiConfig);

const redisHost = Bun.env.REDIS_HOST || "localhost";
const redisPort = Bun.env.REDIS_PORT || "6379";
const redisUrl = `redis://${redisHost}:${redisPort}`;

export const redis = new Redis(redisUrl, {
  connectTimeout: 5000,
  lazyConnect: true,
  db: Bun.env.ENV === "development" ? 1 : 0,
});
