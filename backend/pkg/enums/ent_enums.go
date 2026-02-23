package enums

import (
	"github.com/artesipov-alt/odnoi-krovi-app/ent/user"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain"
)

// GetAllEntPetTypes возвращает все доступные типы животных из ENT
func GetAllEntPetTypes() []domain.PetType {
	return []domain.PetType{domain.PetTypeDog, domain.PetTypeCat}
}

// GetAllEntPetStatuses возвращает все доступные статусы питомцев из ENT
func GetAllEntPetStatuses() []domain.PetStatus {
	return []domain.PetStatus{domain.PetStatusDonor, domain.PetStatusRecipient, domain.PetStatusNone}
}

// GetAllEntGenders возвращает все доступные значения пола из ENT
func GetAllEntGenders() []domain.Gender {
	return []domain.Gender{domain.GenderMale, domain.GenderFemale}
}

// GetAllEntLivingConditions возвращает все доступные условия проживания из ENT
func GetAllEntLivingConditions() []domain.LivingCondition {
	return []domain.LivingCondition{
		domain.LivingConditionIndoor,
		domain.LivingConditionLeashWalking,
		domain.LivingConditionSelfOutdoor,
	}
}

// GetAllEntHealthStatuses возвращает все доступные статусы здоровья из ENT
func GetAllEntHealthStatuses() []domain.HealthStatus {
	return []domain.HealthStatus{
		domain.HealthStatusHealthy,
		domain.HealthStatusIll,
		domain.HealthStatusUnknown,
	}
}

// GetAllEntReproductiveStatuses возвращает все доступные репродуктивные состояния из ENT
func GetAllEntReproductiveStatuses() []domain.ReproductiveStatus {
	return []domain.ReproductiveStatus{
		domain.ReproductiveStatusPregnancy,
		domain.ReproductiveStatusLactation,
		domain.ReproductiveStatusEstrus,
	}
}

// GetAllEntUserRoles возвращает все доступные роли пользователей из ENT
func GetAllEntUserRoles() []user.Role {
	return []user.Role{
		user.RoleUser,
		user.RoleAdmin,
	}
}

// LocalizeEntPetType локализует тип животного из ENT
func LocalizeEntPetType(pt domain.PetType) string {
	switch pt {
	case domain.PetTypeDog:
		return "Собака"
	case domain.PetTypeCat:
		return "Кошка"
	default:
		return string(pt)
	}
}

// LocalizeEntPetStatus локализует статус питомца из ENT
func LocalizeEntPetStatus(ps domain.PetStatus) string {
	switch ps {
	case domain.PetStatusDonor:
		return "Донор"
	case domain.PetStatusRecipient:
		return "Реципиент"
	default:
		return string(ps)
	}
}

// LocalizeEntGender локализует пол животного из ENT
func LocalizeEntGender(g domain.Gender) string {
	switch g {
	case domain.GenderMale:
		return "Самец"
	case domain.GenderFemale:
		return "Самка"
	default:
		return string(g)
	}
}

// LocalizeEntLivingCondition локализует условие проживания из ENT
func LocalizeEntLivingCondition(lc domain.LivingCondition) string {
	switch lc {
	case domain.LivingConditionIndoor:
		return "Домашний"
	case domain.LivingConditionLeashWalking:
		return "Выгул на шлейке"
	case domain.LivingConditionSelfOutdoor:
		return "Самовыгул"
	default:
		return string(lc)
	}
}

// LocalizeEntHealthStatus локализует состояние здоровья из ENT
func LocalizeEntHealthStatus(hs domain.HealthStatus) string {
	switch hs {
	case domain.HealthStatusHealthy:
		return "Здоров"
	case domain.HealthStatusIll:
		return "Есть заболевания"
	case domain.HealthStatusUnknown:
		return "Неизвестно"
	default:
		return string(hs)
	}
}

// LocalizeEntUserRole локализует роль пользователя из ENT
func LocalizeEntUserRole(r user.Role) string {
	switch r {
	case user.RoleUser:
		return "Пользователь"
	case user.RoleAdmin:
		return "Администратор"
	default:
		return string(r)
	}
}

// LocalizeEntReproductiveStatus локализует физиологическое состояние из ENT
func LocalizeEntReproductiveStatus(rs domain.ReproductiveStatus) string {
	switch rs {
	case domain.ReproductiveStatusPregnancy:
		return "Беременность"
	case domain.ReproductiveStatusLactation:
		return "Лактация"
	case domain.ReproductiveStatusEstrus:
		return "Течка"
	default:
		return string(rs)
	}
}
