// internal/storage/storage.go
package storage

import (
	"context"
	"io"
)

type FileStorage interface {
	Upload(ctx context.Context, file io.Reader, filename string, contentType string) (string, error)
	Download(ctx context.Context, filepath string) (io.ReadCloser, error)
	Delete(ctx context.Context, filepath string) error
	GetURL(filepath string) string
	GenerateAvatarURL(ctx context.Context, petID string) (string, error)
}
