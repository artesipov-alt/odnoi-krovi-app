package models

type Animal struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Species   string `json:"species"`
	BirthDate string `json:"birth_date"`
}
