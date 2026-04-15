package model

import (
	"errors"
	"math"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/common"
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
	OwnerID                  string
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

type BloodRequestWithApplications struct {
	BloodRequest
	DonorApplications []donormodel.DonorResponse
}

// Recipient представляет модель чтения реципиент
type BloodRequestWithMatchingDonors struct {
	BloodRequest
	RecipientData     RecipientData
	MatchingDonors    []MatchingDonorReadModel
	DefaultDonorPrefs *DefaultDonorPrefs
}

type RecipientData struct {
	PetName        string
	PetType        common.PetType
	BloodGroupName string
	OwnerName      string
	PhotoURLs      []string
}

type AdvancedInfo struct {
	Description string
	PhotoURLs   []string
}

type DefaultDonorPrefs struct {
	CompensationType common.CompensationType
	Bonuses          []string
	TaxiCompensation bool
}

// MatchingDonorReadModel представляет модель чтения для подходящего донора
type MatchingDonorReadModel struct {
	PetID           string
	PetName         string
	Amount          float64
	DonorBloodGroup string
	PhotoURLs       []string
}

// NewBloodRequest creates a new blood request with default values
func NewBloodRequest(petID string, bloodVolumeNeeded float64, regions []string) *BloodRequest {
	return &BloodRequest{
		PetID:                    petID,
		BloodVolumeNeeded:        math.Round(bloodVolumeNeeded*10) / 10,
		BloodVolumeReserved:      0,
		Regions:                  regions,
		SmallPetsNotifyAllowed:   true,
		Status:                   BloodRequestStatusActive,
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

func (b *BloodRequestWithApplications) RecalculateBloodAmount() {
	var donated float64
	for _, app := range b.DonorApplications {
		if app.IsConfirmed && app.Status == donormodel.DonorResponseStatusCompleted {
			donated += app.Amount
		}
	}
	b.BloodVolumeDonated = math.Round(donated*10) / 10

	reserved := b.BloodVolumeDonated
	for _, app := range b.DonorApplications {
		// Ищем только откликнувшихся доноров
		if (app.IsConfirmed == false && app.Status == donormodel.DonorResponseStatusAccepted) || (app.IsConfirmed == false && app.Status == donormodel.DonorResponseStatusCompleted) {
			reserved += app.Amount
			// Обрезаем до максимального
			if reserved >= b.BloodVolumeNeeded {
				reserved = b.BloodVolumeNeeded
			}
		}
	}
	b.BloodVolumeReserved = math.Round(reserved*10) / 10
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

// BloodRequestFilter represents filter options for listing requests
type BloodRequestFilter struct {
	PetID   string
	Status  BloodRequestStatus
	Regions []string
	Limit   int
	Offset  int
}

func (r *BloodRequestWithMatchingDonors) SetDefaultPrefs(compensationType common.CompensationType, taxiCompensation bool) {
	r.DefaultDonorPrefs = &DefaultDonorPrefs{
		CompensationType: compensationType,
		TaxiCompensation: taxiCompensation,
		Bonuses:          []string{},
	}
}
