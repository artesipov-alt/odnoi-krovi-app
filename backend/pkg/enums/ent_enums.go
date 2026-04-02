package enums

import (
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/common"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/user"
)

// GetAllEntPetTypes возвращает все доступные типы животных из ENT
func GetAllEntPetTypes() []common.PetType {
	return []common.PetType{common.PetTypeDog, common.PetTypeCat}
}

// GetAllEntPetStatuses возвращает все доступные статусы питомцев из ENT
func GetAllEntPetStatuses() []model.PetStatus {
	return []model.PetStatus{model.PetStatusDonor, model.PetStatusRecipient, model.PetStatusNone}
}

// GetAllEntGenders возвращает все доступные значения пола из ENT
func GetAllEntGenders() []model.Gender {
	return []model.Gender{model.GenderMale, model.GenderFemale}
}

// GetAllEntLivingConditions возвращает все доступные условия проживания из ENT
func GetAllEntLivingConditions() []model.LivingCondition {
	return []model.LivingCondition{
		model.LivingConditionIndoor,
		model.LivingConditionLeashWalking,
		model.LivingConditionSelfOutdoor,
	}
}

// GetAllEntHealthStatuses возвращает все доступные статусы здоровья из ENT
func GetAllEntHealthStatuses() []model.HealthStatus {
	return []model.HealthStatus{
		model.HealthStatusHealthy,
		model.HealthStatusIll,
		model.HealthStatusUnknown,
	}
}

// GetAllEntReproductiveStatuses возвращает все доступные репродуктивные состояния из ENT
func GetAllEntReproductiveStatuses() []model.ReproductiveStatus {
	return []model.ReproductiveStatus{
		model.ReproductiveStatusPregnancy,
		model.ReproductiveStatusLactation,
		model.ReproductiveStatusEstrus,
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
func LocalizeEntPetType(pt common.PetType) string {
	switch pt {
	case common.PetTypeDog:
		return "Собака"
	case common.PetTypeCat:
		return "Кошка"
	default:
		return string(pt)
	}
}

// LocalizeEntPetStatus локализует статус питомца из ENT
func LocalizeEntPetStatus(ps model.PetStatus) string {
	switch ps {
	case model.PetStatusDonor:
		return "Донор"
	case model.PetStatusRecipient:
		return "Реципиент"
	default:
		return string(ps)
	}
}

// LocalizeEntGender локализует пол животного из ENT
func LocalizeEntGender(g model.Gender) string {
	switch g {
	case model.GenderMale:
		return "Самец"
	case model.GenderFemale:
		return "Самка"
	default:
		return string(g)
	}
}

// LocalizeEntLivingCondition локализует условие проживания из ENT
func LocalizeEntLivingCondition(lc model.LivingCondition) string {
	switch lc {
	case model.LivingConditionIndoor:
		return "Домашний"
	case model.LivingConditionLeashWalking:
		return "Выгул на шлейке"
	case model.LivingConditionSelfOutdoor:
		return "Самовыгул"
	default:
		return string(lc)
	}
}

// LocalizeEntHealthStatus локализует состояние здоровья из ENT
func LocalizeEntHealthStatus(hs model.HealthStatus) string {
	switch hs {
	case model.HealthStatusHealthy:
		return "Здоров"
	case model.HealthStatusIll:
		return "Есть заболевания"
	case model.HealthStatusUnknown:
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
func LocalizeEntReproductiveStatus(rs model.ReproductiveStatus) string {
	switch rs {
	case model.ReproductiveStatusPregnancy:
		return "Беременность"
	case model.ReproductiveStatusLactation:
		return "Лактация"
	case model.ReproductiveStatusEstrus:
		return "Течка"
	default:
		return string(rs)
	}
}
