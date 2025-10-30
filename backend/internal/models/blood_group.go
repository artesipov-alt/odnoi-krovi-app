package models

// BloodGroup представляет группу крови животного в системе
type BloodGroup struct {
	ID          int    `gorm:"primaryKey;autoIncrement" json:"id" example:"1"`
	PetType     string `gorm:"type:varchar(50);not null" json:"PetType" example:"dog"`
	BloodGroup  string `gorm:"type:varchar(50);not null" json:"bloodGroup" example:"DEA 1.1 Positive"`
	Description string `gorm:"type:text" json:"description,omitempty" example:"Универсальный донор для собак"`
}
