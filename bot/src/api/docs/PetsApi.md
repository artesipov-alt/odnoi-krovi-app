# PetsApi

All URIs are relative to */api/v1*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**petsIdDelete**](PetsApi.md#petsiddelete) | **DELETE** /pets/{id} | Удаление питомца по ID |
| [**petsIdGet**](PetsApi.md#petsidget) | **GET** /pets/{id} | Получение питомца по ID |
| [**petsIdPut**](PetsApi.md#petsidput) | **PUT** /pets/{id} | Обновление данных питомца |
| [**petsUserUserIdGet**](PetsApi.md#petsuseruseridget) | **GET** /pets/user/{user_id} | Получение питомцев пользователя |
| [**petsUserUserIdPost**](PetsApi.md#petsuseruseridpost) | **POST** /pets/user/{user_id} | Создание нового питомца |



## petsIdDelete

> HandlersSuccessResponse petsIdDelete(id)

Удаление питомца по ID

Удаляет питомца из системы

### Example

```ts
import {
  Configuration,
  PetsApi,
} from '';
import type { PetsIdDeleteRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new PetsApi();

  const body = {
    // number | ID питомца
    id: 56,
  } satisfies PetsIdDeleteRequest;

  try {
    const data = await api.petsIdDelete(body);
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
| **id** | `number` | ID питомца | [Defaults to `undefined`] |

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
| **200** | Питомец успешно удален |  -  |
| **400** | Неверный запрос |  -  |
| **404** | Питомец не найден |  -  |
| **500** | Внутренняя ошибка сервера |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## petsIdGet

> ModelsPet petsIdGet(id)

Получение питомца по ID

Возвращает информацию о питомце по его идентификатору

### Example

```ts
import {
  Configuration,
  PetsApi,
} from '';
import type { PetsIdGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new PetsApi();

  const body = {
    // number | ID питомца
    id: 56,
  } satisfies PetsIdGetRequest;

  try {
    const data = await api.petsIdGet(body);
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
| **id** | `number` | ID питомца | [Defaults to `undefined`] |

### Return type

[**ModelsPet**](ModelsPet.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Данные питомца |  -  |
| **400** | Неверный запрос |  -  |
| **404** | Питомец не найден |  -  |
| **500** | Внутренняя ошибка сервера |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## petsIdPut

> HandlersSuccessResponse petsIdPut(id, request)

Обновление данных питомца

Обновляет информацию о питомце

### Example

```ts
import {
  Configuration,
  PetsApi,
} from '';
import type { PetsIdPutRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new PetsApi();

  const body = {
    // number | ID питомца
    id: 56,
    // ServicesPetUpdate | Данные для обновления
    request: ...,
  } satisfies PetsIdPutRequest;

  try {
    const data = await api.petsIdPut(body);
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
| **id** | `number` | ID питомца | [Defaults to `undefined`] |
| **request** | [ServicesPetUpdate](ServicesPetUpdate.md) | Данные для обновления | |

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
| **404** | Питомец не найден |  -  |
| **500** | Внутренняя ошибка сервера |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## petsUserUserIdGet

> Array&lt;ModelsPet&gt; petsUserUserIdGet(userId)

Получение питомцев пользователя

Возвращает всех питомцев конкретного пользователя

### Example

```ts
import {
  Configuration,
  PetsApi,
} from '';
import type { PetsUserUserIdGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new PetsApi();

  const body = {
    // string | ID пользователя
    userId: userId_example,
  } satisfies PetsUserUserIdGetRequest;

  try {
    const data = await api.petsUserUserIdGet(body);
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
| **userId** | `string` | ID пользователя | [Defaults to `undefined`] |

### Return type

[**Array&lt;ModelsPet&gt;**](ModelsPet.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Список питомцев |  -  |
| **400** | Неверный запрос |  -  |
| **404** | Пользователь не найден |  -  |
| **500** | Внутренняя ошибка сервера |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## petsUserUserIdPost

> ModelsPet petsUserUserIdPost(userId, request)

Создание нового питомца

Создает нового питомца для пользователя

### Example

```ts
import {
  Configuration,
  PetsApi,
} from '';
import type { PetsUserUserIdPostRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new PetsApi();

  const body = {
    // string | ID пользователя
    userId: userId_example,
    // ServicesPetCreate | Данные питомца
    request: ...,
  } satisfies PetsUserUserIdPostRequest;

  try {
    const data = await api.petsUserUserIdPost(body);
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
| **userId** | `string` | ID пользователя | [Defaults to `undefined`] |
| **request** | [ServicesPetCreate](ServicesPetCreate.md) | Данные питомца | |

### Return type

[**ModelsPet**](ModelsPet.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **201** | Созданный питомец |  -  |
| **400** | Неверный запрос |  -  |
| **404** | Пользователь не найден |  -  |
| **500** | Внутренняя ошибка сервера |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

