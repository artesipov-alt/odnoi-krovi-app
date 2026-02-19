// internal/services/file_service.go
package services

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/repositories"
)

type FileService interface {
	GetPresignURLs(ctx context.Context, ID string, count int64, preloads ...string) ([]repositories.UploadInfo, error)
	ConfirmUploads(ctx context.Context, ID string, paths []string, preload string) error
}

// FileServiceImpl реализует FileService
type FileServiceImpl struct {
	PetRepo   repositories.PetRepository
	UserRepo  repositories.UserRepository
	BloodRepo BloodRequestRepository
	storage   repositories.FileStorage
}

// NewFileService создает новый FileServiceImpl
func NewFileService(petRepo repositories.PetRepository, userRepo repositories.UserRepository, bloodRepo BloodRequestRepository, storage repositories.FileStorage) *FileServiceImpl {
	return &FileServiceImpl{
		PetRepo:   petRepo,
		UserRepo:  userRepo,
		BloodRepo: bloodRepo,
		storage:   storage,
	}
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

// ConfirmUploads подтверждает загрузку фото для сущности
func (s *FileServiceImpl) ConfirmUploads(ctx context.Context, ID string, paths []string, preload string) error {
	switch preload {
	case "pet_avatar":
		exists, err := s.PetRepo.ExistsByID(ctx, ID)
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
		// Сначала проверяем, существует ли питомец
		_, err = s.PetRepo.GetPetQuery(ctx, ID).Only(ctx)
		if err != nil {
			if ent.IsNotFound(err) {
				return apperrors.ErrPetNotFound
			}
			return apperrors.Internal(err, "failed to get pet for update")
		}
		_, err = s.PetRepo.Update(ctx, ID, &ent.UpdatePetInput{
			PhotoUrls: paths,
		}, nil, nil, nil, nil)
		if err != nil {
			return apperrors.Internal(err, "failed to update pet photos")
		}
	case "user_avatar":
		exists, err := s.UserRepo.ExistsByID(ctx, ID)
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
		exists, err := s.BloodRepo.ExistsByID(ctx, ID)
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
