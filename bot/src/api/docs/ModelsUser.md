
# ModelsUser


## Properties

Name | Type
------------ | -------------
`allowGeo` | boolean
`consentPd` | boolean
`createdAt` | string
`email` | string
`fullName` | string
`id` | string
`locationId` | number
`onBoarding` | boolean
`organizationName` | string
`phone` | string
`role` | [ModelsUserRole](ModelsUserRole.md)
`telegramId` | number

## Example

```typescript
import type { ModelsUser } from ''

// TODO: Update the object below with actual values
const example = {
  "allowGeo": true,
  "consentPd": true,
  "createdAt": 2023-01-01T12:00:00Z,
  "email": user@example.com,
  "fullName": Иван Иванов,
  "id": USR-25-0001,
  "locationId": 1,
  "onBoarding": false,
  "organizationName": ООО Ромашка,
  "phone": +79991234567,
  "role": null,
  "telegramId": 123456789,
} satisfies ModelsUser

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ModelsUser
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


