package validation

import (
	"fmt"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/models"
)

// ValidatePetType проверяет, является ли строка валидным типом животного
func ValidatePetType(petType string) (models.PetType, error) {
	pt := models.PetType(petType)
	switch pt {
	case models.PetTypeDog:
		return "Собака", nil
	case models.PetTypeCat:
		return "Кошка", nil
	default:
		return "", fmt.Errorf("недопустимый тип животного: %s", petType)
	}
}

// ValidateGender проверяет, является ли строка валидным полом
func ValidateGender(gender string) (models.Gender, error) {
	g := models.Gender(gender)
	switch g {
	case models.GenderMale:
		return "Самец", nil
	case models.GenderFemale:
		return "Самка", nil
	default:
		return "", fmt.Errorf("недопустимый пол: %s", gender)
	}
}

// ValidateLivingCondition проверяет, является ли строка валидным условием проживания
func ValidateLivingCondition(condition string) (models.LivingCondition, error) {
	lc := models.LivingCondition(condition)
	switch lc {
	case models.LivingConditionApartment:
		return "Квартира", nil
	case models.LivingConditionHouse:
		return "Дом", nil
	case models.LivingConditionAviary:
		return "Вольер", nil
	case models.LivingConditionOther:
		return "Другое", nil
	default:
		return "", fmt.Errorf("недопустимое условие проживания: %s", condition)
	}
}

// ValidateUserRole проверяет, является ли строка валидной ролью пользователя
func ValidateUserRole(role string) (models.UserRole, error) {
	r := models.UserRole(role)
	switch r {
	case models.UserRoleUser:
		return "Пользователь", nil
	case models.UserRoleClinic:
		return "Клиника", nil
	case models.UserRoleAdmin:
		return "Администратор", nil
	case models.UserRoleDonor:
		return "Донор", nil
	default:
		return "", fmt.Errorf("недопустимая роль: %s", role)
	}
}

// ValidateBloodSearchStatus проверяет, является ли строка валидным статусом поиска крови
func ValidateBloodSearchStatus(status string) (models.BloodSearchStatus, error) {
	s := models.BloodSearchStatus(status)
	switch s {
	case models.BloodSearchStatusActive:
		return "Активный", nil
	case models.BloodSearchStatusCompleted:
		return "Завершен", nil
	case models.BloodSearchStatusCancelled:
		return "Отменен", nil
	case models.BloodSearchStatusExpired:
		return "Истек", nil
	default:
		return "", fmt.Errorf("недопустимый статус поиска: %s", status)
	}
}

// ValidateBloodStockStatus проверяет, является ли строка валидным статусом запаса крови
func ValidateBloodStockStatus(status string) (models.BloodStockStatus, error) {
	s := models.BloodStockStatus(status)
	switch s {
	case models.BloodStockStatusActive:
		return "Активный", nil
	case models.BloodStockStatusReserved:
		return "Зарезервирован", nil
	case models.BloodStockStatusUsed:
		return "Использован", nil
	case models.BloodStockStatusExpired:
		return "Истек", nil
	default:
		return "", fmt.Errorf("недопустимый статус запаса: %s", status)
	}
}

// ValidateDonationStatus проверяет, является ли строка валидным статусом донорства
func ValidateDonationStatus(status string) (models.DonationStatus, error) {
	s := models.DonationStatus(status)
	switch s {
	case models.DonationStatusScheduled:
		return "Запланирован", nil
	case models.DonationStatusCompleted:
		return "Завершен", nil
	case models.DonationStatusCancelled:
		return "Отменен", nil
	case models.DonationStatusNoShow:
		return "Не явился", nil
	default:
		return "", fmt.Errorf("недопустимый статус донорства: %s", status)
	}
}
