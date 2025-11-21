# ReferenceApi

All URIs are relative to */api/v1*

| Method | HTTP request | Description |
|------------- | ------------- | -------------|
| [**referenceBloodComponentsGet**](ReferenceApi.md#referencebloodcomponentsget) | **GET** /reference/blood-components | Получение компонентов крови животных |
| [**referenceBloodGroupsPetTypeGet**](ReferenceApi.md#referencebloodgroupspettypeget) | **GET** /reference/blood-groups/{pet_type} | Получение групп крови животных по типу животного |
| [**referenceBloodSearchStatusesGet**](ReferenceApi.md#referencebloodsearchstatusesget) | **GET** /reference/blood-search-statuses | Получение всех статусов поиска крови |
| [**referenceBloodStockStatusesGet**](ReferenceApi.md#referencebloodstockstatusesget) | **GET** /reference/blood-stock-statuses | Получение всех статусов запаса крови |
| [**referenceBreedsByTypeGet**](ReferenceApi.md#referencebreedsbytypeget) | **GET** /reference/breeds-by-type | Получение пород животных по типу |
| [**referenceBreedsGet**](ReferenceApi.md#referencebreedsget) | **GET** /reference/breeds | Получение всех пород животных |
| [**referenceDonationStatusesGet**](ReferenceApi.md#referencedonationstatusesget) | **GET** /reference/donation-statuses | Получение всех статусов донорства |
| [**referenceGendersGet**](ReferenceApi.md#referencegendersget) | **GET** /reference/genders | Получение всех значений пола |
| [**referenceLivingConditionsGet**](ReferenceApi.md#referencelivingconditionsget) | **GET** /reference/living-conditions | Получение всех условий проживания |
| [**referenceLocationsGet**](ReferenceApi.md#referencelocationsget) | **GET** /reference/locations | Получение всех локаций |
| [**referencePetTypesGet**](ReferenceApi.md#referencepettypesget) | **GET** /reference/pet-types | Получение всех типов животных |
| [**referenceUserRolesGet**](ReferenceApi.md#referenceuserrolesget) | **GET** /reference/user-roles | Получение всех ролей пользователей |



## referenceBloodComponentsGet

> HandlersReferenceResponse referenceBloodComponentsGet()

Получение компонентов крови животных

Возвращает список компонентов крови животных для выбора на фронтенде

### Example

```ts
import {
  Configuration,
  ReferenceApi,
} from '';
import type { ReferenceBloodComponentsGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new ReferenceApi();

  try {
    const data = await api.referenceBloodComponentsGet();
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

[**HandlersReferenceResponse**](HandlersReferenceResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Список компонентов крови |  -  |
| **500** | Внутренняя ошибка сервера |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## referenceBloodGroupsPetTypeGet

> HandlersReferenceResponseDB referenceBloodGroupsPetTypeGet(petType)

Получение групп крови животных по типу животного

Возвращает список групп крови животных для выбора на фронтенде

### Example

```ts
import {
  Configuration,
  ReferenceApi,
} from '';
import type { ReferenceBloodGroupsPetTypeGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new ReferenceApi();

  const body = {
    // string | Тип животного
    petType: petType_example,
  } satisfies ReferenceBloodGroupsPetTypeGetRequest;

  try {
    const data = await api.referenceBloodGroupsPetTypeGet(body);
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
| **petType** | `string` | Тип животного | [Defaults to `undefined`] |

### Return type

[**HandlersReferenceResponseDB**](HandlersReferenceResponseDB.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Список групп крови |  -  |
| **400** | Неверный тип животного |  -  |
| **500** | Внутренняя ошибка сервера |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## referenceBloodSearchStatusesGet

> HandlersReferenceResponse referenceBloodSearchStatusesGet()

Получение всех статусов поиска крови

Возвращает все доступные статусы поиска крови для выбора на фронтенде

### Example

```ts
import {
  Configuration,
  ReferenceApi,
} from '';
import type { ReferenceBloodSearchStatusesGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new ReferenceApi();

  try {
    const data = await api.referenceBloodSearchStatusesGet();
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

[**HandlersReferenceResponse**](HandlersReferenceResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Список статусов поиска крови |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## referenceBloodStockStatusesGet

> HandlersReferenceResponse referenceBloodStockStatusesGet()

Получение всех статусов запаса крови

Возвращает все доступные статусы запаса крови для выбора на фронтенде

### Example

```ts
import {
  Configuration,
  ReferenceApi,
} from '';
import type { ReferenceBloodStockStatusesGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new ReferenceApi();

  try {
    const data = await api.referenceBloodStockStatusesGet();
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

[**HandlersReferenceResponse**](HandlersReferenceResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Список статусов запаса крови |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## referenceBreedsByTypeGet

> HandlersReferenceResponse referenceBreedsByTypeGet(petType)

Получение пород животных по типу

Возвращает список пород животных для указанного типа животного для выбора на фронтенде

### Example

```ts
import {
  Configuration,
  ReferenceApi,
} from '';
import type { ReferenceBreedsByTypeGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new ReferenceApi();

  const body = {
    // string | Тип животного (dog, cat, etc.)
    petType: petType_example,
  } satisfies ReferenceBreedsByTypeGetRequest;

  try {
    const data = await api.referenceBreedsByTypeGet(body);
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
| **petType** | `string` | Тип животного (dog, cat, etc.) | [Defaults to `undefined`] |

### Return type

[**HandlersReferenceResponse**](HandlersReferenceResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Список пород животных |  -  |
| **400** | Неверный тип животного |  -  |
| **500** | Внутренняя ошибка сервера |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## referenceBreedsGet

> HandlersReferenceResponse referenceBreedsGet()

Получение всех пород животных

Возвращает список всех пород животных в базе для выбора на фронтенде

### Example

```ts
import {
  Configuration,
  ReferenceApi,
} from '';
import type { ReferenceBreedsGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new ReferenceApi();

  try {
    const data = await api.referenceBreedsGet();
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

[**HandlersReferenceResponse**](HandlersReferenceResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Список пород животных |  -  |
| **500** | Внутренняя ошибка сервера |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## referenceDonationStatusesGet

> HandlersReferenceResponse referenceDonationStatusesGet()

Получение всех статусов донорства

Возвращает все доступные статусы донорства для выбора на фронтенде

### Example

```ts
import {
  Configuration,
  ReferenceApi,
} from '';
import type { ReferenceDonationStatusesGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new ReferenceApi();

  try {
    const data = await api.referenceDonationStatusesGet();
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

[**HandlersReferenceResponse**](HandlersReferenceResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Список статусов донорства |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## referenceGendersGet

> HandlersReferenceResponse referenceGendersGet()

Получение всех значений пола

Возвращает все доступные значения пола для выбора на фронтенде

### Example

```ts
import {
  Configuration,
  ReferenceApi,
} from '';
import type { ReferenceGendersGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new ReferenceApi();

  try {
    const data = await api.referenceGendersGet();
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

[**HandlersReferenceResponse**](HandlersReferenceResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Список значений пола |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## referenceLivingConditionsGet

> HandlersReferenceResponse referenceLivingConditionsGet()

Получение всех условий проживания

Возвращает все доступные условия проживания для выбора на фронтенде

### Example

```ts
import {
  Configuration,
  ReferenceApi,
} from '';
import type { ReferenceLivingConditionsGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new ReferenceApi();

  try {
    const data = await api.referenceLivingConditionsGet();
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

[**HandlersReferenceResponse**](HandlersReferenceResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Список условий проживания |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## referenceLocationsGet

> HandlersReferenceResponseDB referenceLocationsGet()

Получение всех локаций

Возвращает список всех локаций в системе для выбора на фронтенде

### Example

```ts
import {
  Configuration,
  ReferenceApi,
} from '';
import type { ReferenceLocationsGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new ReferenceApi();

  try {
    const data = await api.referenceLocationsGet();
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

[**HandlersReferenceResponseDB**](HandlersReferenceResponseDB.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Список локаций |  -  |
| **500** | Внутренняя ошибка сервера |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## referencePetTypesGet

> HandlersReferenceResponse referencePetTypesGet()

Получение всех типов животных

Возвращает все доступные типы животных для выбора на фронтенде

### Example

```ts
import {
  Configuration,
  ReferenceApi,
} from '';
import type { ReferencePetTypesGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new ReferenceApi();

  try {
    const data = await api.referencePetTypesGet();
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

[**HandlersReferenceResponse**](HandlersReferenceResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Список типов животных |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


## referenceUserRolesGet

> HandlersReferenceResponse referenceUserRolesGet()

Получение всех ролей пользователей

Возвращает все доступные роли пользователей для выбора на фронтенде

### Example

```ts
import {
  Configuration,
  ReferenceApi,
} from '';
import type { ReferenceUserRolesGetRequest } from '';

async function example() {
  console.log("🚀 Testing  SDK...");
  const api = new ReferenceApi();

  try {
    const data = await api.referenceUserRolesGet();
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

[**HandlersReferenceResponse**](HandlersReferenceResponse.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: `application/json`


### HTTP response details
| Status code | Description | Response headers |
|-------------|-------------|------------------|
| **200** | Список ролей пользователей |  -  |

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)

