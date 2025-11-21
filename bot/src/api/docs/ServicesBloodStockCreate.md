
# ServicesBloodStockCreate


## Properties

Name | Type
------------ | -------------
`bloodTypeId` | number
`clinicId` | number
`expirationDate` | string
`petType` | [ModelsPetType](ModelsPetType.md)
`status` | [ModelsBloodStockStatus](ModelsBloodStockStatus.md)
`volumeMl` | number

## Example

```typescript
import type { ServicesBloodStockCreate } from ''

// TODO: Update the object below with actual values
const example = {
  "bloodTypeId": null,
  "clinicId": null,
  "expirationDate": null,
  "petType": null,
  "status": null,
  "volumeMl": null,
} satisfies ServicesBloodStockCreate

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ServicesBloodStockCreate
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


