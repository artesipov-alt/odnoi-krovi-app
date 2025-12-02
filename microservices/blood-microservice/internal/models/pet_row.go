package models

import (
	"encoding/json"
	"time"

	bloodsearchv1 "github.com/artesipov-alt/odnoi-krovi-app/microservices/blood-microservice/gen/api/bloodsearch/v1"
	"gorm.io/datatypes"
)

// PetRow represents a pet blood search request in the database
type PetRow struct {
	PetID                  string         `gorm:"primaryKey;column:pet_id;type:varchar(255);not null" json:"pet_id"`
	PetType                string         `gorm:"column:pet_type;type:varchar(100);not null" json:"pet_type"`
	BloodGroup             string         `gorm:"column:blood_group;type:varchar(50);not null" json:"blood_group"`
	BloodComponents        datatypes.JSON `gorm:"column:blood_components;type:jsonb" json:"blood_components"`
	BloodVolume            float32        `gorm:"column:blood_volume;type:real" json:"blood_volume"`
	Regions                datatypes.JSON `gorm:"column:regions;type:jsonb" json:"regions"`
	SmallPetsNotifyAllowed bool           `gorm:"column:small_pets_notify_allowed;type:boolean;default:false" json:"small_pets_notify_allowed"`
	Status                 string         `gorm:"column:status;type:varchar(50);default:'active'" json:"status"`
	CreatedAt              time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt              time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

// ToProto converts PetRow model to proto message
func (m *PetRow) ToProto() *bloodsearchv1.PetRow {
	var bloodComponents []string
	var regions []int32

	// Parse JSON fields
	if m.BloodComponents.String() != "" {
		json.Unmarshal(m.BloodComponents, &bloodComponents)
	}
	if m.Regions.String() != "" {
		json.Unmarshal(m.Regions, &regions)
	}

	return &bloodsearchv1.PetRow{
		PetId:                  m.PetID,
		PetType:                m.PetType,
		BloodGroup:             m.BloodGroup,
		BloodComponents:        bloodComponents,
		BloodVolume:            m.BloodVolume,
		Regions:                regions,
		SmallPetsNotifyAllowed: m.SmallPetsNotifyAllowed,
		Status:                 m.Status,
	}
}

// FromProto converts proto message to PetRow model
func (m *PetRow) FromProto(proto *bloodsearchv1.PetRow) error {
	m.PetID = proto.PetId
	m.PetType = proto.PetType
	m.BloodGroup = proto.BloodGroup
	m.BloodVolume = proto.BloodVolume
	m.SmallPetsNotifyAllowed = proto.SmallPetsNotifyAllowed
	m.Status = proto.Status

	// Convert slices to JSON
	if proto.BloodComponents != nil {
		componentsJSON, err := json.Marshal(proto.BloodComponents)
		if err != nil {
			return err
		}
		m.BloodComponents = datatypes.JSON(componentsJSON)
	}

	if proto.Regions != nil {
		regionsJSON, err := json.Marshal(proto.Regions)
		if err != nil {
			return err
		}
		m.Regions = datatypes.JSON(regionsJSON)
	}

	return nil
}

// PetRowsToProto converts slice of PetRow models to proto messages
func PetRowsToProto(models []*PetRow) []*bloodsearchv1.PetRow {
	protoRows := make([]*bloodsearchv1.PetRow, len(models))
	for i, model := range models {
		protoRows[i] = model.ToProto()
	}
	return protoRows
}
