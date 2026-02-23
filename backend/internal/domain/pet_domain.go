package domain

import (
	"time"
)

// PetStatus представляет статус питомца
type PetStatus string

const (
	PetStatusNone       PetStatus = "none"
	PetStatusDonor      PetStatus = "donor"
	PetStatusRecipient  PetStatus = "recipient"
	PetStatusBloodFound PetStatus = "blood_found"
)

// PetType представляет тип животного
type PetType string

const (
	PetTypeDog PetType = "dog"
	PetTypeCat PetType = "cat"
)

// Gender представляет пол питомца
type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
)

// LivingCondition представляет условия проживания
type LivingCondition string

const (
	LivingConditionIndoor       LivingCondition = "indoor"
	LivingConditionLeashWalking LivingCondition = "leash_walking"
	LivingConditionSelfOutdoor  LivingCondition = "self_outdoor"
)

// ReproductiveStatus представляет репродуктивный статус
type ReproductiveStatus string

const (
	ReproductiveStatusPregnancy ReproductiveStatus = "pregnancy"
	ReproductiveStatusLactation ReproductiveStatus = "lactation"
	ReproductiveStatusEstrus    ReproductiveStatus = "estrus"
)

// HealthStatus представляет состояние здоровья
type HealthStatus string

const (
	HealthStatusHealthy HealthStatus = "healthy"
	HealthStatusIll     HealthStatus = "ill"
	HealthStatusUnknown HealthStatus = "unknown"
)

// Pet представляет доменную модель питомца
type Pet struct {
	ID                 string
	Name               string
	PetStatus          PetStatus
	Type               PetType
	WeightKg           float64
	Gender             Gender
	BirthDate          *time.Time
	ChipNumber         string
	PhotoURLs          []string
	LivingCondition    LivingCondition
	ReproductiveStatus ReproductiveStatus
	OwnerID            string
	BreedRefID         *string
	BloodGroupName     *string
	DonorRestrictions  []string
	Bonuses            []string
	Health             *PetHealth
	Treatments         *PetTreatment
	Analyses           []*PetAnalysis
	CreatedAt          *time.Time
	UpdatedAt          *time.Time
	DeletedAt          *time.Time
}

// PetHealth представляет здоровье питомца
type PetHealth struct {
	HealthStatus          HealthStatus
	LastDonation          *time.Time
	Transfused            *bool
	Medications           *string
	SurgicalInterventions *string
}

// PetTreatment представляет лечение питомца
type PetTreatment struct {
	RabiesVaccinationDate     *time.Time
	InfectionVaccinationDate  *time.Time
	EctoparasiteTreatmentDate *time.Time
	DewormingDate             *time.Time
}

// PetAnalysis представляет анализ питомца
type PetAnalysis struct {
	ID           string
	AnalysisName string
	AnalysisType string
	AnalysisDate *time.Time
}

// FactorCode — общий тип-код для факторов и предупреждений
type FactorCode string

const (
	StopFactorTooOld                       FactorCode = "STOP_TOO_OLD"
	StopFactorNoPhoto                      FactorCode = "STOP_NO_PHOTO"
	StopFactorNoInfectionVaccination       FactorCode = "STOP_NO_INFECTION_VACCINATION"
	StopFactorNoRabiesVaccination          FactorCode = "STOP_NO_RABIES_VACCINATION"
	StopFactorVaccinationExpired           FactorCode = "STOP_VACCINATION_EXPIRED"
	StopFactorVaccinationTooRecent         FactorCode = "STOP_VACCINATION_TOO_RECENT"
	StopFactorEctoparasiteTreatmentExpired FactorCode = "STOP_ECTOPARASITE_TREATMENT_EXPIRED"
	StopFactorNoDeworming                  FactorCode = "STOP_NO_DEWORMING"
	StopFactorNoEctoparasiteTreatment      FactorCode = "STOP_NO_ECTOPARASITE_TREATMENT"
	StopFactorDewormingExpired             FactorCode = "STOP_DEWORMING_EXPIRED"
	StopFactorTooYoung                     FactorCode = "STOP_TOO_YOUNG"
	StopFactorPregnancy                    FactorCode = "STOP_PREGNANCY"
	StopFactorLactation                    FactorCode = "STOP_LACTATION"
	StopFactorEstrus                       FactorCode = "STOP_ESTRUS"
	StopFactorHasDiseases                  FactorCode = "STOP_HAS_DISEASES"
	StopFactorDonationTooRecent            FactorCode = "STOP_DONATION_TOO_RECENT"
	StopFactorTransfused                   FactorCode = "STOP_TRANSFUSED"
	StopFactorCurrentlyRecipient           FactorCode = "STOP_CURRENTLY_RECIPIENT"

	WarnFactorTakingMedications    FactorCode = "WARN_TAKING_MEDICATIONS"
	WarnFactorSurgicalIntervention FactorCode = "WARN_SURGICAL_INTERVENTION"
	WarnFactorApproaching8Years    FactorCode = "WARN_APPROACHING_8_YEARS"
	WarnFactorFreeRange            FactorCode = "WARN_FREE_RANGE"
	WarnFactorNoCurrentAnalyses    FactorCode = "WARN_NO_CURRENT_ANALYSES"
)

