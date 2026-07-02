import { authApi, pinologger } from "../instances";

interface AuthData {
  tokenType: string;
  accessToken: string;
  expiresAt: string; // ISO string
  userId: string;
}

const store = new Map<number, AuthData>();

/**
 * Возвращает токен для Telegram ID из кэша или создаёт новый через authUserViaService.
 * Если сохранённый токен протух — удаляет его и создаёт новый.
 */
export async function getOrCreateToken(telegramId: number): Promise<AuthData> {
  const existing = store.get(telegramId);

  if (existing) {
    const expiresAt = new Date(existing.expiresAt).getTime();
    const now = Date.now();

    // Токен ещё жив (с запасом 5 минут)
    if (expiresAt > now + 5 * 60 * 1000) {
      pinologger.debug({ telegramId }, "Using cached token");
      return existing;
    }

    pinologger.info({ telegramId }, "Cached token expired, re-authenticating");
  }

  const result = await authApi.authUserViaService({
    xInternalKey: Bun.env.INTERNAL_TG_BOT_SECRET!,
    serviceSignInBody: {
      providerName: "telegram_bot",
      providerId: String(telegramId),
    },
  });

  if (!result.accessToken) {
    throw new Error("authUserViaService returned no accessToken");
  }

  const data: AuthData = {
    tokenType: result.tokenType!,
    accessToken: result.accessToken,
    expiresAt: result.expiresAt,
    userId: result.userId,
  };

  store.set(telegramId, data);

  pinologger.info({ telegramId }, "New token acquired and cached");
  return data;
}

/**
 * Принудительно удаляет токен из кэша (например, при 401 ошибке).
 */
export function invalidateToken(telegramId: number): void {
  store.delete(telegramId);
  pinologger.info({ telegramId }, "Token invalidated");
}
