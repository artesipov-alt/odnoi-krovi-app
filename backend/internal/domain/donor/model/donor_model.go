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
	ID               string
	RequestID        string
	DonorID          string
	DonorName        string
	DonorPhotos      []string
	DonorBloodGroup  string
	Amount           float64
	WarnFactors      []string
	CompensationType string
	TaxiCompensation bool
	IsConfirmed      bool
	RejectedReason   string
	Status           DonorResponseStatus
	CreatedAt        *time.Time
	UpdatedAt        *time.Time
	DeletedAt        *time.Time
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
		RequestID:        requestID,
		DonorID:          donorID,
		Amount:           math.Round(amount*10) / 10,
		CompensationType: compensationType,
		TaxiCompensation: taxiCompensation,
		Status:           DonorResponseStatusPending,
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
		return errors.New("reason is required")
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
	if d.Status != DonorResponseStatusAccepted {
		return errors.New("Невозможно подтвердить отклик. не верный первичный статус")
	}
	d.Status = DonorResponseStatusCompleted
	d.Amount = amount
	d.IsConfirmed = true
	return nil
}