// FactorDescription представляет описание фактора
type FactorDescription struct {
	Description    string
	SubDescription string
}

// factorDescriptions маппа кодов факторов в их описания
var factorDescriptions = map[FactorCode]FactorDescription{
	StopFactorTooOld: {
		Description:    "Возраст больше 8 лет",
		SubDescription: "Донации после 8 лет рискованны для донора",
	},
	StopFactorNoPhoto: {
		Description:    "Отсутствует фото питомца",
		SubDescription: "",
	},
	StopFactorNoInfectionVaccination: {
		Description:    "Отсутствует вакцинация от инфекций",
		SubDescription: "",
	},
	StopFactorNoRabiesVaccination: {
		Description:    "Отсутствует вакцинация от бешенства",
		SubDescription: "",
	},
	StopFactorVaccinationExpired: {
		Description:    "Прошло больше года после вакцинации",
		SubDescription: "",
	},
	StopFactorVaccinationTooRecent: {
		Description:    "Прошло меньше месяца после вакцинации",
		SubDescription: "",
	},
	StopFactorEctoparasiteTreatmentExpired: {
		Description:    "Прошло больше 3 месяцев после обработки от эктопаразитов",
		SubDescription: "",
	},
	StopFactorNoDeworming: {
		Description:    "Не проведена дегельминтизация",
		SubDescription: "",
	},
	StopFactorNoEctoparasiteTreatment: {
		Description:    "Не обработан от эктопаразитов",
		SubDescription: "",
	},
	StopFactorDewormingExpired: {
		Description:    "Прошло больше 3 месяцев после дегельминтизации",
		SubDescription: "",
	},
	StopFactorTooYoung: {
		Description:    "Возраст меньше года",
		SubDescription: "",
	},
	StopFactorPregnancy: {
		Description:    "Беременность",
		SubDescription: "",
	},
	StopFactorLactation: {
		Description:    "Лактация",
		SubDescription: "",
	},
	StopFactorEstrus: {
		Description:    "Течка",
		SubDescription: "",
	},
	StopFactorHasDiseases: {
		Description:    "Есть заболевания",
		SubDescription: "",
	},
	StopFactorDonationTooRecent: {
		Description:    "Прошло меньше 2 месяцев с последней донации",
		SubDescription: "",
	},
	StopFactorTransfused: {
		Description:    "Питомцу переливали кровь",
		SubDescription: "",
	},
	StopFactorCurrentlyRecipient: {
		Description:    "Питомцу сейчас ищут кровь",
		SubDescription: "",
	},
	WarnFactorTakingMedications: {
		Description:    "Идет прием препаратов",
		SubDescription: "Прием препаратов может говорить о проблемах со здоровьем",
	},
	WarnFactorSurgicalIntervention: {
		Description:    "Было хирургическое вмешательство",
		SubDescription: "Донор может еще восстанавливаться после операции",
	},
	WarnFactorApproaching8Years: {
		Description:    "Скоро исполнится 8 лет",
		SubDescription: "Донации после 8 лет рискованны для донора",
	},
	WarnFactorFreeRange: {
		Description:    "Животное на самовыгуле",
		SubDescription: "",
	},
	WarnFactorNoCurrentAnalyses: {
		Description:    "Отсутствуют актуальные анализы",
		SubDescription: "Рекомендованы проверки раз в год на особо опасные инфекции - их можно сдать перед донацией (в некоторых клиниках за 1 день)",
	},
}

// GetFactorDescription возвращает описание для данного кода фактора
func GetFactorDescription(code FactorCode) FactorDescription {
	return factorDescriptions[code]
}

