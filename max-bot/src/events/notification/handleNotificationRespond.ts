import type { Context } from "@maxhub/max-bot-api";
import { Keyboard } from "@maxhub/max-bot-api";
import { NotificationRespondBodyActionEnum } from "../../../../shared/ts/index";
import { bloodRequestApi, pinologger } from "../../instances";
import { getOrCreateToken, invalidateToken } from "../../utils/authStore";

const BOT_ID =
  Bun.env.ENV === "development" ? "id3200014662_2" : "id3200014662";

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
    if (err?.response?.status === 401 || err?.status === 401) {
      return false;
    }
    throw err;
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
  const maxId = ctx.user?.user_id;
  if (!maxId) {
    pinologger.warn("No user ID in callback query");
    await ctx.answerOnCallback({
      notification: "Не удалось определить пользователя",
    });
    return;
  }

  try {
    const auth = await getOrCreateToken(Number(maxId));
    const authorization = `${auth.tokenType} ${auth.accessToken}`;

    let ok = await callRespond(requestId, action, authorization);

    if (!ok) {
      pinologger.info({ maxId }, "Token expired, refreshing and retrying");
      invalidateToken(Number(maxId));

      const freshAuth = await getOrCreateToken(Number(maxId));
      const freshAuthorization = `${freshAuth.tokenType} ${freshAuth.accessToken}`;

      ok = await callRespond(requestId, action, freshAuthorization);

      if (!ok) {
        throw new Error("Повторный запрос также вернул 401");
      }
    }

    // Успешно — редактируем сообщение
    if (action === "yes") {
      await ctx.editMessage({
        text: "Хорошо, продолжаем поиск!",
        format: "markdown",
        attachments: [], // убираем кнопки
      });
    } else {
      await ctx.editMessage({
        text: "Поиск завершен. При необходимости начните новый поиск для питомца во вкладке «Найти кровь» в приложении.",
        format: "markdown",
        attachments: [
          Keyboard.inlineKeyboard([
            [
              Keyboard.button.link(
                "🩸 Открыть приложение",
                `https://max.ru/${BOT_ID}_bot?startapp`,
              ),
            ],
          ]),
        ],
      });
    }

    pinologger.info(
      { action, requestId, maxId },
      "Notification respond handled successfully",
    );
  } catch (err) {
    pinologger.error(
      { error: err, action, requestId, maxId },
      "Failed to handle notification respond",
    );

    await ctx.answerOnCallback({
      notification: "Произошла ошибка, попробуйте позже",
    });
    return;
  }

  await ctx.answerOnCallback({});
};
