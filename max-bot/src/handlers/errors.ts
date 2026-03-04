import { GrammyError, HttpError } from "grammy";
import { BotError } from "grammy";
import { pinologger } from "../instances";

export const errorHandler = (err: any) => {
  const ctx = err.ctx;
  pinologger.error(`Error while handling update ${ctx.update.update_id}:`);
  const e = err.error;
  if (e instanceof GrammyError) {
    pinologger.error(`Error in request: ${e.description}`);
  } else if (e instanceof HttpError) {
    pinologger.error(`Could not contact Telegram: ${e}`);
  } else if (e instanceof BotError) {
    pinologger.error(`Bot error: ${e}`);
  } else {
    pinologger.error(`Unknown error: ${e}`);
  }
};
