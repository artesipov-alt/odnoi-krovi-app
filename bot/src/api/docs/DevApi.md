# DevApi

All URIs are relative to */api/v1*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**devDeletedUsersGet**](DevApi.md#devdeletedusersget) | **GET** /dev/deleted-users | Получение всех удаленных пользователей |
| [**devResetUserIdPost**](DevApi.md#devresetuseridpost) | **POST** /dev/reset-user/{id} | Сброс пользователя к начальным настройкам |
| [**devRestoreUserIdPost**](DevApi.md#devrestoreuseridpost) | **POST** /dev/restore-user/{id} | Восстановление удаленного пользователя |



## devDeletedUsersGet

> HandlersGetDeletedUsersResponse devDeletedUsersGet()

Получение всех удаленных пользователей

Возвращает список всех мягко удаленных пользователей

### Example

```ts
import {
  Configuration,
  DevApi,
} from '';
import type { DevDeletedUsersGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new DevApi();

  try {
    const data = await api.devDeletedUsersGet();
    console.log(data);
  } catch (error) {
    console.error(error);
  }
}

// Run the test
example().catch(console.error);
```

### Parameters

This endpoint does not need any parameter.

### Return type

[**HandlersGetDeletedUsersResponse**](HandlersGetDeletedUsersResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Список удаленных пользователей |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## devResetUserIdPost

> HandlersDevResponse devResetUserIdPost(id)

Сброс пользователя к начальным настройкам

Сбрасывает пользователя к заводским настройкам на этапе команды старт от бота

### Example

```ts
import {
  Configuration,
  DevApi,
} from '';
import type { DevResetUserIdPostRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new DevApi();

  const body = {
    // string | ID пользователя
    id: id_example,
  } satisfies DevResetUserIdPostRequest;

  try {
    const data = await api.devResetUserIdPost(body);
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

[**HandlersDevResponse**](HandlersDevResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Успешный сброс пользователя |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## devRestoreUserIdPost

> HandlersDevResponse devRestoreUserIdPost(id)

Восстановление удаленного пользователя

Восстанавливает мягко удаленного пользователя, устанавливая deleted_at в NULL

### Example

```ts
import {
  Configuration,
  DevApi,
} from '';
import type { DevRestoreUserIdPostRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new DevApi();

  const body = {
    // string | ID пользователя
    id: id_example,
  } satisfies DevRestoreUserIdPostRequest;

  try {
    const data = await api.devRestoreUserIdPost(body);
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

[**HandlersDevResponse**](HandlersDevResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Успешное восстановление пользователя |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

