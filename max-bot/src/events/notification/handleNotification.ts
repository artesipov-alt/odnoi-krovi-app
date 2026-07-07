import { pinologger } from "../../instances";
import { sendMessageToUser } from "../../max";
import { getAppOpenKeyboard, getNotificationKeyboard } from "../../keyboards";

interface Notification {
  type: string;
  targets: {
    telegramId?: string;
    maxId?: string;
  };
  payload: Record<string, any>;
  createdAt: string;
}

const messages: Record<string, (payload: Record<string, any>) => string> = {
  recipient_donor_waiting: (p) =>
    `На Ваш запрос откликнулся донор *${p.donorPetName}* (группа ${p.donorBloodGroup}). Обсудите с хозяином донора донацию, пока он не передумал.`,

  recipient_inactive_warning: (_p) =>
    `На Ваш поиск откликнулись доноры, но Вы никого из них не выбрали. Обсудите с хозяином подходящего донора детали донации, иначе через 6 часов Ваш поиск будет остановлен, чтобы доноры могли помочь другим нуждающимся.`,

  recipient_search_closed: (_p) =>
    `Вы долго не отвечали донорам, которые хотели Вам помочь, поэтому мы закрыли Вашу заявку. При необходимости начните новый поиск для питомца во вкладке «Найти кровь» в приложении.`,

  donor_not_accepted: (p) =>
    `Вы откликнулись на поиск (*${p.recipientPetName}*, группа ${p.recipientBloodGroup}, ${p.volume} мл), но хозяин реципиента пока не принял Ваше предложение.\nМожете подождать еще немного или отменить донацию и помочь другому питомцу на Портале.`,

  // п.5 — пустая витрина, 24ч не заходил
  recipient_empty_showcase: (p) =>
    `Вас долго не было, Вы еще ищите помощь для питомца *${p.petName}* (группа ${p.bloodGroup}, ${p.volume} мл)?`,

  // п.6 — пустая витрина, 48ч не заходил и не нажал "Да"
  recipient_search_closed_inactive: (_p) =>
    `Вас долго не было, поэтому мы закрыли Вашу заявку. При необходимости начните новый поиск для питомца во вкладке «Найти кровь» в приложении.`,
};

export const handleNotification = async (event: Notification) => {
  const { type, targets, payload } = event;

  if (!targets.maxId || targets.maxId.trim() === "") {
    pinologger.debug(
      { notificationType: type },
      "No maxId in targets, skipping notification",
    );
    return;
  }

  const messageFn = messages[type];
  if (!messageFn) {
    pinologger.warn(
      { notificationType: type, payload },
      "Unknown notification type, skipping",
    );
    return;
  }

  const text = messageFn(payload);

  try {
    // Для recipient_empty_showcase используем клавиатуру с "Да"/"Нет"
    const keyboard =
      type === "recipient_empty_showcase"
        ? getNotificationKeyboard(payload.requestId)
        : getAppOpenKeyboard();

    await sendMessageToUser(targets.maxId, text, {
      attachments: [keyboard],
    });

    pinologger.info(
      {
        maxId: targets.maxId,
        notificationType: type,
      },
      "Sent notification to Max user",
    );
  } catch (err) {
    pinologger.error(
      {
        error: err,
        maxId: targets.maxId,
        notificationType: type,
      },
      "Failed to send notification to Max user",
    );
  }
};
