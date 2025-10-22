package models

// VetClinic represents a veterinary clinic in the system
type VetClinic struct {
	UserID                   int     `gorm:"not null" json:"user_id" example:"1"`
	Name                     string  `gorm:"size:255;not null" json:"name" example:"ВетКлиника ЗооДоктор"`
	Phone                    string  `gorm:"size:20" json:"phone,omitempty" example:"+79991234567"`
	Website                  string  `gorm:"size:255" json:"website,omitempty" example:"https://vetclinic.example.com"`
	WorkHours                string  `gorm:"type:text" json:"work_hours,omitempty" example:"Пн-Пт: 9:00-18:00"`
	Latitude                 float64 `gorm:"type:numeric" json:"latitude,omitempty" example:"55.7558"`
	Longitude                float64 `gorm:"type:numeric" json:"longitude,omitempty" example:"37.6173"`
	TransfusionConditions    string  `gorm:"type:text" json:"transfusion_conditions,omitempty" example:"Условия для переливания крови"`
	DonorBonusPrograms       string  `gorm:"type:text" json:"donor_bonus_programs,omitempty" example:"Бонусные программы для доноров"`
	DonorRequirements        string  `gorm:"type:jsonb" json:"donor_requirements,omitempty" example:"{\"min_weight\": 20, \"age_min\": 1, \"age_max\": 8}"`
	ContactPersonName        string  `gorm:"size:255" json:"contact_person_name,omitempty" example:"Мария Петрова"`
	ContactPersonPosition    string  `gorm:"size:255" json:"contact_person_position,omitempty" example:"Администратор"`
	LocationID               int     `gorm:"not null;default:1" json:"location_id" example:"1"`
	AppointmentRequirementID int     `gorm:"not null;default:1" json:"appointment_requirement_id" example:"1"`
}

// // TableName specifies the table name for VetClinic
// func (VetClinic) TableName() string {
// 	return "vet_clinics"
// }
