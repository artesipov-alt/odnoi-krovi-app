
# ServicesPetCreate


## Properties

Name | Type
------------ | -------------
`ageMonths` | number
`ageYears` | number
`bloodGroup` | string
`breed` | string
`chipNumber` | string
`gender` | [ModelsGender](ModelsGender.md)
`hasChip` | boolean
`isGuideDog` | boolean
`isTherapist` | boolean
`knowsBloodGroup` | boolean
`latitude` | number
`livingCondition` | [ModelsLivingCondition](ModelsLivingCondition.md)
`longitude` | number
`name` | string
`photoUrl` | string
`sterilized` | boolean
`type` | [ModelsPetType](ModelsPetType.md)
`weightKg` | number

## Example

```typescript
import type { ServicesPetCreate } from ''

// TODO: Update the object below with actual values
const example = {
  "ageMonths": null,
  "ageYears": null,
  "bloodGroup": null,
  "breed": null,
  "chipNumber": null,
  "gender": null,
  "hasChip": null,
  "isGuideDog": null,
  "isTherapist": null,
  "knowsBloodGroup": null,
  "latitude": null,
  "livingCondition": null,
  "longitude": null,
  "name": null,
  "photoUrl": null,
  "sterilized": null,
  "type": null,
  "weightKg": null,
} satisfies ServicesPetCreate

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ServicesPetCreate
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


