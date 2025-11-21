# VetClinicsApi

All URIs are relative to */api/v1*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**vetClinicsIdDelete**](VetClinicsApi.md#vetclinicsiddelete) | **DELETE** /vet-clinics/{id} | Удаление клиники по ID |
| [**vetClinicsIdGet**](VetClinicsApi.md#vetclinicsidget) | **GET** /vet-clinics/{id} | Получение профиля клиники по ID |
| [**vetClinicsIdPut**](VetClinicsApi.md#vetclinicsidput) | **PUT** /vet-clinics/{id} | Обновление профиля клиники |
| [**vetClinicsLocationLocationIdGet**](VetClinicsApi.md#vetclinicslocationlocationidget) | **GET** /vet-clinics/location/{location_id} | Получение всех клиник по ID локации |
| [**vetClinicsRegisterPost**](VetClinicsApi.md#vetclinicsregisterpost) | **POST** /vet-clinics/register | Регистрация новой ветеринарной клиники |



## vetClinicsIdDelete

> HandlersSuccessResponse vetClinicsIdDelete(id)

Удаление клиники по ID

Удаляет клинику из системы (soft delete)

### Example

```ts
import {
  Configuration,
  VetClinicsApi,
} from '';
import type { VetClinicsIdDeleteRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new VetClinicsApi();

  const body = {
    // number | ID клиники
    id: 56,
  } satisfies VetClinicsIdDeleteRequest;

  try {
    const data = await api.vetClinicsIdDelete(body);
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
| **id** | `number` | ID клиники | [Defaults to `undefined`] |

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
| **200** | Клиника успешно удалена |  -  |
| **400** | Неверный запрос |  -  |
| **404** | Клиника не найдена |  -  |
| **500** | Внутренняя ошибка сервера |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## vetClinicsIdGet

> ServicesVetClinicProfile vetClinicsIdGet(id)

Получение профиля клиники по ID

Возвращает полный профиль ветеринарной клиники

### Example

```ts
import {
  Configuration,
  VetClinicsApi,
} from '';
import type { VetClinicsIdGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new VetClinicsApi();

  const body = {
    // number | ID клиники
    id: 56,
  } satisfies VetClinicsIdGetRequest;

  try {
    const data = await api.vetClinicsIdGet(body);
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
| **id** | `number` | ID клиники | [Defaults to `undefined`] |

### Return type

[**ServicesVetClinicProfile**](ServicesVetClinicProfile.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Профиль клиники |  -  |
| **400** | Неверный запрос |  -  |
| **404** | Клиника не найдена |  -  |
| **500** | Внутренняя ошибка сервера |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## vetClinicsIdPut

> HandlersSuccessResponse vetClinicsIdPut(id, request)

Обновление профиля клиники

Обновляет информацию о ветеринарной клинике

### Example

```ts
import {
  Configuration,
  VetClinicsApi,
} from '';
import type { VetClinicsIdPutRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new VetClinicsApi();

  const body = {
    // number | ID клиники
    id: 56,
    // ServicesVetClinicUpdate | Данные для обновления
    request: ...,
  } satisfies VetClinicsIdPutRequest;

  try {
    const data = await api.vetClinicsIdPut(body);
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
| **id** | `number` | ID клиники | [Defaults to `undefined`] |
| **request** | [ServicesVetClinicUpdate](ServicesVetClinicUpdate.md) | Данные для обновления | |

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
| **404** | Клиника не найдена |  -  |
| **500** | Внутренняя ошибка сервера |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## vetClinicsLocationLocationIdGet

> Array&lt;ModelsVetClinic&gt; vetClinicsLocationLocationIdGet(locationId)

Получение всех клиник по ID локации

Возвращает список всех ветеринарных клиник в указанной локации

### Example

```ts
import {
  Configuration,
  VetClinicsApi,
} from '';
import type { VetClinicsLocationLocationIdGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new VetClinicsApi();

  const body = {
    // number | ID локации
    locationId: 56,
  } satisfies VetClinicsLocationLocationIdGetRequest;

  try {
    const data = await api.vetClinicsLocationLocationIdGet(body);
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
| **locationId** | `number` | ID локации | [Defaults to `undefined`] |

### Return type

[**Array&lt;ModelsVetClinic&gt;**](ModelsVetClinic.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Список клиник |  -  |
| **400** | Неверный запрос |  -  |
| **500** | Внутренняя ошибка сервера |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## vetClinicsRegisterPost

> ModelsVetClinic vetClinicsRegisterPost(request)

Регистрация новой ветеринарной клиники

Регистрирует новую ветеринарную клинику в системе

### Example

```ts
import {
  Configuration,
  VetClinicsApi,
} from '';
import type { VetClinicsRegisterPostRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new VetClinicsApi();

  const body = {
    // ServicesVetClinicRegistration | Данные клиники
    request: ...,
  } satisfies VetClinicsRegisterPostRequest;

  try {
    const data = await api.vetClinicsRegisterPost(body);
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
| **request** | [ServicesVetClinicRegistration](ServicesVetClinicRegistration.md) | Данные клиники | |

### Return type

[**ModelsVetClinic**](ModelsVetClinic.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **201** | Созданная клиника |  -  |
| **400** | Неверный запрос |  -  |
| **409** | Клиника уже существует |  -  |
| **500** | Внутренняя ошибка сервера |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

