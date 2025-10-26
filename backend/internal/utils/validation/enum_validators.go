package validation

import (
	"fmt"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/models"
)

// ValidatePetType проверяет, является ли строка валидным типом животного
func ValidatePetType(petType string) (models.PetType, error) {
	pt := models.PetType(petType)
	switch pt {
	case models.PetTypeDog, models.PetTypeCat:
		return pt, nil
	default:
		return "", fmt.Errorf("недопустимый тип животного: %s (доступны: dog, cat)", petType)
	}
}

// ValidateGender проверяет, является ли строка валидным полом
func ValidateGender(gender string) (models.Gender, error) {
	g := models.Gender(gender)
	switch g {
	case models.GenderMale, models.GenderFemale:
		return g, nil
	default:
		return "", fmt.Errorf("недопустимый пол: %s (доступны: male, female)", gender)
	}
}

// ValidateLivingCondition проверяет, является ли строка валидным условием проживания
func ValidateLivingCondition(condition string) (models.LivingCondition, error) {
	lc := models.LivingCondition(condition)
	switch lc {
	case models.LivingConditionApartment, models.LivingConditionHouse, models.LivingConditionAviary, models.LivingConditionOther:
		return lc, nil
	default:
		return "", fmt.Errorf("недопустимое условие проживания: %s (доступны: apartment, house, aviary, other)", condition)
	}
}

// ValidateUserRole проверяет, является ли строка валидной ролью пользователя
func ValidateUserRole(role string) (models.UserRole, error) {
	r := models.UserRole(role)
	switch r {
	case models.UserRoleUser, models.UserRoleClinic, models.UserRoleAdmin, models.UserRoleDonor:
		return r, nil
	default:
		return "", fmt.Errorf("недопустимая роль: %s (доступны: user, clinic, admin, donor)", role)
	}
}

// ValidateBloodSearchStatus проверяет, является ли строка валидным статусом поиска крови
func ValidateBloodSearchStatus(status string) (models.BloodSearchStatus, error) {
	s := models.BloodSearchStatus(status)
	switch s {
	case models.BloodSearchStatusActive, models.BloodSearchStatusCompleted, models.BloodSearchStatusCancelled, models.BloodSearchStatusExpired:
		return s, nil
	default:
		return "", fmt.Errorf("недопустимый статус поиска: %s (доступны: active, completed, cancelled, expired)", status)
	}
}

// ValidateBloodStockStatus проверяет, является ли строка валидным статусом запаса крови
func ValidateBloodStockStatus(status string) (models.BloodStockStatus, error) {
	s := models.BloodStockStatus(status)
	switch s {
	case models.BloodStockStatusActive, models.BloodStockStatusReserved, models.BloodStockStatusUsed, models.BloodStockStatusExpired:
		return s, nil
	default:
		return "", fmt.Errorf("недопустимый статус запаса: %s (доступны: active, reserved, used, expired)", status)
	}
}

// ValidateDonationStatus проверяет, является ли строка валидным статусом донорства
func ValidateDonationStatus(status string) (models.DonationStatus, error) {
	s := models.DonationStatus(status)
	switch s {
	case models.DonationStatusScheduled, models.DonationStatusCompleted, models.DonationStatusCancelled, models.DonationStatusNoShow:
		return s, nil
	default:
		return "", fmt.Errorf("недопустимый статус донорства: %s (доступны: scheduled, completed, cancelled, no_show)", status)
	}
}
