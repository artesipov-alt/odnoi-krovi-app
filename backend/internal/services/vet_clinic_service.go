package services

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/models"
	repositories "github.com/artesipov-alt/odnoi-krovi-app/internal/repositories/interfaces"
)

// VetClinicService определяет интерфейс для бизнес-логики ветеринарных клиник
type VetClinicService interface {
	// RegisterClinic регистрирует новую ветеринарную клинику в системе
	RegisterClinic(ctx context.Context, clinicData VetClinicRegistration) (*models.VetClinic, error)

	// GetClinicProfile получает полный профиль ветеринарной клиники
	GetClinicProfile(ctx context.Context, clinicID int) (*VetClinicProfile, error)

	// UpdateClinicProfile обновляет информацию о ветеринарной клинике
	UpdateClinicProfile(ctx context.Context, clinicID int, updates VetClinicUpdate) error

	// GetClinicsByLocationID получает все клиники по ID локации
	GetClinicsByLocationID(ctx context.Context, locationID int) ([]*models.VetClinic, error)

	// DeleteClinic удаляет клинику по ID (soft delete)
	DeleteClinic(ctx context.Context, clinicID int) error
}

// VetClinicRegistration содержит данные для регистрации ветеринарной клиники
type VetClinicRegistration struct {
	Name                     string  `json:"name" validate:"required,min=2,max=255"`
	Phone                    string  `json:"phone" validate:"omitempty,e164"`
	Website                  string  `json:"website,omitempty" validate:"omitempty,url"`
	WorkHours                string  `json:"workHours,omitempty"`
	Latitude                 float64 `json:"latitude,omitempty"`
	Longitude                float64 `json:"longitude,omitempty"`
	TransfusionConditions    string  `json:"transfusionConditions,omitempty"`
	DonorBonusPrograms       string  `json:"donorBonusPrograms,omitempty"`
	ContactPersonName        string  `json:"contactPersonName,omitempty"`
	ContactPersonPosition    string  `json:"contactPersonPosition,omitempty"`
	LocationID               int     `json:"locationId" validate:"required,min=1"`
	AppointmentRequirementID int     `json:"appointmentRequirementId" validate:"required,min=1"`
}

// VetClinicUpdate содержит поля, которые можно обновить для ветеринарной клиники
type VetClinicUpdate struct {
	Name                     *string `json:"name,omitempty" validate:"omitempty,min=2,max=255"`
	Phone                    *string `json:"phone,omitempty" validate:"omitempty,e164"`
	Website                  *string `json:"website,omitempty" validate:"omitempty,url"`
	WorkHours                *string `json:"workHours,omitempty"`
	TransfusionConditions    *string `json:"transfusionConditions,omitempty"`
	DonorBonusPrograms       *string `json:"donorBonusPrograms,omitempty"`
	ContactPersonName        *string `json:"contactPersonName,omitempty"`
	ContactPersonPosition    *string `json:"contactPersonPosition,omitempty"`
	LocationID               *int    `json:"locationId,omitempty" validate:"omitempty,min=1"`
	AppointmentRequirementID *int    `json:"appointmentRequirementId,omitempty" validate:"omitempty,min=1"`
}

// VetClinicProfile представляет полный профиль ветеринарной клиники
type VetClinicProfile struct {
	Clinic *models.VetClinic `json:"clinic"`
}

// VetClinicServiceImpl реализует VetClinicService
type VetClinicServiceImpl struct {
	vetClinicRepo repositories.VetClinicRepository
}

// NewVetClinicService создает новый сервис ветеринарных клиник
func NewVetClinicService(vetClinicRepo repositories.VetClinicRepository) *VetClinicServiceImpl {
	return &VetClinicServiceImpl{
		vetClinicRepo: vetClinicRepo,
	}
}

