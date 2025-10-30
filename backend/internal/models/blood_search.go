package models

import "time"

// BloodSearch представляет запрос на поиск донора крови в системе
type BloodSearch struct {
	ID          int               `gorm:"primaryKey;autoIncrement" json:"id" example:"1"`
	PetID       *int              `gorm:"column:pet_id" json:"petId,omitempty" example:"1"`
	CreatedAt   time.Time         `json:"createdAt" example:"2024-01-01T00:00:00Z"`
	LocationID  *int              `gorm:"column:location_id" json:"locationId,omitempty" example:"1"`
	BloodTypeID *int              `gorm:"column:blood_type_id" json:"bloodTypeId,omitempty" example:"1"`
	Status      BloodSearchStatus `gorm:"column:status" json:"status,omitempty" example:"active"`
}
