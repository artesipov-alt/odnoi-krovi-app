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
		models.LivingConditionIndoor,
		models.LivingConditionLeash,
		models.LivingConditionOutdoor,
	}
}

// GetAllHealthStatuses возвращает все доступные статусы здоровья
func GetAllHealthStatuses() []models.HealthStatus {
	return []models.HealthStatus{
		models.HealthStatusHealthy,
		models.HealthStatusIll,
		models.HealthStatusUnknown,
	}
}

// GetAllReproductiveStatuses возвращает все доступные физиологические состояния
func GetAllReproductiveStatuses() []models.ReproductiveStatus {
	return []models.ReproductiveStatus{
		models.ReproductiveStatusPregnancy,
		models.ReproductiveStatusLactation,
		models.ReproductiveStatusEstrus,
		models.ReproductiveStatusNone,
	}
}

// GetAllUserRoles возвращает все доступные роли пользователей
func GetAllUserRoles() []models.UserRole {
	return []models.UserRole{
		models.UserRoleUser,
		models.UserRoleAdmin,
	}
}

// GetAllPetRoles возвращает все доступные роли животных
func GetAllPetRoles() []models.PetRole {
	return []models.PetRole{
		models.PetRoleRecipient,
		models.PetRoleDonor,
	}
}

// LocalizePetType локализует тип животного в русское название
// Можно валидировать через проверку.
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
func LocalizeGender(gender string) (string, error) {
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
	case models.LivingConditionIndoor:
		return "Домашний", nil
	case models.LivingConditionLeash:
		return "Выгул на шлейке", nil
	case models.LivingConditionOutdoor:
		return "Самовыгул", nil
	default:
		return "", fmt.Errorf("недопустимое условие проживания: %s", condition)
	}
}

// LocalizeHealthStatus локализует состояние здоровья в русское название
func LocalizeHealthStatus(status string) (models.HealthStatus, error) {
	hs := models.HealthStatus(status)
	switch hs {
	case models.HealthStatusHealthy:
		return "Здоров", nil
	case models.HealthStatusIll:
		return "Есть заболевания", nil
	case models.HealthStatusUnknown:
		return "Неизвестно", nil
	default:
		return "", fmt.Errorf("недопустимый статус здоровья: %s", status)
	}
}

// LocalizeUserRole локализует роль пользователя в русское название
func LocalizeUserRole(role string) (models.UserRole, error) {
	r := models.UserRole(role)
	switch r {
	case models.UserRoleUser:
		return "Пользователь", nil
	case models.UserRoleAdmin:
		return "Администратор", nil
	default:
		return "", fmt.Errorf("недопустимая роль: %s", role)
	}
}

// LocalizePetRole локализует роль животного в русское название
func LocalizePetRole(role string) (models.PetRole, error) {
	r := models.PetRole(role)
	switch r {
	case models.PetRoleDonor:
		return "Донор", nil
	case models.PetRoleRecipient:
		return "Реципиент", nil
	default:
		return "", fmt.Errorf("недопустимая роль: %s", role)
	}
}

// LocalizeReproductiveStatus локализует физиологическое состояние в русское название
func LocalizeReproductiveStatus(status string) (models.ReproductiveStatus, error) {
	rs := models.ReproductiveStatus(status)
	switch rs {
	case models.ReproductiveStatusPregnancy:
		return "Беременность", nil
	case models.ReproductiveStatusLactation:
		return "Лактация", nil
	case models.ReproductiveStatusEstrus:
		return "Течка", nil
	case models.ReproductiveStatusNone:
		return "Нет", nil
	default:
		return "", fmt.Errorf("недопустимое физиологическое состояние: %s", status)
	}
}
