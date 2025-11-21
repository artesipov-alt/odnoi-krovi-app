
# ServicesUserUpdate


## Properties

Name | Type
------------ | -------------
`allowGeo` | boolean
`email` | string
`fullName` | string
`locationId` | number
`onBoarding` | boolean
`phone` | string

## Example

```typescript
import type { ServicesUserUpdate } from ''

// TODO: Update the object below with actual values
const example = {
  "allowGeo": null,
  "email": null,
  "fullName": null,
  "locationId": null,
  "onBoarding": null,
  "phone": null,
} satisfies ServicesUserUpdate

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ServicesUserUpdate
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