// GetStopFactors возвращает список стоп-факторов для питомца на основе текущего времени
func (p *Pet) GetStopFactors(now time.Time) []FactorCode {
	var factors []FactorCode
	if code := p.checkPhoto(); code != "" {
		factors = append(factors, code)
	}
	if code := p.checkNoInfectionVaccination(); code != "" {
		factors = append(factors, code)
	}
	if code := p.checkNoRabiesVaccination(); code != "" {
		factors = append(factors, code)
	}
	if code := p.checkVaccinationExpired(now); code != "" {
		factors = append(factors, code)
	}
	if code := p.checkVaccinationTooRecent(now); code != "" {
		factors = append(factors, code)
	}
	if code := p.checkNoDeworming(); code != "" {
		factors = append(factors, code)
	}
	if code := p.checkNoEctoparasiteTreatment(); code != "" {
		factors = append(factors, code)
	}
	if code := p.checkDewormingExpired(now); code != "" {
		factors = append(factors, code)
	}
	if code := p.checkEctoparasiteTreatmentExpired(now); code != "" {
		factors = append(factors, code)
	}
	if code := p.checkStopAge(now); code != "" {
		factors = append(factors, code)
	}
	if code := p.checkReproductiveStatus(); code != "" {
		factors = append(factors, code)
	}
	if code := p.checkStopHealth(); code != "" {
		factors = append(factors, code)
	}
	if code := p.checkDonationHistory(now); code != "" {
		factors = append(factors, code)
	}
	return factors
}

// GetWarnFactors возвращает список варн-факторов для питомца на основе текущего времени
func (p *Pet) GetWarnFactors(now time.Time) []FactorCode {
	var factors []FactorCode
	if code := p.checkWarnAge(now); code != "" {
		factors = append(factors, code)
	}
	if code := p.checkWarnHealth(); code != "" {
		factors = append(factors, code)
	}
	if code := p.checkWarnLivingCondition(); code != "" {
		factors = append(factors, code)
	}
	if code := p.checkWarnAnalyses(now); code != "" {
		factors = append(factors, code)
	}
	return factors
}

// checkPhoto проверяет наличие фотографий питомца
func (p *Pet) checkPhoto() FactorCode {
	if len(p.PhotoURLs) == 0 {
		return StopFactorNoPhoto
	}
	return ""
}

// checkNoInfectionVaccination проверяет наличие вакцинации от инфекций
func (p *Pet) checkNoInfectionVaccination() FactorCode {
	if p.Treatments == nil || p.Treatments.InfectionVaccinationDate == nil {
		return StopFactorNoInfectionVaccination
	}
	return ""
}

// checkNoRabiesVaccination проверяет наличие вакцинации от бешенства
func (p *Pet) checkNoRabiesVaccination() FactorCode {
	if p.Treatments == nil || p.Treatments.RabiesVaccinationDate == nil {
		return StopFactorNoRabiesVaccination
	}
	return ""
}

// checkVaccinationExpired проверяет, не просрочена ли вакцинация
func (p *Pet) checkVaccinationExpired(now time.Time) FactorCode {
	if p.Treatments == nil {
		return ""
	}
	yearAgo := now.AddDate(-1, 0, 0)
	if p.Treatments.RabiesVaccinationDate != nil && p.Treatments.RabiesVaccinationDate.Before(yearAgo) {
		return StopFactorVaccinationExpired
	}
	if p.Treatments.InfectionVaccinationDate != nil && p.Treatments.InfectionVaccinationDate.Before(yearAgo) {
		return StopFactorVaccinationExpired
	}
	return ""
}

// checkVaccinationTooRecent проверяет, не сделана ли вакцинация слишком недавно
func (p *Pet) checkVaccinationTooRecent(now time.Time) FactorCode {
	if p.Treatments == nil {
		return ""
	}
	monthAgo := now.AddDate(0, -1, 0)
	if p.Treatments.RabiesVaccinationDate != nil && p.Treatments.RabiesVaccinationDate.After(monthAgo) {
		return StopFactorVaccinationTooRecent
	}
	if p.Treatments.InfectionVaccinationDate != nil && p.Treatments.InfectionVaccinationDate.After(monthAgo) {
		return StopFactorVaccinationTooRecent
	}
	return ""
}

// checkNoDeworming проверяет, была ли проведена дегельминтизация
func (p *Pet) checkNoDeworming() FactorCode {
	if p.Treatments == nil || p.Treatments.DewormingDate == nil {
		return StopFactorNoDeworming
	}
	return ""
}

// checkNoEctoparasiteTreatment проверяет, была ли обработка от эктопаразитов
func (p *Pet) checkNoEctoparasiteTreatment() FactorCode {
	if p.Treatments == nil || p.Treatments.EctoparasiteTreatmentDate == nil {
		return StopFactorNoEctoparasiteTreatment
	}
	return ""
}

