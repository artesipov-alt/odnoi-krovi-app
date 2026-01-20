package repositories

import (
	"context"
)

type FileStorage interface {
	GetPresignedURLs(ctx context.Context, count int64, id string) ([]UploadInfo, error)
	CheckObjectExists(ctx context.Context, objectPath string) (bool, error)
	GetAvatarPublicURL(id string) string
	GetPublicURLFromPath(path string) string
	SetObjectPublicACL(ctx context.Context, objectPath string) error
	ConfirmUploads(ctx context.Context, paths []string) error
}

// UploadInfo содержит информацию для загрузки файла: подписанную ссылку и путь в S3
type UploadInfo struct {
	UploadURL  string
	ObjectPath string
}
