package enums

import (
	"github.com/artesipov-alt/odnoi-krovi-app/ent/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/pethealth"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/user"
)

// GetAllEntPetTypes возвращает все доступные типы животных из ENT
func GetAllEntPetTypes() []pet.Type {
	return []pet.Type{pet.TypeDog, pet.TypeCat}
}

// GetAllEntPetStatuses возвращает все доступные статусы питомцев из ENT
// func GetAllEntPetStatuses() []pet.PetStatus {
// 	return []pet.PetStatus{pet.PetStatusDonor, pet.PetStatusRecipient, pet.PetStatusNone}
// }

// GetAllEntGenders возвращает все доступные значения пола из ENT
func GetAllEntGenders() []pet.Gender {
	return []pet.Gender{pet.GenderMale, pet.GenderFemale}
}

// GetAllEntLivingConditions возвращает все доступные условия проживания из ENT
func GetAllEntLivingConditions() []pet.LivingCondition {
	return []pet.LivingCondition{
		pet.LivingConditionIndoor,
		pet.LivingConditionLeashWalking,
		pet.LivingConditionSelfOutdoor,
	}
}

// GetAllEntHealthStatuses возвращает все доступные статусы здоровья из ENT
func GetAllEntHealthStatuses() []pethealth.HealthStatus {
	return []pethealth.HealthStatus{
		pethealth.HealthStatusHealthy,
		pethealth.HealthStatusIll,
		pethealth.HealthStatusUnknown,
	}
}

// GetAllEntReproductiveStatuses возвращает все доступные репродуктивные состояния из ENT
func GetAllEntReproductiveStatuses() []pet.ReproductiveStatus {
	return []pet.ReproductiveStatus{
		pet.ReproductiveStatusPregnancy,
		pet.ReproductiveStatusLactation,
		pet.ReproductiveStatusEstrus,
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
func LocalizeEntPetType(pt pet.Type) string {
	switch pt {
	case pet.TypeDog:
		return "Собака"
	case pet.TypeCat:
		return "Кошка"
	default:
		return string(pt)
	}
}

// LocalizeEntPetStatus локализует статус питомца из ENT
// func LocalizeEntPetStatus(ps pet.PetStatus) string {
// 	switch ps {
// 	case pet.PetStatusDonor:
// 		return "Донор"
// 	case pet.PetStatusRecipient:
// 		return "Реципиент"
// 	default:
// 		return string(ps)
// 	}
// }

// LocalizeEntGender локализует пол животного из ENT
func LocalizeEntGender(g pet.Gender) string {
	switch g {
	case pet.GenderMale:
		return "Самец"
	case pet.GenderFemale:
		return "Самка"
	default:
		return string(g)
	}
}

// LocalizeEntLivingCondition локализует условие проживания из ENT
func LocalizeEntLivingCondition(lc pet.LivingCondition) string {
	switch lc {
	case pet.LivingConditionIndoor:
		return "Домашний"
	case pet.LivingConditionLeashWalking:
		return "Выгул на шлейке"
	case pet.LivingConditionSelfOutdoor:
		return "Самовыгул"
	default:
		return string(lc)
	}
}

// LocalizeEntHealthStatus локализует состояние здоровья из ENT
func LocalizeEntHealthStatus(hs pethealth.HealthStatus) string {
	switch hs {
	case pethealth.HealthStatusHealthy:
		return "Здоров"
	case pethealth.HealthStatusIll:
		return "Есть заболевания"
	case pethealth.HealthStatusUnknown:
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
func LocalizeEntReproductiveStatus(rs pet.ReproductiveStatus) string {
	switch rs {
	case pet.ReproductiveStatusPregnancy:
		return "Беременность"
	case pet.ReproductiveStatusLactation:
		return "Лактация"
	case pet.ReproductiveStatusEstrus:
		return "Течка"
	default:
		return string(rs)
	}
}
