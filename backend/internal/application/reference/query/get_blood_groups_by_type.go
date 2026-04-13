package query

import (
	"context"
	"fmt"
)

// BloodGroupItem represents a blood group item
type BloodGroupItem struct {
	ID         string
	BloodGroup string
}

type GetBloodGroupsByPetTypeHandler struct{}

func NewGetBloodGroupsByPetTypeHandler() *GetBloodGroupsByPetTypeHandler {
	return &GetBloodGroupsByPetTypeHandler{}
}

func (h *GetBloodGroupsByPetTypeHandler) Handle(ctx context.Context, petType string) ([]BloodGroupItem, error) {
	switch petType {
	case "dog":
		return []BloodGroupItem{
			{ID: "DEA1+", BloodGroup: "DEA 1+"},
			{ID: "DEA1-", BloodGroup: "DEA 1-"},
			{ID: "UNKNOWN", BloodGroup: "Неизвестная"},
		}, nil
	case "cat":
		return []BloodGroupItem{
			{ID: "A", BloodGroup: "A"},
			{ID: "B", BloodGroup: "B"},
			{ID: "AB", BloodGroup: "AB"},
			{ID: "UNKNOWN", BloodGroup: "Неизвестная"},
		}, nil
	default:
		return nil, fmt.Errorf("unsupported pet type: %s", petType)
	}
}
