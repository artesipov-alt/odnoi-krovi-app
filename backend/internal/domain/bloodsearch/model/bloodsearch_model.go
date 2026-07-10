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

// BloodRequestFilter represents filter options for listing requests
type BloodRequestFilter struct {
	PetID   string
	Status  BloodRequestStatus
	Regions []string
	Limit   int
	Offset  int
}
