package models

import "gorm.io/gorm"

// VetClinic represents a veterinary clinic in the system
type VetClinic struct {
	ClinicID                 int                `gorm:"primaryKey;autoIncrement" json:"clinicId" example:"1"`
	Name                     string             `gorm:"size:255;not null" json:"name" example:"ВетКлиника ЗооДоктор"`
	Phone                    string             `gorm:"size:20" json:"phone,omitempty" example:"+79991234567"`
	Website                  string             `gorm:"size:255" json:"website,omitempty" example:"https://vetclinic.example.com"`
	WorkHours                string             `gorm:"type:text" json:"workHours,omitempty" example:"Пн-Пт: 9:00-18:00"`
	Latitude                 float64            `gorm:"type:numeric" json:"latitude,omitempty" example:"55.7558"`
	Longitude                float64            `gorm:"type:numeric" json:"longitude,omitempty" example:"37.6173"`
	TransfusionConditions    string             `gorm:"type:text" json:"transfusionConditions,omitempty" example:"Условия для переливания крови"`
	DonorBonusPrograms       string             `gorm:"type:text" json:"donorBonusPrograms,omitempty" example:"Бонусные программы для доноров"`
	DonorRequirements        *DonorRequirements `gorm:"type:jsonb;serializer:json" json:"donorRequirements,omitempty"`
	ContactPersonName        string             `gorm:"size:255" json:"contactPersonName,omitempty" example:"Мария Петрова"`
	ContactPersonPosition    string             `gorm:"size:255" json:"contactPersonPosition,omitempty" example:"Администратор"`
	LocationID               int                `gorm:"not null;default:1" json:"locationId" example:"1"`
	AppointmentRequirementID int                `gorm:"not null;default:1" json:"appointmentRequirementId" example:"1"`
	DeletedAt                *gorm.DeletedAt    `gorm:"index" json:"deletedAt,omitempty" swaggerignore:"true"`
}

// // TableName specifies the table name for VetClinic
// func (VetClinic) TableName() string {
// 	return "vet_clinics"
// }
