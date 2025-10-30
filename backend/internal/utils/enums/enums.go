package enums

import (
	"fmt"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/models"
)

// GetAllPetTypes возвращает все доступные типы животных
func GetAllPetTypes() []models.PetType {
	return []models.PetType{models.PetTypeDog, models.PetTypeCat}
}

// GetAllGenders возвращает все доступные значения пола
func GetAllGenders() []models.Gender {
	return []models.Gender{models.GenderMale, models.GenderFemale}
}

// GetAllLivingConditions возвращает все доступные условия проживания
func GetAllLivingConditions() []models.LivingCondition {
	return []models.LivingCondition{
		models.LivingConditionApartment,
		models.LivingConditionHouse,
		models.LivingConditionAviary,
		models.LivingConditionOther,
	}
}

// GetAllUserRoles возвращает все доступные роли пользователей
func GetAllUserRoles() []models.UserRole {
	return []models.UserRole{
		models.UserRoleUser,
		models.UserRoleClinic,
		models.UserRoleAdmin,
		models.UserRoleDonor,
	}
}

// GetAllBloodSearchStatuses возвращает все доступные статусы поиска крови
func GetAllBloodSearchStatuses() []models.BloodSearchStatus {
	return []models.BloodSearchStatus{
		models.BloodSearchStatusActive,
		models.BloodSearchStatusCompleted,
		models.BloodSearchStatusCancelled,
		models.BloodSearchStatusExpired,
	}
}

// GetAllBloodStockStatuses возвращает все доступные статусы запаса крови
func GetAllBloodStockStatuses() []models.BloodStockStatus {
	return []models.BloodStockStatus{
		models.BloodStockStatusActive,
		models.BloodStockStatusReserved,
		models.BloodStockStatusUsed,
		models.BloodStockStatusExpired,
	}
}

// GetAllDonationStatuses возвращает все доступные статусы донорства
func GetAllDonationStatuses() []models.DonationStatus {
	return []models.DonationStatus{
		models.DonationStatusScheduled,
		models.DonationStatusCompleted,
		models.DonationStatusCancelled,
		models.DonationStatusNoShow,
	}
}

// LocalizePetType локализует тип животного в русское название
func LocalizePetType(petType string) (models.PetType, error) {
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

// LocalizeGender локализует пол животного в русское название
func LocalizeGender(gender string) (models.Gender, error) {
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

// LocalizeLivingCondition локализует условие проживания в русское название
func LocalizeLivingCondition(condition string) (models.LivingCondition, error) {
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

// LocalizeUserRole локализует роль пользователя в русское название
func LocalizeUserRole(role string) (models.UserRole, error) {
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

// LocalizeBloodSearchStatus локализует статус поиска крови в русское название
func LocalizeBloodSearchStatus(status string) (models.BloodSearchStatus, error) {
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

// LocalizeBloodStockStatus локализует статус запаса крови в русское название
func LocalizeBloodStockStatus(status string) (models.BloodStockStatus, error) {
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

// LocalizeDonationStatus локализует статус донорства в русское название
func LocalizeDonationStatus(status string) (models.DonationStatus, error) {
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
