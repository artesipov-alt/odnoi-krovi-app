package validator

import (
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/ent"
)

var DefaultStopChecks = []CheckFunc{
	CheckPhoto,
	CheckNoRabiesVaccination,
	CheckNoInfectionVaccination,
	CheckVaccinationExpired,
	CheckVaccinationTooRecent,
	CheckNoDeworming,
	CheckNoEctoparasiteTreatment,
	CheckDewormingExpired,
	CheckEctoparasiteTreatmentExpired,
	CheckStopAge,
	CheckReproductiveStatus,
	CheckStopHealth,
	CheckDonationHistory,
	CheckPetStatus,
}

var DefaultWarnChecks = []CheckFunc{
	CheckWarnAge,
	CheckWarnHealth,
	CheckWarnLivingCondition,
	CheckWarnAnalyses,
}

// CheckPhoto проверяет наличие фотографий питомца. Если фото отсутствуют, возвращает стоп-фактор NO_PHOTO.
func CheckPhoto(p *ent.Pet) FactorCode {
	if len(p.PhotoUrls) == 0 {
		return StopFactorNoPhoto
	}
	return ""
}

// CheckNoInfectionVaccination проверяет наличие вакцинации от инфекций.
func CheckNoInfectionVaccination(p *ent.Pet) FactorCode {
	if p.Edges.Treatments == nil || p.Edges.Treatments.InfectionVaccinationDate == nil {
		return StopFactorNoInfectionVaccination
	}
	return ""
}

// CheckNoRabiesVaccination проверяет наличие вакцинации от бешенства.
func CheckNoRabiesVaccination(p *ent.Pet) FactorCode {
	if p.Edges.Treatments == nil || p.Edges.Treatments.RabiesVaccinationDate == nil {
		return StopFactorNoRabiesVaccination
	}
	return ""
}

// CheckVaccinationExpired проверяет, не просрочена ли вакцинация (более 1 года).
func CheckVaccinationExpired(p *ent.Pet) FactorCode {
	if p.Edges.Treatments == nil {
		return ""
	}
	t := p.Edges.Treatments
	now := time.Now()
	yearAgo := now.AddDate(-1, 0, 0)

	if t.RabiesVaccinationDate != nil && t.RabiesVaccinationDate.Before(yearAgo) {
		return StopFactorVaccinationExpired
	}
	if t.InfectionVaccinationDate != nil && t.InfectionVaccinationDate.Before(yearAgo) {
		return StopFactorVaccinationExpired
	}
	return ""
}

// CheckVaccinationTooRecent проверяет, не сделана ли вакцинация слишком недавно (менее 1 месяца).
func CheckVaccinationTooRecent(p *ent.Pet) FactorCode {
	if p.Edges.Treatments == nil {
		return ""
	}
	t := p.Edges.Treatments
	now := time.Now()
	monthAgo := now.AddDate(0, -1, 0)

	if t.RabiesVaccinationDate != nil && t.RabiesVaccinationDate.After(monthAgo) {
		return StopFactorVaccinationTooRecent
	}
	if t.InfectionVaccinationDate != nil && t.InfectionVaccinationDate.After(monthAgo) {
		return StopFactorVaccinationTooRecent
	}
	return ""
}

// CheckNoDeworming проверяет, была ли проведена дегельминтизация.
func CheckNoDeworming(p *ent.Pet) FactorCode {
	if p.Edges.Treatments == nil || p.Edges.Treatments.DewormingDate == nil {
		return StopFactorNoDeworming
	}
	return ""
}

// CheckNoEctoparasiteTreatment проверяет, была ли обработка от эктопаразитов.
func CheckNoEctoparasiteTreatment(p *ent.Pet) FactorCode {
	if p.Edges.Treatments == nil || p.Edges.Treatments.EctoparasiteTreatmentDate == nil {
		return StopFactorNoEctoparasiteTreatment
	}
	return ""
}

// CheckDewormingExpired проверяет, не просрочена ли дегельминтизация (более 3 месяцев).
func CheckDewormingExpired(p *ent.Pet) FactorCode {
	if p.Edges.Treatments == nil || p.Edges.Treatments.DewormingDate == nil {
		return ""
	}
	now := time.Now()
	threeMonthsAgo := now.AddDate(0, -3, 0)
	if p.Edges.Treatments.DewormingDate.Before(threeMonthsAgo) {
		return StopFactorDewormingExpired
	}
	return ""
}

