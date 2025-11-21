
# HandlersSimpleRegistrationRequest


## Properties

Name | Type
------------ | -------------
`fullName` | string
`telegramId` | number

## Example

```typescript
import type { HandlersSimpleRegistrationRequest } from ''

// TODO: Update the object below with actual values
const example = {
  "fullName": Иван Иванов,
  "telegramId": 123456789,
} satisfies HandlersSimpleRegistrationRequest

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as HandlersSimpleRegistrationRequest
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


