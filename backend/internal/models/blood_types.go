package models

// BloodType представляет тип крови в системе
type BloodType struct {
	ID   int    `gorm:"primaryKey;autoIncrement" json:"id" example:"1"`
	Name string `gorm:"type:varchar(255);not null" json:"name" example:"Плазма"`
}
