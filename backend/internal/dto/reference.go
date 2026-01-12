package dto

// ReferenceCodeItem represents an item of reference data defined in code (string value)
type ReferenceCodeItem struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// ReferenceCodeData represents a response for reference data defined in code
type ReferenceCodeData struct {
	Data []ReferenceCodeItem `json:"data"`
}

// ReferenceLocalItem represents an item of reference data (string value)
type ReferenceLocalItem struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// ReferenceLocalData represents a response for reference data
type ReferenceLocalData struct {
	Data []ReferenceLocalItem `json:"data"`
}

// ReferenceDBItem represents an item of reference data from the database (int value)
type ReferenceDBItem struct {
	Value int    `json:"value"`
	Label string `json:"label"`
}

// ReferenceDBData represents a response for reference data from the database
type ReferenceDBData struct {
	Data []ReferenceDBItem `json:"data"`
}

// ReferenceCodeResponse is a wrapper for Huma OpenAPI documentation for a list of string-valued reference items defined in code.
type ReferenceCodeResponse struct {
	Body ReferenceCodeData
}

// ReferenceLocalResponse is a wrapper for Huma OpenAPI documentation for a list of string-valued reference items.
type ReferenceLocalResponse struct {
	Body ReferenceLocalData
}

// ReferenceDBResponse is a wrapper for Huma OpenAPI documentation for a list of int-valued reference items from the database.
type ReferenceDBResponse struct {
	Body ReferenceDBData
}

// PetTypePath defines the path parameter for pet type
type PetTypePath struct {
	PetType string `path:"pet_type" doc:"Тип животного" example:"dog"`
}

// PetTypeQuery defines the query parameter for pet type
type PetTypeQuery struct {
	PetType string `query:"petType" doc:"Тип животного (dog, cat, etc.)" example:"dog"`
}
