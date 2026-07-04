import { Bot } from "@maxhub/max-bot-api";
import {
  Configuration,
  AuthV1Api,
  BloodRequestV1Api,
} from "../../shared/ts/index";

import type { Context } from "@maxhub/max-bot-api";
import pino from "pino";
import Redis from "ioredis";

// Max API до 19 июля 2026 переключается с `platform-api.max.ru` на
// `platform-api2.max.ru`. SDK по умолчанию использует старый домен, поэтому
// передаём `baseUrl` явно — иначе после дедлайна `sendMessage`, `editMessage`,
// `getMyInfo` и т.д. перестанут работать.
export const MAX_API_BASE_URL = "https://platform-api2.max.ru";

export const bot = new Bot<Context>(process.env.MAX_BOT_TOKEN!, {
  clientOptions: {
    baseUrl: MAX_API_BASE_URL,
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

const redisHost = process.env.REDIS_HOST || "localhost";
const redisPort = process.env.REDIS_PORT || "6379";
const redisUrl = `redis://${redisHost}:${redisPort}`;

export const redis = new Redis(redisUrl, {
  connectTimeout: 5000,
  lazyConnect: true,
  db: process.env.ENV === "development" ? 1 : 0,
});

// API Configuration
const apiConfig = new Configuration({
  basePath: process.env.API_BASE_URL || "http://localhost:8080/api/v1",
  headers: {
    "Content-Type": "application/json",
  },
  fetchApi: (url: string, init?: RequestInit) =>
    fetch(url, { ...init, signal: AbortSignal.timeout(5000) }),
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
export const bloodRequestApi = new BloodRequestV1Api(apiConfig);
