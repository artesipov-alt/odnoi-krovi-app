
# ModelsDonorRequirements


## Properties

Name | Type
------------ | -------------
`bloodTypes` | Array&lt;string&gt;
`healthConditions` | Array&lt;string&gt;
`maxAge` | number
`minAge` | number
`minWeight` | number
`vaccinations` | Array&lt;string&gt;

## Example

```typescript
import type { ModelsDonorRequirements } from ''

// TODO: Update the object below with actual values
const example = {
  "bloodTypes": null,
  "healthConditions": null,
  "maxAge": null,
  "minAge": null,
  "minWeight": null,
  "vaccinations": null,
} satisfies ModelsDonorRequirements

console.log(example)

// Convert the instance to a JSON string
const exampleJSON: string = JSON.stringify(example)
console.log(exampleJSON)

// Parse the JSON string back to an object
const exampleParsed = JSON.parse(exampleJSON) as ModelsDonorRequirements
console.log(exampleParsed)
```

[[Back to top]](#) [[Back to API list]](../README.md#api-endpoints) [[Back to Model list]](../README.md#models) [[Back to README]](../README.md)


