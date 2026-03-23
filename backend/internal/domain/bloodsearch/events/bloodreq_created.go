// internal/domain/events/blood_request.go
package events

import "time"

type BloodRequestCreated struct {
	RequestID string
	PetType   string
	BloodType string
	CityID    int64
	CreatedAt time.Time
}

func (e BloodRequestCreated) EventName() string     { return "BloodRequestCreated" }
func (e BloodRequestCreated) OccurredAt() time.Time { return e.CreatedAt }
