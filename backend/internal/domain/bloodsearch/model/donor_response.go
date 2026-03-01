package model

import "time"

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
	ID         string
	RequestID  string
	DonorID    string
	Conditions []string
	Status     DonorResponseStatus
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
