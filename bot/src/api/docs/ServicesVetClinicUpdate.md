
# ServicesVetClinicUpdate


## Properties

Name | Type
------------ | -------------
`appointmentRequirementId` | number
`contactPersonName` | string
`contactPersonPosition` | string
`donorBonusPrograms` | string
`locationId` | number
`name` | string
`phone` | string
`transfusionConditions` | string
`website` | string
`workHours` | string

## Example

```typescript
import type { ServicesVetClinicUpdate } from ''

// TODO: Update the object below with actual values
const example = {
  "appointmentRequirementId": null,
  "contactPersonName": null,
  "contactPersonPosition": null,
  "donorBonusPrograms": null,
  "locationId": null,
  "name": null,
  "phone": null,
  "transfusionConditions": null,
  "website": null,
  "workHours": null,
} satisfies ServicesVetClinicUpdate

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ServicesVetClinicUpdate
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


