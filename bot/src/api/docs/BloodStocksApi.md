# BloodStocksApi

All URIs are relative to */api/v1*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**bloodStocksBloodTypeBloodTypeIdGet**](BloodStocksApi.md#bloodstocksbloodtypebloodtypeidget) | **GET** /blood-stocks/blood-type/{blood_type_id} | Получение запасов крови по типу крови |
| [**bloodStocksClinicClinicIdGet**](BloodStocksApi.md#bloodstocksclinicclinicidget) | **GET** /blood-stocks/clinic/{clinic_id} | Получение запасов крови клиники |
| [**bloodStocksGet**](BloodStocksApi.md#bloodstocksget) | **GET** /blood-stocks | Получение всех запасов крови |
| [**bloodStocksIdDelete**](BloodStocksApi.md#bloodstocksiddelete) | **DELETE** /blood-stocks/{id} | Удаление запаса крови |
| [**bloodStocksIdGet**](BloodStocksApi.md#bloodstocksidget) | **GET** /blood-stocks/{id} | Получение запаса крови по ID |
| [**bloodStocksIdPut**](BloodStocksApi.md#bloodstocksidput) | **PUT** /blood-stocks/{id} | Обновление запаса крови |
| [**bloodStocksPost**](BloodStocksApi.md#bloodstockspost) | **POST** /blood-stocks | Создание нового запаса крови |
| [**bloodStocksSearchGet**](BloodStocksApi.md#bloodstockssearchget) | **GET** /blood-stocks/search | Поиск запасов крови с фильтрами |



## bloodStocksBloodTypeBloodTypeIdGet

> Array&lt;ModelsBloodStock&gt; bloodStocksBloodTypeBloodTypeIdGet(bloodTypeId)

Получение запасов крови по типу крови

Возвращает все запасы крови для конкретного типа крови

### Example

```ts
import {
  Configuration,
  BloodStocksApi,
} from '';
import type { BloodStocksBloodTypeBloodTypeIdGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new BloodStocksApi();

  const body = {
    // number | ID типа крови
    bloodTypeId: 56,
  } satisfies BloodStocksBloodTypeBloodTypeIdGetRequest;

  try {
    const data = await api.bloodStocksBloodTypeBloodTypeIdGet(body);
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
| **bloodTypeId** | `number` | ID типа крови | [Defaults to `undefined`] |

### Return type

[**Array&lt;ModelsBloodStock&gt;**](ModelsBloodStock.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Список запасов крови |  -  |
| **400** | Неверный запрос |  -  |
| **404** | Тип крови не найден |  -  |
| **500** | Внутренняя ошибка сервера |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## bloodStocksClinicClinicIdGet

> Array&lt;ModelsBloodStock&gt; bloodStocksClinicClinicIdGet(clinicId)

Получение запасов крови клиники

Возвращает все запасы крови для конкретной клиники

### Example

```ts
import {
  Configuration,
  BloodStocksApi,
} from '';
import type { BloodStocksClinicClinicIdGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new BloodStocksApi();

  const body = {
    // number | ID клиники
    clinicId: 56,
  } satisfies BloodStocksClinicClinicIdGetRequest;

  try {
    const data = await api.bloodStocksClinicClinicIdGet(body);
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
| **clinicId** | `number` | ID клиники | [Defaults to `undefined`] |

### Return type

[**Array&lt;ModelsBloodStock&gt;**](ModelsBloodStock.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Список запасов крови клиники |  -  |
| **400** | Неверный запрос |  -  |
| **404** | Клиника не найдена |  -  |
| **500** | Внутренняя ошибка сервера |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## bloodStocksGet

> Array&lt;ModelsBloodStock&gt; bloodStocksGet()

Получение всех запасов крови

Возвращает список всех запасов крови в системе

### Example

```ts
import {
  Configuration,
  BloodStocksApi,
} from '';
import type { BloodStocksGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new BloodStocksApi();

  try {
    const data = await api.bloodStocksGet();
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

[**Array&lt;ModelsBloodStock&gt;**](ModelsBloodStock.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Список запасов крови |  -  |
| **500** | Внутренняя ошибка сервера |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## bloodStocksIdDelete

> HandlersSuccessResponse bloodStocksIdDelete(id)

Удаление запаса крови

Удаляет запас крови из системы

### Example

```ts
import {
  Configuration,
  BloodStocksApi,
} from '';
import type { BloodStocksIdDeleteRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new BloodStocksApi();

  const body = {
    // number | ID запаса крови
    id: 56,
  } satisfies BloodStocksIdDeleteRequest;

  try {
    const data = await api.bloodStocksIdDelete(body);
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
| **id** | `number` | ID запаса крови | [Defaults to `undefined`] |

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
| **200** | Запас крови успешно удален |  -  |
| **400** | Неверный запрос |  -  |
| **404** | Запас крови не найден |  -  |
| **500** | Внутренняя ошибка сервера |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## bloodStocksIdGet

> ModelsBloodStock bloodStocksIdGet(id)

Получение запаса крови по ID

Возвращает информацию о конкретном запасе крови

### Example

```ts
import {
  Configuration,
  BloodStocksApi,
} from '';
import type { BloodStocksIdGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new BloodStocksApi();

  const body = {
    // number | ID запаса крови
    id: 56,
  } satisfies BloodStocksIdGetRequest;

  try {
    const data = await api.bloodStocksIdGet(body);
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
| **id** | `number` | ID запаса крови | [Defaults to `undefined`] |

### Return type

[**ModelsBloodStock**](ModelsBloodStock.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Запас крови |  -  |
| **400** | Неверный запрос |  -  |
| **404** | Запас крови не найден |  -  |
| **500** | Внутренняя ошибка сервера |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## bloodStocksIdPut

> HandlersSuccessResponse bloodStocksIdPut(id, request)

Обновление запаса крови

Обновляет информацию о запасе крови

### Example

```ts
import {
  Configuration,
  BloodStocksApi,
} from '';
import type { BloodStocksIdPutRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new BloodStocksApi();

  const body = {
    // number | ID запаса крови
    id: 56,
    // ServicesBloodStockUpdate | Данные для обновления
    request: ...,
  } satisfies BloodStocksIdPutRequest;

  try {
    const data = await api.bloodStocksIdPut(body);
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
| **id** | `number` | ID запаса крови | [Defaults to `undefined`] |
| **request** | [ServicesBloodStockUpdate](ServicesBloodStockUpdate.md) | Данные для обновления | |

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
| **200** | Запас крови успешно обновлен |  -  |
| **400** | Неверный запрос |  -  |
| **404** | Запас крови не найден |  -  |
| **500** | Внутренняя ошибка сервера |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## bloodStocksPost

> ModelsBloodStock bloodStocksPost(request)

Создание нового запаса крови

Создает новый запас крови в системе

### Example

```ts
import {
  Configuration,
  BloodStocksApi,
} from '';
import type { BloodStocksPostRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new BloodStocksApi();

  const body = {
    // ServicesBloodStockCreate | Данные запаса крови
    request: ...,
  } satisfies BloodStocksPostRequest;

  try {
    const data = await api.bloodStocksPost(body);
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
| **request** | [ServicesBloodStockCreate](ServicesBloodStockCreate.md) | Данные запаса крови | |

### Return type

[**ModelsBloodStock**](ModelsBloodStock.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: `application/json`
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **201** | Созданный запас крови |  -  |
| **400** | Неверный запрос |  -  |
| **404** | Клиника или тип крови не найдены |  -  |
| **500** | Внутренняя ошибка сервера |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## bloodStocksSearchGet

> Array&lt;ModelsBloodStock&gt; bloodStocksSearchGet(clinicId, petType, bloodTypeId, status, minVolume, maxVolume, minPrice, maxPrice)

Поиск запасов крови с фильтрами

Выполняет поиск запасов крови по различным параметрам (клиника, тип животного, тип крови, статус, объем, цена)

### Example

```ts
import {
  Configuration,
  BloodStocksApi,
} from '';
import type { BloodStocksSearchGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new BloodStocksApi();

  const body = {
    // number | ID клиники (optional)
    clinicId: 56,
    // string | Тип животного (dog/cat) (optional)
    petType: petType_example,
    // number | ID типа крови (optional)
    bloodTypeId: 56,
    // string | Статус (active/reserved/used/expired) (optional)
    status: status_example,
    // number | Минимальный объем (мл) (optional)
    minVolume: 56,
    // number | Максимальный объем (мл) (optional)
    maxVolume: 56,
    // number | Минимальная цена (руб) (optional)
    minPrice: 8.14,
    // number | Максимальная цена (руб) (optional)
    maxPrice: 8.14,
  } satisfies BloodStocksSearchGetRequest;

  try {
    const data = await api.bloodStocksSearchGet(body);
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
| **clinicId** | `number` | ID клиники | [Optional] [Defaults to `undefined`] |
| **petType** | `string` | Тип животного (dog/cat) | [Optional] [Defaults to `undefined`] |
| **bloodTypeId** | `number` | ID типа крови | [Optional] [Defaults to `undefined`] |
| **status** | `string` | Статус (active/reserved/used/expired) | [Optional] [Defaults to `undefined`] |
| **minVolume** | `number` | Минимальный объем (мл) | [Optional] [Defaults to `undefined`] |
| **maxVolume** | `number` | Максимальный объем (мл) | [Optional] [Defaults to `undefined`] |
| **minPrice** | `number` | Минимальная цена (руб) | [Optional] [Defaults to `undefined`] |
| **maxPrice** | `number` | Максимальная цена (руб) | [Optional] [Defaults to `undefined`] |

### Return type

[**Array&lt;ModelsBloodStock&gt;**](ModelsBloodStock.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Список найденных запасов крови |  -  |
| **400** | Неверный запрос |  -  |
| **500** | Внутренняя ошибка сервера |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