// RegisterClinic регистрирует новую ветеринарную клинику в системе
func (s *VetClinicServiceImpl) RegisterClinic(ctx context.Context, clinicData VetClinicRegistration) (*models.VetClinic, error) {
	// Создаем новую клинику
	clinic := &models.VetClinic{
		Name:                     clinicData.Name,
		Phone:                    clinicData.Phone,
		Website:                  clinicData.Website,
		WorkHours:                clinicData.WorkHours,
		Latitude:                 clinicData.Latitude,
		Longitude:                clinicData.Longitude,
		TransfusionConditions:    clinicData.TransfusionConditions,
		DonorBonusPrograms:       clinicData.DonorBonusPrograms,
		ContactPersonName:        clinicData.ContactPersonName,
		ContactPersonPosition:    clinicData.ContactPersonPosition,
		LocationID:               clinicData.LocationID,
		AppointmentRequirementID: clinicData.AppointmentRequirementID,
	}

	if err := s.vetClinicRepo.Create(ctx, clinic); err != nil {
		return nil, apperrors.Internal(err, "не удалось создать клинику")
	}

	return clinic, nil
}

// DeleteClinic удаляет клинику по ID (soft delete)
func (s *VetClinicServiceImpl) DeleteClinic(ctx context.Context, clinicID int) error {
	// Проверяем, существует ли клиника
	clinic, err := s.vetClinicRepo.GetByID(ctx, clinicID)
	if err != nil {
		return apperrors.Internal(err, "не удалось получить клинику")
	}

	if clinic == nil {
		return apperrors.NewClinicNotFoundError(clinicID)
	}

	// Удаляем клинику
	if err := s.vetClinicRepo.Delete(ctx, clinicID); err != nil {
		return apperrors.Internal(err, "не удалось удалить клинику")
	}

	return nil
}

// GetClinicProfile получает полный профиль ветеринарной клиники
func (s *VetClinicServiceImpl) GetClinicProfile(ctx context.Context, clinicID int) (*VetClinicProfile, error) {
	clinic, err := s.vetClinicRepo.GetByID(ctx, clinicID)
	if err != nil {
		return nil, apperrors.Internal(err, "не удалось получить клинику")
	}

	if clinic == nil {
		return nil, apperrors.NewClinicNotFoundError(clinicID)
	}

	profile := &VetClinicProfile{
		Clinic: clinic,
	}

	return profile, nil
}

// UpdateClinicProfile обновляет информацию о ветеринарной клинике
func (s *VetClinicServiceImpl) UpdateClinicProfile(ctx context.Context, clinicID int, updates VetClinicUpdate) error {
	// Получаем существующую клинику
	clinic, err := s.vetClinicRepo.GetByID(ctx, clinicID)
	if err != nil {
		return apperrors.Internal(err, "не удалось получить клинику")
	}

	if clinic == nil {
		return apperrors.NewClinicNotFoundError(clinicID)
	}

	// Применяем обновления
	if updates.Name != nil {
		clinic.Name = *updates.Name
	}
	if updates.Phone != nil {
		clinic.Phone = *updates.Phone
	}
	if updates.Website != nil {
		clinic.Website = *updates.Website
	}
	if updates.WorkHours != nil {
		clinic.WorkHours = *updates.WorkHours
	}
	if updates.TransfusionConditions != nil {
		clinic.TransfusionConditions = *updates.TransfusionConditions
	}
	if updates.DonorBonusPrograms != nil {
		clinic.DonorBonusPrograms = *updates.DonorBonusPrograms
	}
	if updates.ContactPersonName != nil {
		clinic.ContactPersonName = *updates.ContactPersonName
	}
	if updates.ContactPersonPosition != nil {
		clinic.ContactPersonPosition = *updates.ContactPersonPosition
	}
	if updates.LocationID != nil {
		clinic.LocationID = *updates.LocationID
	}
	if updates.AppointmentRequirementID != nil {
		clinic.AppointmentRequirementID = *updates.AppointmentRequirementID
	}

	// Сохраняем обновленную клинику
	if err := s.vetClinicRepo.Update(ctx, clinic); err != nil {
		return apperrors.Internal(err, "не удалось обновить клинику")
	}

	return nil
}

// GetClinicsByLocationID получает все клиники по ID локации
func (s *VetClinicServiceImpl) GetClinicsByLocationID(ctx context.Context, locationID int) ([]*models.VetClinic, error) {
	clinics, err := s.vetClinicRepo.GetByLocationID(ctx, locationID)
	if err != nil {
		return nil, apperrors.Internal(err, "не удалось получить клиники по локации")
	}
	return clinics, nil
}
