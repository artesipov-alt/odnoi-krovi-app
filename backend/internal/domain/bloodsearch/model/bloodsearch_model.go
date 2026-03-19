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
	BloodRequestStatusActive BloodRequestStatus = "active"
	BloodRequestStatusClosed BloodRequestStatus = "closed"
	BloodRequestStatusDraft  BloodRequestStatus = "draft"
)

// BloodRequest represents a request to search for blood donors
type BloodRequest struct {
	ID                       string
	PetID                    string
	BloodVolumeNeeded        int32
	BloodVolumeReserved      int32
	Regions                  []string
	SmallPetsNotifyAllowed   bool
	Status                   BloodRequestStatus
	Description              string
	PhotoURLs                []string
	BloodGroupNames          []string
	BloodComponentIDs        []string
	OnBoarding               []string
	PrioritySearch           bool
	IncludeUnknownBloodGroup bool
	CreatedAt                *time.Time
	UpdatedAt                *time.Time
	DeletedAt                *time.Time
	DonorApplications        []DonorApplication
}

// DonorApplication представляет отклик донора
type DonorApplication struct {
	ID               string
	RequestID        string
	DonorID          string
	DonorName        string
	DonorPhotos      []string
	DonorBloodGroup  string
	Amount           int32
	WarnFactors      []string
	CompensationType string
	TaxiCompensation bool
	Status           string
	CreatedAt        *time.Time
	UpdatedAt        *time.Time
	DeletedAt        *time.Time
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

// ReserveVolume reserves blood volume
func (b *BloodRequest) ReserveVolume(amount int32) error {
	b.BloodVolumeReserved += amount
	// Auto-close if fully reserved
	if b.BloodVolumeReserved >= b.BloodVolumeNeeded {
		b.Status = BloodRequestStatusClosed
	}
	return nil
}

// AddPhoto adds a photo URL to the request
func (b *BloodRequest) AddPhoto(url string) {
	b.PhotoURLs = append(b.PhotoURLs, url)
}

// SetBloodGroups sets compatible blood groups
func (b *BloodRequest) SetBloodGroups(groups []string) {
	b.BloodGroupNames = groups
}

// BloodRequestFilter represents filter options for listing requests
type BloodRequestFilter struct {
	PetID   string
	Status  BloodRequestStatus
	Regions []string
	Limit   int
	Offset  int
}
