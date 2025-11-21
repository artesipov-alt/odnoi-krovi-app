
# ModelsBloodStock


## Properties

Name | Type
------------ | -------------
`bloodTypeId` | number
`clinicId` | number
`expirationDate` | string
`id` | number
`petType` | [ModelsPetType](ModelsPetType.md)
`priceRub` | number
`status` | [ModelsBloodStockStatus](ModelsBloodStockStatus.md)
`volumeMl` | number

## Example

```typescript
import type { ModelsBloodStock } from ''

// TODO: Update the object below with actual values
const example = {
  "bloodTypeId": 1,
  "clinicId": 1,
  "expirationDate": 2024-12-31,
  "id": 1,
  "petType": null,
  "priceRub": 5000.0,
  "status": null,
  "volumeMl": 500,
} satisfies ModelsBloodStock

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ModelsBloodStock
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


