// internal/storage/storage.go
package storage

import (
	"context"
	"io"
	"time"
)

type FileStorage interface {
	Upload(ctx context.Context, file io.Reader, filename string, contentType string) (string, error)
	UploadPublic(ctx context.Context, file io.Reader, filename string, contentType string) (string, error)
	Download(ctx context.Context, filepath string) (io.ReadCloser, error)
	Delete(ctx context.Context, filepath string) error
	GetURL(filepath string) string
	GetPresignedURL(ctx context.Context, filepath string, expire time.Duration) (string, error)
	GetAvatarUploadInfo(ctx context.Context, id string) (string, string, error)
	MakeAvatarPublic(ctx context.Context, id string) error
	CheckObjectExists(ctx context.Context, objectPath string) (bool, error)
	CheckAvatarExists(ctx context.Context, id string) (bool, error)
	GetAvatarPublicURL(id string) string
	UploadAvatar(ctx context.Context, id string, file io.Reader, makePublic bool) (string, error)
	SetObjectPublic(ctx context.Context, filepath string) error
	SetObjectPrivate(ctx context.Context, filepath string) error
}
