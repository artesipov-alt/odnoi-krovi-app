import { Context } from "grammy";
import { InlineKeyboard } from "grammy";
import type { ApiResponse } from "../../../../shared/ts/runtime";
import { NotificationRespondBodyActionEnum } from "../../../../shared/ts/index";
import { bloodRequestApi, pinologger } from "../../instances";
import { getOrCreateToken, invalidateToken } from "../../utils/authStore";

const OPEN_APP_URL =
  Bun.env.ENV === "development" || Bun.env.ENV === "dev"
    ? "https://dev.1krovi.app"
    : "https://1krovi.app";

/**
 * Выполняет respondToNotification с переданным токеном.
 * Возвращает true в случае успеха, false — если пришёл 401.
 */
async function callRespond(
  requestId: string,
  action: string,
  authorization: string,
): Promise<boolean> {
  try {
    await bloodRequestApi.respondToNotification(
      {
        reqId: requestId,
        notificationRespondBody: {
          action: action as NotificationRespondBodyActionEnum,
        },
      },
      {
        headers: {
          Authorization: authorization,
        },
      },
    );
    return true;
  } catch (err: any) {
    // Пробуем вытащить статус из ошибки
    if (err?.response?.status === 401 || err?.status === 401) {
      return false; // протухший токен
    }
    throw err; // остальные ошибки прокидываем
  }
}

/**
 * Обрабатывает нажатие на кнопки "Да" / "Нет" в уведомлении recipient_empty_showcase.
 */
export const handleNotificationRespond = async (
  ctx: Context,
  action: string,
  requestId: string,
) => {
  const chatId = ctx.from?.id;
  if (!chatId) {
    pinologger.warn("No chat ID in callback query");
    await ctx.answerCallbackQuery({
      text: "Не удалось определить пользователя",
      show_alert: true,
    });
    return;
  }

  try {
    // Получаем токен (из кэша или новый)
    const auth = await getOrCreateToken(chatId);
    const authorization = `${auth.tokenType} ${auth.accessToken}`;

    // Пробуем выполнить запрос
    let ok = await callRespond(requestId, action, authorization);

    // Если 401 — инвалидируем кэш, получаем новый токен, ретраим
    if (!ok) {
      pinologger.info({ chatId }, "Token expired, refreshing and retrying");
      invalidateToken(chatId);

      const freshAuth = await getOrCreateToken(chatId);
      const freshAuthorization = `${freshAuth.tokenType} ${freshAuth.accessToken}`;

      ok = await callRespond(requestId, action, freshAuthorization);

      if (!ok) {
        throw new Error("Повторный запрос также вернул 401");
      }
    }

    // Успешно — редактируем сообщение
    if (action === "yes") {
      await ctx.editMessageText("Хорошо, продолжаем поиск!", {
        parse_mode: "Markdown",
        reply_markup: undefined,
      });
    } else {
      await ctx.editMessageText(
        "Поиск завершен. При необходимости начните новый поиск для питомца во вкладке «Найти кровь» в приложении.",
        {
          parse_mode: "Markdown",
          reply_markup: new InlineKeyboard().webApp(
            "🩸 Открыть приложение",
            OPEN_APP_URL,
          ),
        },
      );
    }

    pinologger.info(
      { action, requestId, chatId },
      "Notification respond handled successfully",
    );
  } catch (err) {
    pinologger.error(
      { error: err, action, requestId, chatId },
      "Failed to handle notification respond",
    );

    await ctx.answerCallbackQuery({
      text: "Произошла ошибка, попробуйте позже",
      show_alert: true,
    });
    return;
  }

  await ctx.answerCallbackQuery();
};
