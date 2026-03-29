package model

import (
	"errors"
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
	Amount           int32
	WarnFactors      []string
	CompensationType string
	TaxiCompensation bool
	IsConfirmed      bool
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
func NewDonorResponse(requestID, donorID, compensationType string, amount int32, taxiCompensation bool) (*DonorResponse, error) {
	if requestID == "" {
		return nil, errors.New("request ID is required")
	}
	if donorID == "" {
		return nil, errors.New("donor ID is required")
	}
	return &DonorResponse{
		RequestID:        requestID,
		DonorID:          donorID,
		Amount:           amount,
		CompensationType: compensationType,
		TaxiCompensation: taxiCompensation,
		Status:           DonorResponseStatusPending,
	}, nil
}
