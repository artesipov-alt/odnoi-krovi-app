package validator

import (
	"github.com/artesipov-alt/odnoi-krovi-app/ent"
)

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

type DonorValidator interface {
	GetStopFactors(p *ent.Pet) []FactorCode
	GetWarnFactors(p *ent.Pet) []FactorCode
}

// CheckFunc — общая сигнатура функции проверки
type CheckFunc func(p *ent.Pet) FactorCode

type PetValidatorImpl struct {
	stopChecks []CheckFunc
	warnChecks []CheckFunc
}

// NewDonorValidator — конструктор, который принимает набор проверок
func NewDonorValidator(stopChecks []CheckFunc, warnChecks []CheckFunc) *PetValidatorImpl {
	return &PetValidatorImpl{
		stopChecks: stopChecks,
		warnChecks: warnChecks,
	}
}

// GetStopFactors прогоняет питомца по стоп-проверкам
func (v *PetValidatorImpl) GetStopFactors(p *ent.Pet) []FactorCode {
	var factors []FactorCode
	for _, check := range v.stopChecks {
		if code := check(p); code != "" {
			factors = append(factors, code)
		}
	}
	return factors
}

// GetWarnFactors прогоняет питомца по варн-проверкам
func (v *PetValidatorImpl) GetWarnFactors(p *ent.Pet) []FactorCode {
	var factors []FactorCode
	for _, check := range v.warnChecks {
		if code := check(p); code != "" {
			factors = append(factors, code)
		}
	}
	return factors
}
