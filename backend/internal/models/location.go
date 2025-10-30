package models

// Location представляет локацию в системе
type Location struct {
	ID   int    `gorm:"primaryKey;autoIncrement" json:"id" example:"1"`
	Name string `gorm:"size:255;not null" json:"name" example:"Москва"`
}
