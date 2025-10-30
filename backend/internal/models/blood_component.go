package models

// BloodComponent представляет компонент крови в системе
type BloodComponent struct {
	ID   int    `gorm:"primaryKey;autoIncrement" json:"id" example:"1"`
	Name string `gorm:"type:varchar(255);not null" json:"name" example:"Плазма"`
}
