// Общие типы и диспатчер для событий Redis-канала `events`.
// События публикуются бэкендом в формате EventEnvelope и разбираются обоими ботами.

export const EVENT_TYPES = {
  BLOOD_REQUEST_CREATED: "blood_request_created",
  DONOR_APPLY: "donor_response_apply",
  DONOR_SELECTED: "donor_selected",
  DONATION_CONFIRMED: "donation_confirmed",
  RECIPIENT_APPLY: "recipient_response_apply",
  DONOR_CANCEL: "donor_cancel",
  DONOR_REJECT: "donor_reject",
  DONOR_NOT_CONFIRMED: "donor_not_confirmed",
  DONOR_COMPLETED: "donor_completed",
  USER_CONTACT: "user_contact",
} as const;

export type EventType = (typeof EVENT_TYPES)[keyof typeof EVENT_TYPES];

export interface EventEnvelope<P = unknown> {
  type: EventType;
  payload: P;
  createdAt: string;
}

export type EventHandler<P = unknown> = (payload: P) => Promise<void> | void;

// Типизированная dispatch-таблица. Строгий Record<EventType, ...> даёт
// exhaustiveness-проверку: добавление нового EventType в shared приведёт к
// ошибке компиляции в обоих ботах, пока они не зарегистрируют handler.
export type EventHandlerMap = Record<EventType, EventHandler>;

// Достаёт handler по `envelope.type` и вызывает его с payload.
// При неизвестном типе вызывает `onUnknown` (если задан) и возвращается.
export const dispatchEvent = async (
  envelope: EventEnvelope,
  handlers: EventHandlerMap,
  onUnknown?: (type: string) => void,
): Promise<void> => {
  const handler = handlers[envelope.type];
  if (!handler) {
    onUnknown?.(envelope.type);
    return;
  }
  await handler(envelope.payload);
};
