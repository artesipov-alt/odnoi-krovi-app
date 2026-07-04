import { usersApi, pinologger } from "../instances";

interface AuthData {
  tokenType: string;
  accessToken: string;
  expiresAt: string; // ISO string
  userId: string;
}

const store = new Map<number, AuthData>();

/**
 * Возвращает токен для Max ID из кэша или создаёт новый через authUserViaService.
 * Если сохранённый токен протух — удаляет его и создаёт новый.
 */
export async function getOrCreateToken(maxId: number): Promise<AuthData> {
  const existing = store.get(maxId);

  if (existing) {
    const expiresAt = new Date(existing.expiresAt).getTime();
    const now = Date.now();

    // Токен ещё жив (с запасом 5 минут)
    if (expiresAt > now + 5 * 60 * 1000) {
      pinologger.debug({ maxId }, "Using cached token");
      return existing;
    }

    pinologger.info({ maxId }, "Cached token expired, re-authenticating");
  }

  const result = await usersApi.authUserViaService({
    xInternalKey: process.env.INTERNAL_MAX_BOT_SECRET!,
    serviceSignInBody: {
      providerName: "max_bot",
      providerId: String(maxId),
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

  store.set(maxId, data);

  pinologger.info({ maxId }, "New token acquired and cached");
  return data;
}

/**
 * Принудительно удаляет токен из кэша (например, при 401 ошибке).
 */
export function invalidateToken(maxId: number): void {
  store.delete(maxId);
  pinologger.info({ maxId }, "Token invalidated");
}
