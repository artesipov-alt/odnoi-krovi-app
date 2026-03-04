//import Redis from "ioredis";
import { limit } from "@grammyjs/ratelimiter";
import type { Context } from "grammy";

//Добавить Redis позже для масштабирования
//const redis = new Redis(...);

export const limitter = limit({
  // Allow only 3 messages to be handled every 2 seconds.
  timeFrame: 1500,
  limit: 2,

  // "MEMORY_STORE" is the default value. If you do not want to use Redis, do not pass storageClient at all.
  storageClient: "MEMORY_STORE",

  // This is called when the limit is exceeded.
  onLimitExceeded: async (ctx: Context) => {
    await ctx.reply("Вы превысили число запросов, пожалуйста подождите!");
  },

  // Note that the key should be a number in string format such as "123456789".
  keyGenerator: (ctx) => {
    return ctx.from?.id.toString();
  },
});
