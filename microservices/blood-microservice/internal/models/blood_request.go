package models

import (
	"encoding/json"
	"time"

	bloodrequestv1 "github.com/artesipov-alt/odnoi-krovi-app/microservices/blood-microservice/gen/api/bloodrequest/v1"
	"gorm.io/datatypes"
)

// BloodRequest represents a pet blood search request in the database
type BloodRequest struct {
	PetID                  string         `gorm:"primaryKey;column:pet_id;type:varchar(100);not null" json:"petId"`
	PetType                string         `gorm:"column:pet_type;type:varchar(100);not null" json:"petType"`
	BloodGroup             datatypes.JSON `gorm:"column:blood_group;type:jsonb" json:"bloodGroup"`
	BloodComponents        datatypes.JSON `gorm:"column:blood_components;type:jsonb" json:"bloodComponents"`
	BloodVolumeNeeded      int32          `gorm:"column:blood_volume_needed;type:integer" json:"bloodVolumeNeeded"`
	BloodVolumeReserved    int32          `gorm:"column:blood_volume_reserved;type:integer" json:"bloodVolumeReserved"`
	Regions                datatypes.JSON `gorm:"column:regions;type:jsonb" json:"regions"`
	SmallPetsNotifyAllowed bool           `gorm:"column:small_pets_notify_allowed;type:boolean;default:false" json:"smallPetsNotifyAllowed"`
	Description            string         `gorm:"column:description;type:text" json:"description"`
	Status                 string         `gorm:"column:status;type:varchar(50);default:'active'" json:"status"`
	CreatedAt              time.Time      `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt              time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

// ToProto converts BloodRequest model to proto message
func (m *BloodRequest) ToProto() *bloodrequestv1.BloodRequest {
	var bloodGroup []string
	var bloodComponents []int32
	var regions []int32

	// Parse JSON fields
	if m.BloodGroup.String() != "" {
		json.Unmarshal(m.BloodGroup, &bloodGroup)
	}
	if m.BloodComponents.String() != "" {
		json.Unmarshal(m.BloodComponents, &bloodComponents)
	}
	if m.Regions.String() != "" {
		json.Unmarshal(m.Regions, &regions)
	}

	return &bloodrequestv1.BloodRequest{
		PetId:                  m.PetID,
		PetType:                m.PetType,
		BloodGroup:             bloodGroup,
		BloodComponents:        bloodComponents,
		BloodVolumeNeeded:      m.BloodVolumeNeeded,
		BloodVolumeReserved:    m.BloodVolumeReserved,
		Regions:                regions,
		SmallPetsNotifyAllowed: m.SmallPetsNotifyAllowed,
		Description:            m.Description,
		Status:                 m.Status,
	}
}

// FromProto converts proto message to BloodRequest model
func (m *BloodRequest) FromProto(proto *bloodrequestv1.BloodRequest) error {
	m.PetID = proto.PetId
	m.PetType = proto.PetType
	m.BloodVolumeNeeded = proto.BloodVolumeNeeded
	m.BloodVolumeReserved = proto.BloodVolumeReserved
	m.SmallPetsNotifyAllowed = proto.SmallPetsNotifyAllowed
	m.Description = proto.Description
	m.Status = proto.Status

	// Convert slices to JSON
	if proto.BloodGroup != nil {
		bloodGroupJSON, err := json.Marshal(proto.BloodGroup)
		if err != nil {
			return err
		}
		m.BloodGroup = datatypes.JSON(bloodGroupJSON)
	}

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

// BloodRequestsToProto converts slice of BloodRequest models to proto messages
func BloodRequestsToProto(models []*BloodRequest) []*bloodrequestv1.BloodRequest {
	protoRows := make([]*bloodrequestv1.BloodRequest, len(models))
	for i, model := range models {
		protoRows[i] = model.ToProto()
	}
	return protoRows
}
