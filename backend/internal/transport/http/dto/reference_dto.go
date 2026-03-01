package dto

// ============================================
// Path Parameters
// ============================================

// PetTypePath представляет параметр пути с типом питомца
type PetTypePath struct {
	PetType string `path:"pet_type" doc:"Тип животного" enum:"dog,cat" example:"dog"`
}

// ============================================
// Query Parameters
// ============================================

// PetTypeQuery представляет параметр запроса с типом питомца
type PetTypeQuery struct {
	PetType string `query:"petType" doc:"Тип животного" enum:"dog,cat" example:"dog"`
}

// ============================================
// Get Pet Types
// ============================================

// GetPetTypesInput представляет запрос на получение типов питомцев
type GetPetTypesInput struct{}

// GetPetTypesOutput представляет ответ со списком типов питомцев
type GetPetTypesOutput struct {
	Body PetTypesList
}

// PetTypesList представляет список типов питомцев
type PetTypesList struct {
	Data []ReferenceItem `json:"data" doc:"Список типов питомцев"`
}

// ============================================
// Get Genders
// ============================================

// GetGendersInput представляет запрос на получение полов
type GetGendersInput struct{}

// GetGendersOutput представляет ответ со списком полов
type GetGendersOutput struct {
	Body GendersList
}

// GendersList представляет список полов
type GendersList struct {
	Data []ReferenceItem `json:"data" doc:"Список полов"`
}

// ============================================
// Get Living Conditions
// ============================================

// GetLivingConditionsInput представляет запрос на получение условий проживания
type GetLivingConditionsInput struct{}

// GetLivingConditionsOutput представляет ответ со списком условий проживания
type GetLivingConditionsOutput struct {
	Body LivingConditionsList
}

// LivingConditionsList представляет список условий проживания
type LivingConditionsList struct {
	Data []ReferenceItem `json:"data" doc:"Список условий проживания"`
}

// ============================================
// Get User Roles
// ============================================

// GetUserRolesInput представляет запрос на получение ролей пользователей
type GetUserRolesInput struct{}

// GetUserRolesOutput представляет ответ со списком ролей пользователей
type GetUserRolesOutput struct {
	Body UserRolesList
}

// UserRolesList представляет список ролей пользователей
type UserRolesList struct {
	Data []ReferenceItem `json:"data" doc:"Список ролей пользователей"`
}

// ============================================
// Get Pet Roles
// ============================================

// GetPetRolesInput представляет запрос на получение ролей питомцев
type GetPetRolesInput struct{}

// GetPetRolesOutput представляет ответ со списком ролей питомцев
type GetPetRolesOutput struct {
	Body PetRolesList
}

// PetRolesList представляет список ролей питомцев
type PetRolesList struct {
	Data []ReferenceItem `json:"data" doc:"Список ролей питомцев"`
}

// ============================================
// Get Health Statuses
// ============================================

// GetHealthStatusesInput представляет запрос на получение статусов здоровья
type GetHealthStatusesInput struct{}

// GetHealthStatusesOutput представляет ответ со списком статусов здоровья
type GetHealthStatusesOutput struct {
	Body HealthStatusesList
}

// HealthStatusesList представляет список статусов здоровья
type HealthStatusesList struct {
	Data []ReferenceItem `json:"data" doc:"Список статусов здоровья"`
}

// ============================================
// Get Reproductive Statuses
// ============================================

// GetReproductiveStatusesInput представляет запрос на получение репродуктивных статусов
type GetReproductiveStatusesInput struct{}

// GetReproductiveStatusesOutput представляет ответ со списком репродуктивных статусов
type GetReproductiveStatusesOutput struct {
	Body ReproductiveStatusesList
}

// ReproductiveStatusesList представляет список репродуктивных статусов
type ReproductiveStatusesList struct {
	Data []ReferenceItem `json:"data" doc:"Список репродуктивных статусов"`
}

// ============================================
// Get Breeds
// ============================================

// GetBreedsInput представляет запрос на получение списка пород
type GetBreedsInput struct{}

// GetBreedsOutput представляет ответ со списком пород
type GetBreedsOutput struct {
	Body BreedsList
}

// BreedsList представляет список пород
type BreedsList struct {
	Data []ReferenceItem `json:"data" doc:"Список пород"`
}

// ============================================
// Get Breeds By Type
// ============================================

// GetBreedsByTypeInput представляет запрос на получение пород по типу питомца
type GetBreedsByTypeInput struct {
	PetTypePath
}

// GetBreedsByTypeOutput представляет ответ со списком пород
type GetBreedsByTypeOutput struct {
	Body BreedsList
}

// ============================================
// Get Blood Groups
// ============================================

// GetBloodGroupsInput представляет запрос на получение групп крови
type GetBloodGroupsInput struct {
	PetTypePath
}

// GetBloodGroupsOutput представляет ответ со списком групп крови
type GetBloodGroupsOutput struct {
	Body BloodGroupsList
}

// BloodGroupsList представляет список групп крови
type BloodGroupsList struct {
	Data []ReferenceItem `json:"data" doc:"Список групп крови"`
}

// ============================================
// Get Blood Components
// ============================================

// GetBloodComponentsInput представляет запрос на получение компонентов крови
type GetBloodComponentsInput struct{}

// GetBloodComponentsOutput представляет ответ со списком компонентов крови
type GetBloodComponentsOutput struct {
	Body BloodComponentsList
}

// BloodComponentsList представляет список компонентов крови
type BloodComponentsList struct {
	Data []ReferenceItem `json:"data" doc:"Список компонентов крови"`
}

// ============================================
// Get Locations
// ============================================

// GetLocationsInput представляет запрос на получение списка локаций
type GetLocationsInput struct{}

// GetLocationsOutput представляет ответ со списком локаций
type GetLocationsOutput struct {
	Body LocationsList
}

// LocationsList представляет список локаций
type LocationsList struct {
	Data []ReferenceItem `json:"data" doc:"Список локаций"`
}

// ============================================
// Get Donor Restrictions
// ============================================

// GetDonorRestrictionsInput представляет запрос на получение ограничений для доноров
type GetDonorRestrictionsInput struct{}

// GetDonorRestrictionsOutput представляет ответ с ограничениями для доноров
type GetDonorRestrictionsOutput struct {
	Body DonorRestrictionsDetail
}

// DonorRestrictionsDetail представляет детали ограничений для доноров
type DonorRestrictionsDetail struct {
	StopFactors []FactorDescription `json:"stopFactors" doc:"Список стоп-факторов"`
	WarnFactors []FactorDescription `json:"warnFactors" doc:"Список варн-факторов"`
}

// ============================================
// Common Types
// ============================================

// ReferenceItem представляет элемент справочника
type ReferenceItem struct {
	Value string `json:"value" doc:"Значение элемента" example:"labrador"`
	Label string `json:"label" doc:"Отображаемое название" example:"Лабрадор"`
}

// FactorDescription представляет описание фактора
type FactorDescription struct {
	Code           string `json:"code" doc:"Код фактора" example:"StopFactorTooOld"`
	Description    string `json:"description" doc:"Описание фактора" example:"Питомец слишком стар для донации"`
	SubDescription string `json:"subDescription,omitempty" doc:"Дополнительное описание"`
}
