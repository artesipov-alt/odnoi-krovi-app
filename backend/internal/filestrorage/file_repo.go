package filestrorage

import (
	"context"
	"time"
)

type Repository interface {
	GetPresignedURLs(ctx context.Context, count int64, id string) ([]UploadInfo, error)
	CheckObjectExists(ctx context.Context, objectPath string) (bool, error)
	GetAvatarPublicURL(id string) string
	GetPublicURLFromPath(path string) string
	BuildPhotoURLs(paths []string, updatedAt time.Time) []string
	SetObjectPublicACL(ctx context.Context, objectPath string) error
	ConfirmUploads(ctx context.Context, paths []string) error
}

type MediaService interface {
	Compress(ctx context.Context, path string) error
}

// UploadInfo содержит информацию для загрузки файла: подписанную ссылку и путь в S3
type UploadInfo struct {
	UploadURL  string
	ObjectPath string
}
