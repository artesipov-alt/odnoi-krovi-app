
# ModelsPet


## Properties

Name | Type
------------ | -------------
`ageMonths` | number
`ageYears` | number
`bloodGroup` | string
`breed` | string
`chipNumber` | string
`dewormingDate` | string
`ectoparasiteDate` | string
`gender` | [ModelsGender](ModelsGender.md)
`hasChip` | boolean
`id` | number
`isGuideDog` | boolean
`isTherapist` | boolean
`knowsBloodGroup` | boolean
`lastTransfusionDate` | string
`latitude` | number
`livingCondition` | [ModelsLivingCondition](ModelsLivingCondition.md)
`longitude` | number
`name` | string
`ownerId` | string
`photoUrl` | string
`sterilized` | boolean
`type` | [ModelsPetType](ModelsPetType.md)
`vaccinationDate` | string
`weightKg` | number

## Example

```typescript
import type { ModelsPet } from ''

// TODO: Update the object below with actual values
const example = {
  "ageMonths": 6,
  "ageYears": 3,
  "bloodGroup": DEA 1.1,
  "breed": Лабрадор,
  "chipNumber": 123456789,
  "dewormingDate": 2023-01-01T12:00:00Z,
  "ectoparasiteDate": 2023-01-01T12:00:00Z,
  "gender": null,
  "hasChip": false,
  "id": 1,
  "isGuideDog": false,
  "isTherapist": false,
  "knowsBloodGroup": false,
  "lastTransfusionDate": 2023-01-01T12:00:00Z,
  "latitude": 55.7558,
  "livingCondition": null,
  "longitude": 37.6173,
  "name": Бобик,
  "ownerId": 1,
  "photoUrl": https://example.com/photo.jpg,
  "sterilized": false,
  "type": null,
  "vaccinationDate": 2023-01-01T12:00:00Z,
  "weightKg": 25.5,
} satisfies ModelsPet

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ModelsPet
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


