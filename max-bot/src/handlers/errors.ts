import { MaxError } from "@maxhub/max-bot-api";
import { pinologger } from "../instances";

export const errorHandler = (err: any) => {
  const ctx = err.ctx;
  pinologger.error(
    `Error while handling update ${ctx.update?.update_id || "unknown"}:`,
  );
  const e = err.error;
  if (e instanceof MaxError) {
    pinologger.error(`Error in request: ${e.description}`);
  }
};
