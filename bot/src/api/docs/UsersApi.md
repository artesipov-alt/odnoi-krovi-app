# UsersApi

All URIs are relative to */api/v1*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**userIdDelete**](UsersApi.md#useriddelete) | **DELETE** /user/{id} | Удаление пользователя по ID |
| [**userIdGet**](UsersApi.md#useridget) | **GET** /user/{id} | Получение пользователя по ID |
| [**userIdPut**](UsersApi.md#useridput) | **PUT** /user/{id} | Обновление данных пользователя |
| [**userRegisterPost**](UsersApi.md#userregisterpost) | **POST** /user/register | Регистрация нового пользователя |
| [**userRegisterSimplePost**](UsersApi.md#userregistersimplepost) | **POST** /user/register/simple | Простая регистрация пользователя |
| [**userTelegramGet**](UsersApi.md#usertelegramget) | **GET** /user/telegram | Получение пользователя по Telegram ID |



## userIdDelete

> HandlersSuccessResponse userIdDelete(id)

Удаление пользователя по ID

Удаляет пользователя из системы (soft delete)

### Example

```ts
import {
  Configuration,
  UsersApi,
} from '';
import type { UserIdDeleteRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new UsersApi();

  const body = {
    // string | ID пользователя
    id: id_example,
  } satisfies UserIdDeleteRequest;

  try {
    const data = await api.userIdDelete(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **id** | `string` | ID пользователя | [Defaults to `undefined`] |

### Return type

[**HandlersSuccessResponse**](HandlersSuccessResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Пользователь успешно удален |  -  |
| **400** | Неверный запрос |  -  |
| **404** | Пользователь не найден |  -  |
| **500** | Внутренняя ошибка сервера |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## userIdGet

> ModelsUser userIdGet(id)

Получение пользователя по ID

Возвращает информацию о пользователе по его идентификатору

### Example

```ts
import {
  Configuration,
  UsersApi,
} from '';
import type { UserIdGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new UsersApi();

  const body = {
    // string | ID пользователя
    id: id_example,
  } satisfies UserIdGetRequest;

  try {
    const data = await api.userIdGet(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **id** | `string` | ID пользователя | [Defaults to `undefined`] |

### Return type

[**ModelsUser**](ModelsUser.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Данные пользователя |  -  |
| **400** | Неверный запрос |  -  |
| **404** | Пользователь не найден |  -  |
| **500** | Внутренняя ошибка сервера |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## userIdPut

> HandlersSuccessResponse userIdPut(id, request)

Обновление данных пользователя

Обновляет информацию о пользователе

### Example

```ts
import {
  Configuration,
  UsersApi,
} from '';
import type { UserIdPutRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new UsersApi();

  const body = {
    // string | ID пользователя
    id: id_example,
    // ServicesUserUpdate | Данные для обновления
    request: ...,
  } satisfies UserIdPutRequest;

  try {
    const data = await api.userIdPut(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **id** | `string` | ID пользователя | [Defaults to `undefined`] |
| **request** | [ServicesUserUpdate](ServicesUserUpdate.md) | Данные для обновления | |

### Return type

[**HandlersSuccessResponse**](HandlersSuccessResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Данные успешно обновлены |  -  |
| **400** | Неверный запрос |  -  |
| **404** | Пользователь не найден |  -  |
| **500** | Внутренняя ошибка сервера |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## userRegisterPost

> ModelsUser userRegisterPost(request)

Регистрация нового пользователя

Регистрирует нового пользователя в системе

### Example

```ts
import {
  Configuration,
  UsersApi,
} from '';
import type { UserRegisterPostRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new UsersApi();

  const body = {
    // ServicesUserRegistration | Данные для регистрации пользователя
    request: ...,
  } satisfies UserRegisterPostRequest;

  try {
    const data = await api.userRegisterPost(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **request** | [ServicesUserRegistration](ServicesUserRegistration.md) | Данные для регистрации пользователя | |

### Return type

[**ModelsUser**](ModelsUser.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **201** | Зарегистрированный пользователь |  -  |
| **400** | Неверный запрос |  -  |
| **409** | Пользователь уже существует |  -  |
| **500** | Внутренняя ошибка сервера |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## userRegisterSimplePost

> ModelsUser userRegisterSimplePost(request)

Простая регистрация пользователя

Создает пользователя с Telegram ID и именем (для команды Start)

### Example

```ts
import {
  Configuration,
  UsersApi,
} from '';
import type { UserRegisterSimplePostRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new UsersApi();

  const body = {
    // HandlersSimpleRegistrationRequest | Данные для простой регистрации
    request: ...,
  } satisfies UserRegisterSimplePostRequest;

  try {
    const data = await api.userRegisterSimplePost(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **request** | [HandlersSimpleRegistrationRequest](HandlersSimpleRegistrationRequest.md) | Данные для простой регистрации | |

### Return type

[**ModelsUser**](ModelsUser.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **201** | Зарегистрированный пользователь |  -  |
| **400** | Неверный запрос |  -  |
| **409** | Пользователь уже существует |  -  |
| **500** | Внутренняя ошибка сервера |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## userTelegramGet

> ModelsUser userTelegramGet(telegramId)

Получение пользователя по Telegram ID

Возвращает информацию о пользователе по его Telegram ID

### Example

```ts
import {
  Configuration,
  UsersApi,
} from '';
import type { UserTelegramGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new UsersApi();

  const body = {
    // number | Telegram ID пользователя
    telegramId: 789,
  } satisfies UserTelegramGetRequest;

  try {
    const data = await api.userTelegramGet(body);
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters


| Name | Type | Description  | Notes |
|------------- | ------------- | ------------- | -------------|
| **telegramId** | `number` | Telegram ID пользователя | [Defaults to `undefined`] |

### Return type

[**ModelsUser**](ModelsUser.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Данные пользователя |  -  |
| **400** | Неверный запрос |  -  |
| **404** | Пользователь не найден |  -  |
| **500** | Внутренняя ошибка сервера |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

