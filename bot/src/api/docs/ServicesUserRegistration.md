
# ServicesUserRegistration


## Properties

Name | Type
------------ | -------------
`consentPd` | boolean
`email` | string
`fullName` | string
`locationId` | number
`phone` | string
`role` | [ModelsUserRole](ModelsUserRole.md)

## Example

```typescript
import type { ServicesUserRegistration } from ''

// TODO: Update the object below with actual values
const example = {
  "consentPd": null,
  "email": null,
  "fullName": null,
  "locationId": null,
  "phone": null,
  "role": null,
} satisfies ServicesUserRegistration

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ServicesUserRegistration
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


