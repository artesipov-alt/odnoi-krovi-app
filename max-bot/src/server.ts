import {
  createServer,
  type IncomingMessage,
  type ServerResponse,
} from "node:http";
import { Bot, Context } from "@maxhub/max-bot-api";
import { pinologger } from "./instances";

const WEBHOOK_PATH = "/webhook";
const WEBHOOK_SECRET = process.env.WEBHOOK_SECRET;
const PORT = Number(process.env.PORT ?? 6000);

/**
 * Проверяет секрет вебхука
 */
function verifySecret(req: IncomingMessage): boolean {
  if (!WEBHOOK_SECRET) return true; // Если секрет не задан в env, пропускаем проверку
  const secretHeader = req.headers["x-max-bot-api-secret"];
  return secretHeader === WEBHOOK_SECRET;
}

/**
 * Читает тело запроса из IncomingMessage в строку.
 * Max Bot SDK ожидает уже распарсенный JSON-объект (см. server.ts в Bun-версии).
 */
function readBody(req: IncomingMessage): Promise<string> {
  return new Promise((resolve, reject) => {
    const chunks: Buffer[] = [];
    req.on("data", (chunk: Buffer) => chunks.push(chunk));
    req.on("end", () => resolve(Buffer.concat(chunks).toString("utf-8")));
    req.on("error", reject);
  });
}

/**
 * Запускает HTTP сервер на Node.js для обработки вебхуков Max API.
 *
 * Семантика идентична прежней реализации на `Bun.serve`:
 * - `GET /api/status` → "OK"
 * - `POST /webhook` → проверка секрета → JSON.parse → bot.handleUpdate
 * - всё остальное → 404
 */
export function startServer(bot: Bot<Context>) {
  const server = createServer(async (req, res) => {
    try {
      if (req.method === "GET" && req.url === "/api/status") {
        res.statusCode = 200;
        res.end("OK");
        return;
      }

      if (req.method === "POST" && req.url === WEBHOOK_PATH) {
        if (!verifySecret(req)) {
          pinologger.warn("Webhook secret mismatch or missing");
          res.statusCode = 401;
          res.end("Unauthorized");
          return;
        }

        try {
          const raw = await readBody(req);
          const update = JSON.parse(raw);
          pinologger.info(
            { update_type: (update as any).update_type },
            "Received webhook update",
          );
          // handleUpdate приватный в SDK, используем any для доступа
          await (bot as any).handleUpdate(update);
          res.statusCode = 200;
          res.end("OK");
        } catch (err) {
          pinologger.error({ error: err }, "Error handling webhook update");
          res.statusCode = 500;
          res.end("Internal Server Error");
        }
        return;
      }

      res.statusCode = 404;
      res.end("Not Found");
    } catch (err) {
      pinologger.error({ error: err }, "Unhandled error in HTTP server");
      if (!res.headersSent) {
        res.statusCode = 500;
        res.end("Internal Server Error");
      } else {
        res.end();
      }
    }
  });

  server.listen(PORT, () => {
    pinologger.info(`Server running at http://localhost:${PORT}`);
    pinologger.info(
      `Webhook endpoint: http://localhost:${PORT}${WEBHOOK_PATH}`,
    );
  });

  // Совместимость с прежним контрактом startServer — возвращаем объект
  // с .stop() и .url, чтобы Graceful shutdown в bot.ts не пришлось
  // сильно переделывать.
  return Object.assign(server, {
    url: `http://localhost:${PORT}`,
    stop: (cb?: () => void) => server.close(() => cb?.()),
  }) as unknown as { stop: (cb?: () => void) => void; url: string };
}
