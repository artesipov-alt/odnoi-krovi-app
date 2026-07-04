import { Bot, Context } from "@maxhub/max-bot-api";
import { pinologger } from "./instances";

const WEBHOOK_PATH = "/webhook";
const WEBHOOK_SECRET = process.env.WEBHOOK_SECRET;

/**
 * Проверяет секрет вебхука
 */
function verifySecret(req: Request): boolean {
  if (!WEBHOOK_SECRET) return true; // Если секрет не задан в env, пропускаем проверку
  const secretHeader = req.headers.get("X-Max-Bot-Api-Secret");
  return secretHeader === WEBHOOK_SECRET;
}

/**
 * Запускает HTTP сервер на Bun для обработки вебхуков
 */
export function startServer(bot: Bot<Context>) {
  const server = Bun.serve({
    port: 6000,
    routes: {
      // Маршрут для проверки статуса
      "/api/status": new Response("OK"),

      // Основной маршрут вебхука
      [WEBHOOK_PATH]: {
        POST: async (req) => {
          // 1. Проверка секрета
          if (!verifySecret(req)) {
            pinologger.warn("Webhook secret mismatch or missing");
            return new Response("Unauthorized", { status: 401 });
          }

          // 2. Обработка обновления
          try {
            const update = await req.json();
            pinologger.info(
              { update_type: (update as any).update_type },
              "Received webhook update",
            );

            // Передаем обновление в бота
            // handleUpdate является приватным в данном SDK, используем any для доступа
            await (bot as any).handleUpdate(update);

            return new Response("OK", { status: 200 });
          } catch (err) {
            pinologger.error({ error: err }, "Error handling webhook update");
            return new Response("Internal Server Error", { status: 500 });
          }
        },
      },
    },

    // Fallback для остальных маршрутов
    fetch(req) {
      return new Response("Not Found", { status: 404 });
    },
  });

  pinologger.info(`Server running at ${server.url}`);
  pinologger.info(`Webhook endpoint: ${server.url}${WEBHOOK_PATH.slice(1)}`);

  return server;
}
