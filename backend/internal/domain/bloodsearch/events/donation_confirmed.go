package events

import "time"

type DonorInfo struct {
	ProviderMaxID string
	Name          string
}

type DonationConfirmed struct {
	Initiator string
	DonorData DonorInfo
	Volume    float64
	CreatedAt time.Time
}

func (e DonationConfirmed) EventName() string     { return "DonationConfirmed" }
func (e DonationConfirmed) OccurredAt() time.Time { return e.CreatedAt }
