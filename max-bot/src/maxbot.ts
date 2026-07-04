import { pinologger } from "./instances";

/**
 * Полный список типов обновлений, на которые должен быть подписан бот.
 *
 * Критично: в отличие от Telegram, Max API не присылает обновления, на
 * которые бот явно не подписан через `POST /subscriptions`. Без
 * `message_callback` события о нажатиях на inline-кнопки не доходят до бота.
 */
const UPDATE_TYPES = [
  "message_created",
  "message_callback",
  "message_removed",
  "message_edited",
  "bot_started",
  "bot_added",
  "bot_removed",
  "user_added",
  "user_removed",
  "chat_title_changed",
] as const;

const SUBSCRIPTIONS_URL = "https://platform-api2.max.ru/subscriptions";

interface RegisterWebhookResult {
  success: boolean;
  message?: string;
}

/**
 * Регистрирует (или перерегистрирует) вебхук для Max-бота.
 *
 * Идемпотентно: повторный вызов с тем же URL просто заменяет текущую
 * подписку. Используется для восстановления после деплоя и при первом
 * старте бота — без этого callback'и от inline-кнопок не приходят.
 *
 * Если `webhookUrl` не задан — функция no-op (подписка настраивается
 * вручную, например через masterbot).
 */
export async function registerWebhook(
  token: string | undefined,
  webhookUrl: string | undefined,
  secret: string | undefined,
): Promise<RegisterWebhookResult | null> {
  if (!webhookUrl) {
    pinologger.debug(
      "MAX_BOT_WEBHOOK_URL is not set, skipping webhook registration",
    );
    return null;
  }

  if (!token) {
    pinologger.warn(
      "MAX_BOT_TOKEN is not set, cannot register webhook",
    );
    return null;
  }

  pinologger.info({ webhookUrl }, "Registering Max webhook subscription");

  try {
    const res = await fetch(SUBSCRIPTIONS_URL, {
      method: "POST",
      headers: {
        Authorization: token,
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        url: webhookUrl,
        update_types: UPDATE_TYPES,
        ...(secret ? { secret } : {}),
      }),
    });

    const data = (await res.json().catch(() => ({}))) as RegisterWebhookResult;

    if (!res.ok || data.success === false) {
      pinologger.error(
        { status: res.status, data, webhookUrl },
        "Failed to register Max webhook",
      );
      return data;
    }

    pinologger.info(
      { webhookUrl, updateTypes: UPDATE_TYPES },
      "Max webhook registered successfully",
    );
    return data;
  } catch (err) {
    pinologger.error(
      { error: err, webhookUrl },
      "Error while registering Max webhook",
    );
    return null;
  }
}
