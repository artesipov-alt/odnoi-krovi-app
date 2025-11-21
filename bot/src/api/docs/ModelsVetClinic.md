
# ModelsVetClinic


## Properties

Name | Type
------------ | -------------
`appointmentRequirementId` | number
`clinicId` | number
`contactPersonName` | string
`contactPersonPosition` | string
`donorBonusPrograms` | string
`donorRequirements` | [ModelsDonorRequirements](ModelsDonorRequirements.md)
`latitude` | number
`locationId` | number
`longitude` | number
`name` | string
`phone` | string
`transfusionConditions` | string
`website` | string
`workHours` | string

## Example

```typescript
import type { ModelsVetClinic } from ''

// TODO: Update the object below with actual values
const example = {
  "appointmentRequirementId": 1,
  "clinicId": 1,
  "contactPersonName": Мария Петрова,
  "contactPersonPosition": Администратор,
  "donorBonusPrograms": Бонусные программы для доноров,
  "donorRequirements": null,
  "latitude": 55.7558,
  "locationId": 1,
  "longitude": 37.6173,
  "name": ВетКлиника ЗооДоктор,
  "phone": +79991234567,
  "transfusionConditions": Условия для переливания крови,
  "website": https://vetclinic.example.com,
  "workHours": Пн-Пт: 9:00-18:00,
} satisfies ModelsVetClinic

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ModelsVetClinic
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


