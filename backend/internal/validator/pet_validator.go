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
