package model

import (
	"errors"
	"math"
	"time"
)

// DonorResponseStatus представляет статус отклика донора
type DonorResponseStatus string

const (
	DonorResponseStatusPending   DonorResponseStatus = "pending"
	DonorResponseStatusAccepted  DonorResponseStatus = "accepted"
	DonorResponseStatusCompleted DonorResponseStatus = "completed"
	DonorResponseStatusRejected  DonorResponseStatus = "rejected"
	DonorResponseStatusCancelled DonorResponseStatus = "cancelled"
	DonorResponseStatusFailed    DonorResponseStatus = "failed"
)

// DonorResponse представляет отклик донора на заявку поиска крови
type DonorResponse struct {
	ID              string
	RequestID       string
	DonorID         string
	DonorName       string
	DonorPhotos     []string
	DonorBloodGroup string
	Amount          float64
	WarnFactors     []string
	DonorPrefs      DonorPrefs
	IsConfirmed     bool
	RejectedReason  string
	Status          DonorResponseStatus
	CreatedAt       *time.Time
	UpdatedAt       *time.Time
	DeletedAt       *time.Time
}

type DonorPrefs struct {
	PreferredLocationIDs []string
	CompensationType     string
	TaxiCompensation     bool
}

// DonorPreloadFilter представляет параметры для предзагрузки связанных данных
type DonorPreloadFilter struct {
	Status string
	Limit  int
	Offset int
}

// NewDonorResponse creates a new donor response with validation
func NewDonorResponse(requestID, donorID, compensationType string, amount float64, taxiCompensation bool) (*DonorResponse, error) {
	if requestID == "" {
		return nil, errors.New("request ID is required")
	}
	if donorID == "" {
		return nil, errors.New("donor ID is required")
	}
	return &DonorResponse{
		RequestID: requestID,
		DonorID:   donorID,
		Amount:    math.Round(amount*10) / 10,
		DonorPrefs: DonorPrefs{
			CompensationType: compensationType,
			TaxiCompensation: taxiCompensation,
		},
		Status: DonorResponseStatusPending,
	}, nil
}

func (d *DonorResponse) Accept() error {
	if d.Status == DonorResponseStatusPending || d.Status == DonorResponseStatusCompleted && d.IsConfirmed == false {
		d.Status = DonorResponseStatusAccepted
		return nil
	}
	return errors.New("Невозможно принять заявку, не верный первичный статус")
}

func (d *DonorResponse) Reject(reason string) error {
	if reason == "" {
		reason = "Реципиент отклонил донацию. "
	}

	switch d.Status {
	case DonorResponseStatusPending:
		d.Status = DonorResponseStatusRejected
		d.RejectedReason = reason
	case DonorResponseStatusAccepted:
		d.Status = DonorResponseStatusRejected
		d.RejectedReason = reason
	case DonorResponseStatusCompleted:
		// отклонили результат — возвращаем в работу
		d.Status = DonorResponseStatusAccepted
		d.RejectedReason = reason
	default:
		return errors.New("Невозможно отклонить заявку, не верный первичный статус")
	}

	return nil
}

func (d *DonorResponse) Complete(amount float64) error {
	if d.Status != DonorResponseStatusAccepted {
		return errors.New("Невозможно завершить отклик. не верный первичный статус")
	}
	d.Status = DonorResponseStatusCompleted
	d.Amount = amount
	return nil
}

func (d *DonorResponse) Confirm(amount float64) error {
	if d.Status != DonorResponseStatusAccepted && !(d.Status == DonorResponseStatusCompleted && d.IsConfirmed == false) {
		return errors.New("Невозможно подтвердить отклик. не верный первичный статус")
	}
	d.Status = DonorResponseStatusCompleted
	d.Amount = amount
	d.IsConfirmed = true
	return nil
}

func (d *DonorResponse) Cancel(reason string) error {
	if reason == "" {
		reason = "Донор самостоятельно отменил донацию. "
	}
	if d.Status == DonorResponseStatusPending || d.Status == DonorResponseStatusAccepted {
		d.RejectedReason = reason
		d.Status = DonorResponseStatusCancelled
		return nil
	}
	return errors.New("Невозможно отменить отклик. не верный первичный статус")
}

// IsActiveForDonation checks if the donor response is active for donation purposes.
// ВНИМАНИЕ: составной критерий (Accepted || Pending || (Completed && !IsConfirmed))
// используется очень широко и сознательно НЕ продублирован в SQL.
func (d *DonorResponse) IsActiveForDonation() bool {
	if d == nil {
		return false
	}
	return d.Status == DonorResponseStatusAccepted ||
		d.Status == DonorResponseStatusPending ||
		(d.Status == DonorResponseStatusCompleted && !d.IsConfirmed)
}

func (d *DonorResponse) IsInactive() bool {
	return d.Status == DonorResponseStatusRejected ||
		d.Status == DonorResponseStatusCancelled ||
		d.Status == DonorResponseStatusFailed ||
		d.IsFullyCompleted()
}

// IsFullyCompleted checks if donation is confirmed by recipient.
// ВНИМАНИЕ: критерий продублирован в SQL — EntDonorResponseRepository.CountFullyCompletedByOwnerID.
func (d *DonorResponse) IsFullyCompleted() bool {
	return d.Status == DonorResponseStatusCompleted && d.IsConfirmed
}

// IsCompleted checks if the donor response has completed status (regardless of confirmation)
func (d *DonorResponse) IsCompleted() bool {
	return d.Status == DonorResponseStatusCompleted
}

// IsAwaitingConfirmation checks if the donor has responded (accepted or completed)
// but the recipient hasn't confirmed yet.
func (d *DonorResponse) IsAwaitingConfirmation() bool {
	if d == nil {
		return false
	}
	return !d.IsConfirmed &&
		(d.Status == DonorResponseStatusAccepted ||
			d.Status == DonorResponseStatusCompleted)
}

// ApplyOwnerPrefs заполняет DonorPrefs (регионы, компенсация, такси) из предпочтений владельца.
// Используется, когда DonorResponse загружен из БД, где эти поля не хранятся.
func (d *DonorResponse) ApplyOwnerPrefs(locationIDs []string, compensationType string, taxiCompensation bool) {
	if d == nil {
		return
	}
	if locationIDs == nil {
		d.DonorPrefs.PreferredLocationIDs = []string{}
	} else {
		d.DonorPrefs.PreferredLocationIDs = locationIDs
	}
	d.DonorPrefs.CompensationType = compensationType
	d.DonorPrefs.TaxiCompensation = taxiCompensation
}
