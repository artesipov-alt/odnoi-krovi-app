package filestorage

import (
	"context"
	"time"
)

type Repository interface {
	// Write methods

	// делает объект публично доступным
	SetObjectPublicACL(ctx context.Context, objectPath string) error

	// подтверждает завершение загрузки файлов
	ConfirmUploads(ctx context.Context, paths []string) error

	// Read methods

	// возвращает подписанные URL для загрузки
	GetPresignedURLs(ctx context.Context, count int64, id string) ([]UploadInfo, error)

	// проверяет, существует ли объект по пути
	CheckObjectExists(ctx context.Context, objectPath string) (bool, error)

	// возвращает публичный URL аватара по ID
	GetAvatarPublicURL(id string) string

	// возвращает публичный URL из пути к объекту
	GetPublicURLFromPath(path string) string

	// формирует массив URL фотографий из путей с учетом времени обновления
	BuildPhotoURLs(paths []string, updatedAt time.Time) []string
}

type MediaService interface {
	Compress(ctx context.Context, path string) error
}

// UploadInfo содержит информацию для загрузки файла: подписанную ссылку и путь в S3
type UploadInfo struct {
	UploadURL  string
	ObjectPath string
}
