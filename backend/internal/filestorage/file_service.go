package filestorage

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain"
)

// FileService реализует FileStorage
type FileService struct {
	PetRepo   domain.PetRepository
	UserRepo  domain.UserRepository
	BloodRepo domain.BloodRequestRepository
	storage   domain.FileStorage
}

// NewFileService создает новый FileService
func NewFileService(petRepo domain.PetRepository, userRepo domain.UserRepository, bloodRepo domain.BloodRequestRepository, storage domain.FileStorage) *FileService {
	return &FileService{
		PetRepo:   petRepo,
		UserRepo:  userRepo,
		BloodRepo: bloodRepo,
		storage:   storage,
	}
}

// GetPresignURLs Возвращает ссылки для загрузки фотографий.
func (s *FileService) GetPresignURLs(ctx context.Context, ID string, count int64, preloads ...string) ([]UploadInfo, error) {
	if len(preloads) == 0 {
		return nil, apperrors.BadRequest("preload type is required")
	}

	var uploadInfos []UploadInfo

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

// ConfirmUploads подтверждает загрузку фото для сущности
func (s *FileService) ConfirmUploads(ctx context.Context, ID string, paths []string, preload string) error {

	var (
		exists bool
		err    error
	)

	switch preload {
	case "pet_avatar":
		exists, err = s.PetRepo.ExistsByID(ctx, ID)
		if err != nil {
			return apperrors.Internal(err, "failed to check pet existence")
		}
		if !exists {
			return apperrors.ErrPetNotFound
		}
		err = s.storage.ConfirmUploads(ctx, paths)
		if err != nil {
			return apperrors.Internal(err, "failed to confirm uploads for pet")
		}
		// Обновить PhotoUrls в питомце через репозиторий (замена на новые)
		err = s.PetRepo.AddPhotoURLs(ctx, ID, paths)
		if err != nil {
			return apperrors.Internal(err, "failed to update pet photos")
		}
	case "user_avatar":
		exists, err = s.UserRepo.ExistsByID(ctx, ID)
		if err != nil {
			return apperrors.Internal(err, "failed to check user existence")
		}
		if !exists {
			return apperrors.ErrUserNotFound
		}
		err = s.storage.ConfirmUploads(ctx, paths)
		if err != nil {
			return apperrors.Internal(err, "failed to confirm uploads for user")
		}
		// Обновить PhotoUrls в пользователе через репозиторий
		if err := s.UserRepo.AddPhotoURLs(ctx, ID, paths); err != nil {
			return apperrors.Internal(err, "failed to update user photos")
		}
	case "blood_req":
		exists, err = s.BloodRepo.ExistsByID(ctx, ID)
		if err != nil {
			return apperrors.Internal(err, "failed to check blood request existence")
		}
		if !exists {
			return apperrors.ErrBloodRequestNotFound
		}
		err = s.storage.ConfirmUploads(ctx, paths)
		if err != nil {
			return apperrors.Internal(err, "failed to confirm uploads for blood request")
		}
		// Обновить PhotoUrls в заявке через репозиторий
		if err := s.BloodRepo.AddPhotoURLs(ctx, ID, paths); err != nil {
			return apperrors.Internal(err, "failed to update blood request photos")
		}
	default:
		return apperrors.BadRequest("unsupported preload type")
	}

	return nil
}
