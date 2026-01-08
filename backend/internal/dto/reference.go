package dto

// ReferenceResponse represents a response for reference data
type ReferenceResponse struct {
	Data []ReferenceItem `json:"data"`
}

// ReferenceItem represents an item in reference data
type ReferenceItem struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

// ReferenceResponseDB represents a response for reference data from DB
type ReferenceResponseDB struct {
	Data []ReferenceItemDB `json:"data"`
}

// ReferenceItemDB represents an item in reference data from DB
type ReferenceItemDB struct {
	Value int    `json:"value"`
	Label string `json:"label"`
}
