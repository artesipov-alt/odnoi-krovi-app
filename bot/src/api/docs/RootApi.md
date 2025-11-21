# RootApi

All URIs are relative to */api/v1*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**rootGet**](RootApi.md#rootget) | **GET** / | Корневой эндпоинт |



## rootGet

> string rootGet()

Корневой эндпоинт

Возвращает информационное сообщение о сервере

### Example

```ts
import {
  Configuration,
  RootApi,
} from '';
import type { RootGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new RootApi();

  try {
    const data = await api.rootGet();
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

**string**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `text/html`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Информационное сообщение с HTML-ссылкой |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