// checkDewormingExpired проверяет, не просрочена ли дегельминтизация
func (p *Pet) checkDewormingExpired(now time.Time) FactorCode {
	if p.Treatments == nil || p.Treatments.DewormingDate == nil {
		return ""
	}
	threeMonthsAgo := now.AddDate(0, -3, 0)
	if p.Treatments.DewormingDate.Before(threeMonthsAgo) {
		return StopFactorDewormingExpired
	}
	return ""
}

// checkEctoparasiteTreatmentExpired проверяет, не просрочена ли обработка от эктопаразитов
func (p *Pet) checkEctoparasiteTreatmentExpired(now time.Time) FactorCode {
	if p.Treatments == nil || p.Treatments.EctoparasiteTreatmentDate == nil {
		return ""
	}
	threeMonthsAgo := now.AddDate(0, -3, 0)
	if p.Treatments.EctoparasiteTreatmentDate.Before(threeMonthsAgo) {
		return StopFactorEctoparasiteTreatmentExpired
	}
	return ""
}

// checkStopAge проверяет возраст для стоп-факторов
func (p *Pet) checkStopAge(now time.Time) FactorCode {
	if p.BirthDate == nil {
		return ""
	}
	if p.BirthDate.After(now.AddDate(-1, 0, 0)) {
		return StopFactorTooYoung
	}
	if p.BirthDate.Before(now.AddDate(-8, 0, 0)) {
		return StopFactorTooOld
	}
	return ""
}

// checkWarnAge проверяет возраст для предупреждений
func (p *Pet) checkWarnAge(now time.Time) FactorCode {
	if p.BirthDate == nil {
		return ""
	}
	if p.BirthDate.Before(now.AddDate(-7, 0, 0)) && !p.BirthDate.Before(now.AddDate(-8, 0, 0)) {
		return WarnFactorApproaching8Years
	}
	return ""
}

// checkReproductiveStatus проверяет репродуктивный статус
func (p *Pet) checkReproductiveStatus() FactorCode {
	switch p.ReproductiveStatus {
	case ReproductiveStatusPregnancy:
		return StopFactorPregnancy
	case ReproductiveStatusLactation:
		return StopFactorLactation
	case ReproductiveStatusEstrus:
		return StopFactorEstrus
	}
	return ""
}

// checkStopHealth проверяет здоровье для стоп-факторов
func (p *Pet) checkStopHealth() FactorCode {
	if p.Health == nil {
		return ""
	}
	if p.Health.HealthStatus != "" && p.Health.HealthStatus != HealthStatusHealthy {
		return StopFactorHasDiseases
	}
	if p.Health.Transfused != nil && *p.Health.Transfused {
		return StopFactorTransfused
	}
	return ""
}

// checkWarnHealth проверяет здоровье для предупреждений
func (p *Pet) checkWarnHealth() FactorCode {
	if p.Health == nil {
		return ""
	}
	if p.Health.Medications != nil && *p.Health.Medications != "" {
		return WarnFactorTakingMedications
	}
	if p.Health.SurgicalInterventions != nil && *p.Health.SurgicalInterventions != "" {
		return WarnFactorSurgicalIntervention
	}
	return ""
}

// checkDonationHistory проверяет историю донаций
func (p *Pet) checkDonationHistory(now time.Time) FactorCode {
	if p.Health == nil || p.Health.LastDonation == nil {
		return ""
	}
	twoMonthsAgo := now.AddDate(0, -2, 0)
	if p.Health.LastDonation.After(twoMonthsAgo) {
		return StopFactorDonationTooRecent
	}
	return ""
}

// checkWarnLivingCondition проверяет условия проживания
func (p *Pet) checkWarnLivingCondition() FactorCode {
	if p.LivingCondition == LivingConditionSelfOutdoor {
		return WarnFactorFreeRange
	}
	return ""
}

// checkWarnAnalyses проверяет анализы
func (p *Pet) checkWarnAnalyses(now time.Time) FactorCode {
	if len(p.Analyses) == 0 {
		return WarnFactorNoCurrentAnalyses
	}
	hasRecent := false
	yearAgo := now.AddDate(-1, 0, 0)
	for _, a := range p.Analyses {
		if a.AnalysisDate != nil && !a.AnalysisDate.Before(yearAgo) {
			hasRecent = true
			break
		}
	}
	if !hasRecent {
		return WarnFactorNoCurrentAnalyses
	}
	return ""
}
