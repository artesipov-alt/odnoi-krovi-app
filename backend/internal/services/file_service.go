// internal/services/file_service.go
package services

import (
	"context"
	"io"
)

type FileService interface {
	UploadPetAvatar(ctx context.Context, petID string, file io.Reader) (string, error)
	ValidateImage(file io.Reader) error
	ResizeImage(file io.Reader, width, height int) (io.Reader, error)
}
