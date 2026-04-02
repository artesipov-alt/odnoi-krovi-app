package model

import (
	"errors"
	"time"

	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
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
	BloodVolumeNeeded        int32
	BloodVolumeReserved      int32
	BloodVolumeDonated       int32
	Regions                  []string
	Description              string
	SmallPetsNotifyAllowed   bool
	Status                   BloodRequestStatus
	PhotoURLs                []string
	BloodGroupNames          []string
	BloodComponentIDs        []string
	OnBoarding               []string
	DonorApplications        []donormodel.DonorResponse
	PrioritySearch           bool
	IncludeUnknownBloodGroup bool
	CreatedAt                *time.Time
	UpdatedAt                *time.Time
	DeletedAt                *time.Time
}

// NewBloodRequest creates a new blood request with default values
func NewBloodRequest(petID string, bloodVolumeNeeded int32, regions []string) *BloodRequest {
	return &BloodRequest{
		PetID:                    petID,
		BloodVolumeNeeded:        bloodVolumeNeeded,
		BloodVolumeReserved:      0,
		Regions:                  regions,
		SmallPetsNotifyAllowed:   true,
		Status:                   BloodRequestStatusActive,
		PhotoURLs:                []string{},
		BloodGroupNames:          []string{},
		BloodComponentIDs:        []string{},
		OnBoarding:               []string{},
		PrioritySearch:           false,
		IncludeUnknownBloodGroup: false,
	}
}

// IsActive checks if the request is active
func (b *BloodRequest) IsActive() bool {
	return b.Status == BloodRequestStatusActive
}

// Close marks the request as closed
func (b *BloodRequest) Close() {
	b.Status = BloodRequestStatusClosed
}

// Activate marks the request as active
func (b *BloodRequest) activate() {
	b.Status = BloodRequestStatusActive
}

func (b *BloodRequest) MarkReservedFull() {
	b.Status = BloodRequestStatusReservedFull
}

// AddPhoto adds a photo URL to the request
func (b *BloodRequest) AddPhoto(url string) {
	b.PhotoURLs = append(b.PhotoURLs, url)
}

// SetBloodGroups sets compatible blood groups
func (b *BloodRequest) SetBloodGroups(groups []string) {
	b.BloodGroupNames = groups
}

func (b *BloodRequest) RecalculateBloodAmount() {
	var donated int32
	for _, app := range b.DonorApplications {
		if app.IsConfirmed && app.Status == donormodel.DonorResponseStatusCompleted {
			donated += app.Amount
		}
	}
	b.BloodVolumeDonated = donated

	var reserved int32
	for _, app := range b.DonorApplications {
		// Ищем только откликнувшихся доноров
		if !app.IsConfirmed && app.Status == donormodel.DonorResponseStatusAccepted {
			reserved += app.Amount + b.BloodVolumeDonated
			// Обрезаем до максимального
			if reserved >= b.BloodVolumeNeeded {
				reserved = b.BloodVolumeNeeded
			}
		}
	}
	b.BloodVolumeReserved = reserved
}

func (b *BloodRequest) RecalculateStatus() {
	if b.BloodVolumeReserved == b.BloodVolumeNeeded {
		b.MarkReservedFull()
	}
	if b.BloodVolumeDonated == b.BloodVolumeNeeded {
		b.Close()
	}
}

// BloodRequestFilter represents filter options for listing requests
type BloodRequestFilter struct {
	PetID   string
	Status  BloodRequestStatus
	Regions []string
	Limit   int
	Offset  int
}
