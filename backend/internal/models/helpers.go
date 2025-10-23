package models

import "fmt"

// ValidatePetType проверяет, является ли строка валидным типом животного
func ValidatePetType(petType string) (PetType, error) {
	pt := PetType(petType)
	switch pt {
	case PetTypeDog, PetTypeCat:
		return pt, nil
	default:
		return "", fmt.Errorf("недопустимый тип животного: %s (доступны: dog, cat)", petType)
	}
}

// ValidateGender проверяет, является ли строка валидным полом
func ValidateGender(gender string) (Gender, error) {
	g := Gender(gender)
	switch g {
	case GenderMale, GenderFemale:
		return g, nil
	default:
		return "", fmt.Errorf("недопустимый пол: %s (доступны: male, female)", gender)
	}
}

// ValidateLivingCondition проверяет, является ли строка валидным условием проживания
func ValidateLivingCondition(condition string) (LivingCondition, error) {
	lc := LivingCondition(condition)
	switch lc {
	case LivingConditionApartment, LivingConditionHouse, LivingConditionAviary, LivingConditionOther:
		return lc, nil
	default:
		return "", fmt.Errorf("недопустимое условие проживания: %s (доступны: apartment, house, aviary, other)", condition)
	}
}

// ValidateUserRole проверяет, является ли строка валидной ролью пользователя
func ValidateUserRole(role string) (UserRole, error) {
	r := UserRole(role)
	switch r {
	case UserRoleUser, UserRoleClinic, UserRoleAdmin, UserRoleDonor:
		return r, nil
	default:
		return "", fmt.Errorf("недопустимая роль: %s (доступны: user, clinic, admin, donor)", role)
	}
}

// ValidateBloodSearchStatus проверяет, является ли строка валидным статусом поиска крови
func ValidateBloodSearchStatus(status string) (BloodSearchStatus, error) {
	s := BloodSearchStatus(status)
	switch s {
	case BloodSearchStatusActive, BloodSearchStatusCompleted, BloodSearchStatusCancelled, BloodSearchStatusExpired:
		return s, nil
	default:
		return "", fmt.Errorf("недопустимый статус поиска: %s (доступны: active, completed, cancelled, expired)", status)
	}
}

// ValidateBloodStockStatus проверяет, является ли строка валидным статусом запаса крови
func ValidateBloodStockStatus(status string) (BloodStockStatus, error) {
	s := BloodStockStatus(status)
	switch s {
	case BloodStockStatusActive, BloodStockStatusReserved, BloodStockStatusUsed, BloodStockStatusExpired:
		return s, nil
	default:
		return "", fmt.Errorf("недопустимый статус запаса: %s (доступны: active, reserved, used, expired)", status)
	}
}

// ValidateDonationStatus проверяет, является ли строка валидным статусом донорства
func ValidateDonationStatus(status string) (DonationStatus, error) {
	s := DonationStatus(status)
	switch s {
	case DonationStatusScheduled, DonationStatusCompleted, DonationStatusCancelled, DonationStatusNoShow:
		return s, nil
	default:
		return "", fmt.Errorf("недопустимый статус донорства: %s (доступны: scheduled, completed, cancelled, no_show)", status)
	}
}

// GetAllPetTypes возвращает все доступные типы животных
func GetAllPetTypes() []PetType {
	return []PetType{PetTypeDog, PetTypeCat}
}

// GetAllGenders возвращает все доступные значения пола
func GetAllGenders() []Gender {
	return []Gender{GenderMale, GenderFemale}
}

// GetAllLivingConditions возвращает все доступные условия проживания
func GetAllLivingConditions() []LivingCondition {
	return []LivingCondition{
		LivingConditionApartment,
		LivingConditionHouse,
		LivingConditionAviary,
		LivingConditionOther,
	}
}

// GetAllUserRoles возвращает все доступные роли пользователей
func GetAllUserRoles() []UserRole {
	return []UserRole{UserRoleUser, UserRoleClinic, UserRoleAdmin, UserRoleDonor}
}

// GetAllBloodSearchStatuses возвращает все доступные статусы поиска крови
func GetAllBloodSearchStatuses() []BloodSearchStatus {
	return []BloodSearchStatus{
		BloodSearchStatusActive,
		BloodSearchStatusCompleted,
		BloodSearchStatusCancelled,
		BloodSearchStatusExpired,
	}
}

// GetAllBloodStockStatuses возвращает все доступные статусы запаса крови
func GetAllBloodStockStatuses() []BloodStockStatus {
	return []BloodStockStatus{
		BloodStockStatusActive,
		BloodStockStatusReserved,
		BloodStockStatusUsed,
		BloodStockStatusExpired,
	}
}

// GetAllDonationStatuses возвращает все доступные статусы донорства
func GetAllDonationStatuses() []DonationStatus {
	return []DonationStatus{
		DonationStatusScheduled,
		DonationStatusCompleted,
		DonationStatusCancelled,
		DonationStatusNoShow,
	}
}

// String методы для красивого вывода

// String возвращает строковое представление типа животного на русском
func (pt PetType) String() string {
	switch pt {
	case PetTypeDog:
		return "Собака"
	case PetTypeCat:
		return "Кошка"
	default:
		return string(pt)
	}
}

// String возвращает строковое представление пола на русском
func (g Gender) String() string {
	switch g {
	case GenderMale:
		return "Самец"
	case GenderFemale:
		return "Самка"
	default:
		return string(g)
	}
}

// String возвращает строковое представление условия проживания на русском
func (lc LivingCondition) String() string {
	switch lc {
	case LivingConditionApartment:
		return "Квартира"
	case LivingConditionHouse:
		return "Дом"
	case LivingConditionAviary:
		return "Вольер"
	case LivingConditionOther:
		return "Другое"
	default:
		return string(lc)
	}
}

// String возвращает строковое представление роли пользователя на русском
func (r UserRole) String() string {
	switch r {
	case UserRoleUser:
		return "Пользователь"
	case UserRoleClinic:
		return "Клиника"
	case UserRoleAdmin:
		return "Администратор"
	case UserRoleDonor:
		return "Донор"
	default:
		return string(r)
	}
}

// String возвращает строковое представление статуса поиска на русском
func (s BloodSearchStatus) String() string {
	switch s {
	case BloodSearchStatusActive:
		return "Активный"
	case BloodSearchStatusCompleted:
		return "Завершен"
	case BloodSearchStatusCancelled:
		return "Отменен"
	case BloodSearchStatusExpired:
		return "Истек"
	default:
		return string(s)
	}
}

// String возвращает строковое представление статуса запаса на русском
func (s BloodStockStatus) String() string {
	switch s {
	case BloodStockStatusActive:
		return "Активный"
	case BloodStockStatusReserved:
		return "Зарезервирован"
	case BloodStockStatusUsed:
		return "Использован"
	case BloodStockStatusExpired:
		return "Истек срок"
	default:
		return string(s)
	}
}

// String возвращает строковое представление статуса донорства на русском
func (s DonationStatus) String() string {
	switch s {
	case DonationStatusScheduled:
		return "Запланировано"
	case DonationStatusCompleted:
		return "Завершено"
	case DonationStatusCancelled:
		return "Отменено"
	case DonationStatusNoShow:
		return "Не явился"
	default:
		return string(s)
	}
}
