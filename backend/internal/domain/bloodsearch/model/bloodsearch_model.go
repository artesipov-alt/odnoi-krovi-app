package model

import (
	"errors"
	"time"
)

// ErrInsufficientVolume is returned when trying to reserve more blood than needed
var ErrInsufficientVolume = errors.New("insufficient blood volume available")

// BloodRequestStatus represents the status of a blood search request
type BloodRequestStatus string

const (
	BloodRequestStatusActive       BloodRequestStatus = "active"
	BloodRequestStatusClosed       BloodRequestStatus = "closed"
	BloodRequestStatusDraft        BloodRequestStatus = "draft"
	BloodRequestStatusReservedFull BloodRequestStatus = "reserved_full"
)

// BloodRequest represents a request to search for blood donors
type BloodRequest struct {
	ID                       string
	PetID                    string
	BloodVolumeNeeded        float64
	BloodVolumeReserved      float64
	BloodVolumeDonated       float64
	Regions                  []string
	SmallPetsNotifyAllowed   bool
	Status                   BloodRequestStatus
	BloodGroupNames          []string
	BloodComponentIDs        []string
	OnBoarding               []string
	PrioritySearch           bool
	IncludeUnknownBloodGroup bool
	AdvancedInfo             AdvancedInfo
	CreatedAt                *time.Time
	UpdatedAt                *time.Time
	DeletedAt                *time.Time
}

type AdvancedInfo struct {
	Description string
	PhotoURLs   []string
}

// IsActive checks if the request is active
func (b *BloodRequest) IsActive() bool {
	return b.Status == BloodRequestStatusActive
}

// IsClosed checks if the request is closed
func (b *BloodRequest) IsClosed() bool {
	if b == nil {
		return false
	}
	return b.Status == BloodRequestStatusClosed
}

// IsReservedFull checks if the request is fully reserved —
// нужный объём крови уже обеспечен откликнувшимися донорами.
func (b *BloodRequest) IsReservedFull() bool {
	if b == nil {
		return false
	}
	return b.Status == BloodRequestStatusReservedFull
}

// Close marks the request as closed
func (b *BloodRequest) Close() {
	b.Status = BloodRequestStatusClosed
}

// Activate marks the request as active
func (b *BloodRequest) Activate() {
	b.Status = BloodRequestStatusActive
}

func (b *BloodRequest) MarkReservedFull() {
	b.Status = BloodRequestStatusReservedFull
}

// AddPhoto adds a photo URL to the request
func (b *BloodRequest) AddPhoto(url string) {
	b.AdvancedInfo.PhotoURLs = append(b.AdvancedInfo.PhotoURLs, url)
}

// SetBloodGroups sets compatible blood groups
func (b *BloodRequest) SetBloodGroups(groups []string) {
	b.BloodGroupNames = groups
}

func (b *BloodRequest) RecalculateStatus() {
	if b.Status == BloodRequestStatusClosed {
		return
	}
	if b.BloodVolumeReserved >= b.BloodVolumeNeeded {
		b.MarkReservedFull()
	} else {
		b.Activate()
	}
	if b.BloodVolumeDonated >= b.BloodVolumeNeeded {
		b.Close()
	}
}

func (b *BloodRequest) SetBloodVolume(donated float64, reserved float64) {
	b.BloodVolumeDonated = donated
	b.BloodVolumeReserved = reserved
}

// IsCoversNededAmount (Количество-крови) Бизнес-логика, донор должен покрывать весь объем или хотя бы половину от остатка если реципиент разрешил
func (b *BloodRequest) IsCoversNededAmount(avilableDonorAmount float64) bool {
	halfVolume := (b.BloodVolumeNeeded - b.BloodVolumeReserved) / 2
	if b.BloodVolumeReserved+avilableDonorAmount >= b.BloodVolumeNeeded {
		return true
	} else if b.SmallPetsNotifyAllowed && avilableDonorAmount >= halfVolume {
		return true
	}
	return false
}

// BloodRequestFilter represents filter options for listing requests
type BloodRequestFilter struct {
	PetID   string
	Status  BloodRequestStatus
	Regions []string
	Limit   int
	Offset  int
}
