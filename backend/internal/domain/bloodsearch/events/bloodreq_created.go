// internal/domain/events/blood_request.go
package events

import "time"

type BloodRequestCreated struct {
	RequestID      string
	BloodTypes     []string
	Regions        []string
	AvilableDonors []Peers
	CreatedAt      time.Time
}

type Peers struct {
	TelegramID int32
	MaxID      int32
}

func (e BloodRequestCreated) EventName() string     { return "BloodRequestCreated" }
func (e BloodRequestCreated) OccurredAt() time.Time { return e.CreatedAt }
