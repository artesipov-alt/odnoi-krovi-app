package models

import "time"

// BloodSearch представляет запрос на поиск донора крови в системе
type BloodSearch struct {
	ID          int               `gorm:"primaryKey;autoIncrement" json:"id" example:"1"`
	PetID       *int              `gorm:"column:pet_id" json:"pet_id,omitempty" example:"1"`
	CreatedAt   time.Time         `json:"created_at" example:"2024-01-01T00:00:00Z"`
	LocationID  *int              `gorm:"column:location_id" json:"location_id,omitempty" example:"1"`
	BloodTypeID *int              `gorm:"column:blood_type_id" json:"blood_type_id,omitempty" example:"1"`
	Status      BloodSearchStatus `gorm:"column:status" json:"status,omitempty" example:"active"`
}
