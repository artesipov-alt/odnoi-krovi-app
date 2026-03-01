package bloodsearch

import "errors"

var (
	// ErrBloodRequestNotFound is returned when a blood request is not found
	ErrBloodRequestNotFound = errors.New("blood request not found")

	// ErrBloodRequestAlreadyExists is returned when trying to create a duplicate request
	ErrBloodRequestAlreadyExists = errors.New("blood request already exists for this pet")

	// ErrInvalidBloodRequestStatus is returned when an invalid status is provided
	ErrInvalidBloodRequestStatus = errors.New("invalid blood request status")

	// ErrDonorResponseNotFound is returned when a donor response is not found
	ErrDonorResponseNotFound = errors.New("donor response not found")

	// ErrDonorResponseAlreadyExists is returned when donor already responded
	ErrDonorResponseAlreadyExists = errors.New("donor already responded to this request")

	// ErrInsufficientVolume is returned when trying to reserve more blood than needed
	ErrInsufficientVolume = errors.New("insufficient blood volume available")

	// ErrRequestNotActive is returned when trying to modify a non-active request
	ErrRequestNotActive = errors.New("blood request is not active")
)
