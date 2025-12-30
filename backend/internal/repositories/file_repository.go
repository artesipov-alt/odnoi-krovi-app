package repositories

import (
	"context"
)

type FileStorage interface {
	GetAvatarUploadInfo(ctx context.Context, id string) (string, string, error)
	CheckObjectExists(ctx context.Context, objectPath string) (bool, error)
	GetAvatarPublicURL(id string) string
	GetPublicURLFromPath(path string) string
	SetObjectPublicACL(ctx context.Context, objectPath string) error
}
