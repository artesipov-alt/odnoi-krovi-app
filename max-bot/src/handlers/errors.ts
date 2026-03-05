import { MaxError } from "@maxhub/max-bot-api";
import { pinologger } from "../instances";

export const errorHandler = (err: any) => {
  const ctx = err.ctx;
  if (!ctx) {
    pinologger.error("Error without context");
    return;
  }
  const updateId = ctx.update?.update_id || ctx.update?.timestamp || "unknown";
  pinologger.error(`Error while handling update ${updateId}:`);
  const e = err.error;
  if (e instanceof MaxError) {
    pinologger.error(`Error in request: ${e.description}`);
  }
};
