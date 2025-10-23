package models

// Breed представляет породу животного в системе
type Breed struct {
	ID   int     `gorm:"primaryKey;autoIncrement" json:"id" example:"1"`
	Name string  `gorm:"size:100" json:"name,omitempty" example:"Лабрадор"`
	Type PetType `gorm:"type:varchar(50)" json:"type,omitempty" example:"dog"`
}
