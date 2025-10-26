package enums

import "github.com/artesipov-alt/odnoi-krovi-app/internal/models"

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
