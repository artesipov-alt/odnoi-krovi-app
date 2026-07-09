import { ResponseError } from "../../../shared/ts/runtime";

/**
 * Структура ошибки, которую возвращает бэкенд
 */
interface BackendError {
  Code?: string;
  Message?: string;
  message?: string;
  Details?: Record<string, unknown>;
  HTTPStatus?: number;
}

/**
 * Коды ошибок бэкенда для человекочитаемых сообщений
 */
const ERROR_MESSAGES: Record<string, string> = {
  NOT_FOUND: "Не найдено",
  ALREADY_EXISTS: "Уже существует",
  VALIDATION_ERROR: "Ошибка валидации",
  UNAUTHORIZED: "Требуется авторизация",
  FORBIDDEN: "Доступ запрещён",
  INTERNAL_ERROR: "Внутренняя ошибка сервера",
  BAD_REQUEST: "Неверный запрос",
  CONFLICT: "Конфликт",
};

/**
 * Парсит ошибку от API и возвращает сообщение от бэкенда
 */
export const parseApiError = async (error: unknown): Promise<string> => {
  // Проверяем, есть ли Response в ошибке
  let response: Response | undefined;

  if (error instanceof ResponseError) {
    response = error.response;
  } else if (error instanceof Response) {
    response = error;
  } else if (error && typeof error === "object" && "response" in error) {
    response = (error as { response: Response }).response;
  }

  // Если есть response - пробуем достать тело
  if (response) {
    try {
      const body = (await response.clone().json()) as BackendError;

      // Пробуем достать Message или message - это основное сообщение от бэкенда
      const message = body.Message || body.message;
      if (message) {
        return message;
      }
    } catch {
      // Не удалось распарсить JSON - используем статус
    }

    // Фоллбек на статус код
    return getStatusMessage(response.status);
  }

  // Ошибка без response
  if (error instanceof Error) {
    return error.message;
  }

  return "Произошла неизвестная ошибка";
};

/**
 * Возвращает сообщение на основе HTTP статуса
 */
const getStatusMessage = (status: number): string => {
  switch (status) {
    case 400:
      return "Неверный запрос";
    case 401:
      return "Требуется авторизация";
    case 403:
      return "Доступ запрещён";
    case 404:
      return "Ресурс не найден";
    case 409:
      return "Конфликт данных";
    case 422:
      return "Ошибка валидации";
    case 500:
      return "Внутренняя ошибка сервера";
    case 502:
    case 503:
      return "Сервер временно недоступен";
    default:
      return `Ошибка (код ${status})`;
  }
};
