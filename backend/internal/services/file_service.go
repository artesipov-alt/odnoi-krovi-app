// internal/services/file_service.go
package services

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/repositories"
)

type FileService interface {
	// UploadPetAvatar(ctx context.Context, petID string, file io.Reader) (string, error)
	// ValidateImage(file io.Reader) error
	// ResizeImage(file io.Reader, width, height int) (io.Reader, error)
	GetPresignURLs(ctx context.Context, ID string, count int64, preloads ...string) ([]repositories.UploadInfo, error)
}

// FileServiceImpl реализует FileService
type FileServiceImpl struct {
	PetRepo   repositories.PetRepository
	UserRepo  repositories.UserRepository
	BloodRepo repositories.BloodRequestRepository
	storage   repositories.FileStorage
}

// GetPresignURLs Возвращает ссылки для загрузки фотографий.
func (s *FileServiceImpl) GetPresignURLs(ctx context.Context, ID string, count int64, preloads ...string) ([]repositories.UploadInfo, error) {
	if len(preloads) == 0 {
		return nil, apperrors.BadRequest("preload type is required")
	}

	var uploadInfos []repositories.UploadInfo

	switch preloads[0] {
	case "pet_avatar":
		exists, err := s.PetRepo.ExistsByID(ctx, ID)
		if err != nil {
			return nil, apperrors.Internal(err, "failed to check pet existence")
		}
		if !exists {
			return nil, apperrors.ErrPetNotFound
		}
		uploadInfos, err = s.storage.GetPresignedURLs(ctx, count, ID)
		if err != nil {
			return nil, apperrors.Internal(err, "failed to get upload URLs for pet")
		}
	case "user_avatar":
		exists, err := s.UserRepo.ExistsByID(ctx, ID)
		if err != nil {
			return nil, apperrors.Internal(err, "failed to check user existence")
		}
		if !exists {
			return nil, apperrors.ErrUserNotFound
		}
		uploadInfos, err = s.storage.GetPresignedURLs(ctx, count, ID)
		if err != nil {
			return nil, apperrors.Internal(err, "failed to get upload URLs for user")
		}
	case "blood_req":
		exists, err := s.BloodRepo.ExistsByID(ctx, ID)
		if err != nil {
			return nil, apperrors.Internal(err, "failed to check blood request existence")
		}
		if !exists {
			return nil, apperrors.ErrBloodRequestNotFound
		}
		uploadInfos, err = s.storage.GetPresignedURLs(ctx, count, ID)
		if err != nil {
			return nil, apperrors.Internal(err, "failed to get upload URLs for blood request")
		}
	default:
		return nil, apperrors.BadRequest("unsupported preload type")
	}

	if len(uploadInfos) == 0 {
		return nil, apperrors.Internal(nil, "no upload info returned")
	}

	return uploadInfos, nil
}
