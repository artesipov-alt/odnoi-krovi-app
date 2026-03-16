package model

import (
	"errors"
	"time"
)

// DonorResponseStatus представляет статус отклика донора
type DonorResponseStatus string

const (
	DonorResponseStatusActive    DonorResponseStatus = "active"
	DonorResponseStatusAccepted  DonorResponseStatus = "accepted"
	DonorResponseStatusRejected  DonorResponseStatus = "rejected"
	DonorResponseStatusCancelled DonorResponseStatus = "cancelled"
)

// DonorResponse представляет отклик донора на заявку поиска крови
type DonorResponse struct {
	ID               string
	RequestID        string
	DonorID          string
	CompensationType string
	TaxiCompensation bool
	Status           DonorResponseStatus
	CreatedAt        *time.Time
	UpdatedAt        *time.Time
	DeletedAt        *time.Time
}

// NewDonorResponse creates a new donor response with validation
func NewDonorResponse(requestID, donorID, compensationType string, taxiCompensation bool) (*DonorResponse, error) {
	if requestID == "" {
		return nil, errors.New("request ID is required")
	}
	if donorID == "" {
		return nil, errors.New("donor ID is required")
	}
	return &DonorResponse{
		RequestID:        requestID,
		DonorID:          donorID,
		CompensationType: compensationType,
		TaxiCompensation: taxiCompensation,
		Status:           DonorResponseStatusActive,
	}, nil
}