// CheckEctoparasiteTreatmentExpired проверяет, не просрочена ли обработка от эктопаразитов (более 3 месяцев).
func CheckEctoparasiteTreatmentExpired(p *ent.Pet) FactorCode {
	if p.Edges.Treatments == nil || p.Edges.Treatments.EctoparasiteTreatmentDate == nil {
		return ""
	}
	now := time.Now()
	threeMonthsAgo := now.AddDate(0, -3, 0)
	if p.Edges.Treatments.EctoparasiteTreatmentDate.Before(threeMonthsAgo) {
		return StopFactorEctoparasiteTreatmentExpired
	}
	return ""
}

// CheckStopAge проверяет возраст питомца для стоп-факторов. Возвращает TOO_YOUNG, если возраст менее 1 года, или TOO_OLD, если более 8 лет.
func CheckStopAge(p *ent.Pet) FactorCode {
	if p.BirthDate == nil {
		return ""
	}
	now := time.Now()

	if p.BirthDate.After(now.AddDate(-1, 0, 0)) {
		return StopFactorTooYoung
	}
	if p.BirthDate.Before(now.AddDate(-8, 0, 0)) {
		return StopFactorTooOld
	}
	return ""
}

// CheckWarnAge проверяет возраст питомца для предупреждений. Возвращает APPROACHING_8_YEARS, если возраст от 7 до 8 лет.
func CheckWarnAge(p *ent.Pet) FactorCode {
	if p.BirthDate == nil {
		return ""
	}
	now := time.Now()

	if p.BirthDate.Before(now.AddDate(-7, 0, 0)) && !p.BirthDate.Before(now.AddDate(-8, 0, 0)) {
		return WarnFactorApproaching8Years
	}
	return ""
}

// CheckReproductiveStatus проверяет репродуктивный статус питомца. Возвращает стоп-фактор для беременности, лактации или течки.
func CheckReproductiveStatus(p *ent.Pet) FactorCode {
	if p.ReproductiveStatus == "" {
		return ""
	}
	switch p.ReproductiveStatus {
	case "pregnancy":
		return StopFactorPregnancy
	case "lactation":
		return StopFactorLactation
	case "estrus":
		return StopFactorEstrus
	}
	return ""
}

// CheckStopHealth проверяет здоровье питомца для стоп-факторов. Возвращает HAS_DISEASES, если статус не здоровый, или TRANSFUSED, если питомец был перелит.
func CheckStopHealth(p *ent.Pet) FactorCode {
	if p.Edges.Health == nil {
		return ""
	}
	h := p.Edges.Health

	if h.HealthStatus != "" && h.HealthStatus != "healthy" {
		return StopFactorHasDiseases
	}
	if h.Transfused {
		return StopFactorTransfused
	}
	return ""
}

// CheckWarnHealth проверяет здоровье питомца для предупреждений. Возвращает TAKING_MEDICATIONS или SURGICAL_INTERVENTION, если есть соответствующие записи.
func CheckWarnHealth(p *ent.Pet) FactorCode {
	if p.Edges.Health == nil {
		return ""
	}
	h := p.Edges.Health

	if h.Medications != "" {
		return WarnFactorTakingMedications
	}
	if h.SurgicalInterventions != "" {
		return WarnFactorSurgicalIntervention
	}
	return ""
}

// CheckDonationHistory проверяет историю донаций питомца. Если последняя донация была менее 2 месяцев назад, возвращает стоп-фактор DONATION_TOO_RECENT.
func CheckDonationHistory(p *ent.Pet) FactorCode {
	if p.Edges.Health == nil || p.Edges.Health.LastDonation == nil {
		return ""
	}
	now := time.Now()
	twoMonthsAgo := now.AddDate(0, -2, 0)
	if p.Edges.Health.LastDonation.After(twoMonthsAgo) {
		return StopFactorDonationTooRecent
	}
	return ""
}

// CheckWarnLivingCondition проверяет условия проживания питомца. Если питомец на самовыгуле, возвращает предупреждение FREE_RANGE.
func CheckWarnLivingCondition(p *ent.Pet) FactorCode {
	if p.LivingCondition == "self_outdoor" {
		return WarnFactorFreeRange
	}
	return ""
}

// CheckWarnAnalyses проверяет анализы питомца. Если нет анализов или все старше 1 года, возвращает предупреждение NO_CURRENT_ANALYSES.
func CheckWarnAnalyses(p *ent.Pet) FactorCode {
	if len(p.Edges.Analyses) == 0 {
		return WarnFactorNoCurrentAnalyses
	}
	now := time.Now()
	hasRecent := false
	yearAgo := now.AddDate(-1, 0, 0)
	for _, a := range p.Edges.Analyses {
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

// CheckPetStatus проверяет статус питомца. Если статус "recipient", возвращает стоп-фактор CURRENTLY_RECIPIENT.
func CheckPetStatus(p *ent.Pet) FactorCode {
	if p.PetStatus == "recipient" {
		return StopFactorCurrentlyRecipient
	}
	return ""
}
