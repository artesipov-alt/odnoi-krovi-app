// Package mapper provides conversion functions between domain models and DTOs.
package mapper

import (
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/transport/http/dto"
)

// PetMapper handles conversions between domain Pet model and DTOs.
type PetMapper struct{}

// NewPetMapper creates a new PetMapper instance.
func NewPetMapper() *PetMapper {
	return &PetMapper{}
}

// ToDTO converts a domain Pet model to a DTO.
func (m *PetMapper) ToDTO(petmodel model.Pet) dto.Pet {
	petDTO := dto.Pet{
		ID:                 petmodel.ID,
		Name:               petmodel.Name,
		ChipNumber:         petmodel.ChipNumber,
		PhotoURLs:          petmodel.PhotoURLs,
		WeightKg:           petmodel.WeightKg,
		BirthDate:          petmodel.BirthDate,
		PetStatus:          string(petmodel.PetStatus),
		LivingCondition:    string(petmodel.LivingCondition),
		Gender:             string(petmodel.Gender),
		Type:               string(petmodel.Type),
		ReproductiveStatus: string(petmodel.ReproductiveStatus),
		Bonuses:            petmodel.Bonuses,
		CreatedAt:          petmodel.CreatedAt,
		UpdatedAt:          petmodel.UpdatedAt,
		DeletedAt:          petmodel.DeletedAt,
	}

	// Map StopFactors and WarnFactors into DonorRestrictions
	if len(petmodel.StopFactors) > 0 || len(petmodel.WarnFactors) > 0 {
		var stopFactors []dto.RestrictionFactor
		var warnFactors []dto.RestrictionFactor
		for _, code := range petmodel.StopFactors {
			desc := model.GetFactorDescription(model.FactorCode(code))
			factor := dto.RestrictionFactor{
				Code:           code,
				Description:    desc.Description,
				SubDescription: desc.SubDescription,
			}
			stopFactors = append(stopFactors, factor)
		}
		for _, code := range petmodel.WarnFactors {
			desc := model.GetFactorDescription(model.FactorCode(code))
			factor := dto.RestrictionFactor{
				Code:           code,
				Description:    desc.Description,
				SubDescription: desc.SubDescription,
			}
			warnFactors = append(warnFactors, factor)
		}
		petDTO.DonorRestrictions = &dto.DonorRestrictions{
			StopFactors: stopFactors,
			WarnFactors: warnFactors,
		}
	}

	if petmodel.BreedRefID != nil {
		petDTO.BreedID = *petmodel.BreedRefID
	}
	if petmodel.BloodGroupName != nil {
		petDTO.BloodGroup = *petmodel.BloodGroupName
	}

	if petmodel.Health != nil {
		healthStatus := string(petmodel.Health.HealthStatus)
		petDTO.Health = &dto.PetHealth{
			HealthStatus:          &healthStatus,
			LastDonation:          petmodel.Health.LastDonation,
			Transfused:            petmodel.Health.Transfused,
			Medications:           petmodel.Health.Medications,
			SurgicalInterventions: petmodel.Health.SurgicalInterventions,
		}
	}

	if petmodel.Treatments != nil {
		petDTO.Treatments = &dto.PetTreatment{
			RabiesVaccinationDate:     petmodel.Treatments.RabiesVaccinationDate,
			InfectionVaccinationDate:  petmodel.Treatments.InfectionVaccinationDate,
			EctoparasiteTreatmentDate: petmodel.Treatments.EctoparasiteTreatmentDate,
			DewormingDate:             petmodel.Treatments.DewormingDate,
		}
	}

	// Analyses mapping
	if petmodel.Analyses != nil {
		petDTO.Analyses = &dto.PetAnalysisGroup{}
		for _, a := range petmodel.Analyses {
			analysisName := string(a.AnalysisName)
			analysisType := string(a.AnalysisType)
			dtoAnalysis := &dto.PetAnalysis{
				ID:           nilable(a.ID),
				AnalysisName: &analysisName,
				AnalysisType: &analysisType,
				AnalysisDate: a.AnalysisDate,
			}
			switch a.AnalysisName {
			case "leukemia":
				petDTO.Analyses.Leukemia = append(petDTO.Analyses.Leukemia, dtoAnalysis)
			case "immunodeficiency":
				petDTO.Analyses.Immunodeficiency = append(petDTO.Analyses.Immunodeficiency, dtoAnalysis)
			case "hemoplasmosis":
				petDTO.Analyses.Hemoplasmosis = append(petDTO.Analyses.Hemoplasmosis, dtoAnalysis)
			case "bartonellosis":
				petDTO.Analyses.Bartonellosis = append(petDTO.Analyses.Bartonellosis, dtoAnalysis)
			case "babesiosis":
				petDTO.Analyses.Babesiosis = append(petDTO.Analyses.Babesiosis, dtoAnalysis)
			case "dirofilaria":
				petDTO.Analyses.Dirofilaria = append(petDTO.Analyses.Dirofilaria, dtoAnalysis)
			case "ehrlichiosis":
				petDTO.Analyses.Ehrlichiosis = append(petDTO.Analyses.Ehrlichiosis, dtoAnalysis)
			case "anaplasmosis":
				petDTO.Analyses.Anaplasmosis = append(petDTO.Analyses.Anaplasmosis, dtoAnalysis)
			}
		}
	}

	return petDTO
}

// ToDTOs converts a slice of domain Pet models to DTOs.
func (m *PetMapper) ToDTOs(pets []*model.Pet) []dto.Pet {
	if pets == nil {
		return nil
	}
	petDTOs := make([]dto.Pet, len(pets))
	for i, p := range pets {
		if p != nil {
			petDTOs[i] = m.ToDTO(*p)
		}
	}
	return petDTOs
}

// ToDomain converts a PetCreate DTO to a domain Pet model.
func (m *PetMapper) ToDomainCreate(petDto dto.PetCreate) *model.Pet {
	petmodel := &model.Pet{
		Name:            petDto.Name,
		Type:            model.PetType(petDto.Type),
		WeightKg:        petDto.WeightKg,
		Gender:          model.Gender(petDto.Gender),
		ChipNumber:      petDto.ChipNumber,
		LivingCondition: model.LivingCondition(petDto.LivingCondition),
	}

	if petDto.BreedID != "" {
		petmodel.BreedRefID = &petDto.BreedID
	}
	if petDto.BloodGroup != "" {
		petmodel.BloodGroupName = &petDto.BloodGroup
	}

	if petDto.BirthDate != nil {
		petmodel.BirthDate = petDto.BirthDate
	}

	if petDto.ReproductiveStatus != "" {
		petmodel.ReproductiveStatus = model.ReproductiveStatus(petDto.ReproductiveStatus)
	}
	if petDto.AgeMonths != 0 || petDto.AgeYears != 0 {
		petmodel.BirthDate = calculateBirthDateFromAge(&petDto.AgeYears, &petDto.AgeMonths)
	}

	// Handle PetHealth
	if petDto.Health != nil {
		healthmodel := &model.PetHealth{}
		if petDto.Health.HealthStatus != nil {
			healthmodel.HealthStatus = model.HealthStatus(*petDto.Health.HealthStatus)
		}
		if petDto.Health.LastDonation != nil {
			healthmodel.LastDonation = petDto.Health.LastDonation
		}
		if petDto.Health.Transfused != nil {
			healthmodel.Transfused = petDto.Health.Transfused
		}
		if petDto.Health.Medications != nil {
			healthmodel.Medications = petDto.Health.Medications
		}
		if petDto.Health.SurgicalInterventions != nil {
			healthmodel.SurgicalInterventions = petDto.Health.SurgicalInterventions
		}
		petmodel.Health = healthmodel
	}

	// Handle PetTreatment
	if petDto.Treatments != nil {
		treatmentmodel := &model.PetTreatment{}
		if petDto.Treatments.RabiesVaccinationDate != nil {
			treatmentmodel.RabiesVaccinationDate = petDto.Treatments.RabiesVaccinationDate
		}
		if petDto.Treatments.InfectionVaccinationDate != nil {
			treatmentmodel.InfectionVaccinationDate = petDto.Treatments.InfectionVaccinationDate
		}
		if petDto.Treatments.EctoparasiteTreatmentDate != nil {
			treatmentmodel.EctoparasiteTreatmentDate = petDto.Treatments.EctoparasiteTreatmentDate
		}
		if petDto.Treatments.DewormingDate != nil {
			treatmentmodel.DewormingDate = petDto.Treatments.DewormingDate
		}
		petmodel.Treatments = treatmentmodel
	}

	// Handle PetAnalysis
	analysesDomain := []*model.PetAnalysis{}
	if petDto.Analyses != nil {
		processGroup := func(group []*dto.PetAnalysis, name string) {
			for _, a := range group {
				if a != nil && a.AnalysisDate != nil {
					analysis := &model.PetAnalysis{
						AnalysisName: name,
						AnalysisDate: a.AnalysisDate,
					}
					if a.AnalysisType != nil {
						analysis.AnalysisType = *a.AnalysisType
					}
					analysesDomain = append(analysesDomain, analysis)
				}
			}
		}

		processGroup(petDto.Analyses.Leukemia, "leukemia")
		processGroup(petDto.Analyses.Immunodeficiency, "immunodeficiency")
		processGroup(petDto.Analyses.Hemoplasmosis, "hemoplasmosis")
		processGroup(petDto.Analyses.Bartonellosis, "bartonellosis")
		processGroup(petDto.Analyses.Babesiosis, "babesiosis")
		processGroup(petDto.Analyses.Dirofilaria, "dirofilaria")
		processGroup(petDto.Analyses.Ehrlichiosis, "ehrlichiosis")
		processGroup(petDto.Analyses.Anaplasmosis, "anaplasmosis")
		petmodel.Analyses = analysesDomain
	}

	return petmodel
}

// ToDomainUpdate applies PetUpdate DTO fields to a domain Pet model.
func (m *PetMapper) ToDomainUpdate(petDto dto.PetUpdate, petmodel *model.Pet) {
	if petmodel == nil {
		return
	}

	if petDto.Name != nil {
		petmodel.Name = *petDto.Name
	}
	if petDto.Type != nil {
		petmodel.Type = model.PetType(*petDto.Type)
	}
	if petDto.WeightKg != nil {
		petmodel.WeightKg = *petDto.WeightKg
	}
	if petDto.Gender != nil {
		petmodel.Gender = model.Gender(*petDto.Gender)
	}
	if petDto.ChipNumber != nil {
		petmodel.ChipNumber = *petDto.ChipNumber
	}
	if petDto.LivingCondition != nil {
		petmodel.LivingCondition = model.LivingCondition(*petDto.LivingCondition)
	}
	if petDto.BreedID != nil {
		petmodel.BreedRefID = petDto.BreedID
	}
	if petDto.BloodGroup != nil {
		petmodel.BloodGroupName = petDto.BloodGroup
	}
	if petDto.BirthDate != nil {
		petmodel.BirthDate = petDto.BirthDate
	}
	if petDto.ReproductiveStatus != nil {
		petmodel.ReproductiveStatus = model.ReproductiveStatus(*petDto.ReproductiveStatus)
	}
	if petDto.Bonuses != nil {
		petmodel.Bonuses = *petDto.Bonuses
	}
	if petDto.AgeMonths != nil || petDto.AgeYears != nil {
		petmodel.BirthDate = calculateBirthDateFromAge(petDto.AgeYears, petDto.AgeMonths)
	}

	// Handle PetHealth
	if petDto.Health != nil {
		healthmodel := &model.PetHealth{}
		if petDto.Health.HealthStatus != nil {
			healthmodel.HealthStatus = model.HealthStatus(*petDto.Health.HealthStatus)
		}
		if petDto.Health.LastDonation != nil {
			healthmodel.LastDonation = petDto.Health.LastDonation
		}
		if petDto.Health.Transfused != nil {
			healthmodel.Transfused = petDto.Health.Transfused
		}
		if petDto.Health.Medications != nil {
			healthmodel.Medications = petDto.Health.Medications
		}
		if petDto.Health.SurgicalInterventions != nil {
			healthmodel.SurgicalInterventions = petDto.Health.SurgicalInterventions
		}
		petmodel.Health = healthmodel
	}

	// Handle PetTreatment
	if petDto.Treatments != nil {
		treatmentmodel := &model.PetTreatment{}
		if petDto.Treatments.RabiesVaccinationDate != nil {
			treatmentmodel.RabiesVaccinationDate = petDto.Treatments.RabiesVaccinationDate
		}
		if petDto.Treatments.InfectionVaccinationDate != nil {
			treatmentmodel.InfectionVaccinationDate = petDto.Treatments.InfectionVaccinationDate
		}
		if petDto.Treatments.EctoparasiteTreatmentDate != nil {
			treatmentmodel.EctoparasiteTreatmentDate = petDto.Treatments.EctoparasiteTreatmentDate
		}
		if petDto.Treatments.DewormingDate != nil {
			treatmentmodel.DewormingDate = petDto.Treatments.DewormingDate
		}
		petmodel.Treatments = treatmentmodel
	}

	// Handle PetAnalysis
	analysesDomain := []*model.PetAnalysis{}
	if petDto.Analyses != nil {
		processGroup := func(group []*dto.PetAnalysis, name string) {
			for _, a := range group {
				if a != nil && a.AnalysisDate != nil {
					analysis := &model.PetAnalysis{
						AnalysisName: name,
						AnalysisDate: a.AnalysisDate,
					}
					if a.AnalysisType != nil {
						analysis.AnalysisType = *a.AnalysisType
					}
					analysesDomain = append(analysesDomain, analysis)
				}
			}
		}

		processGroup(petDto.Analyses.Leukemia, "leukemia")
		processGroup(petDto.Analyses.Immunodeficiency, "immunodeficiency")
		processGroup(petDto.Analyses.Hemoplasmosis, "hemoplasmosis")
		processGroup(petDto.Analyses.Bartonellosis, "bartonellosis")
		processGroup(petDto.Analyses.Babesiosis, "babesiosis")
		processGroup(petDto.Analyses.Dirofilaria, "dirofilaria")
		processGroup(petDto.Analyses.Ehrlichiosis, "ehrlichiosis")
		processGroup(petDto.Analyses.Anaplasmosis, "anaplasmosis")
		petmodel.Analyses = analysesDomain
	}
}

// ToSimplifiedDTO converts a domain Pet model to a simplified DTO for nested usage.
// This is used when embedding pet info in other DTOs (e.g., UserDTO).
func (m *PetMapper) ToSimplifiedDTO(pet *model.Pet) dto.Pet {
	if pet == nil {
		return dto.Pet{}
	}

	dtoPet := dto.Pet{
		ID:         pet.ID,
		Name:       pet.Name,
		ChipNumber: pet.ChipNumber,
		PhotoURLs:  pet.PhotoURLs,
		WeightKg:   pet.WeightKg,
		BirthDate:  pet.BirthDate,
		CreatedAt:  pet.CreatedAt,
		UpdatedAt:  pet.UpdatedAt,
	}

	if pet.BreedRefID != nil {
		dtoPet.BreedID = *pet.BreedRefID
	}

	if pet.BloodGroupName != nil {
		dtoPet.BloodGroup = *pet.BloodGroupName
	}

	if pet.LivingCondition != "" {
		dtoPet.LivingCondition = string(pet.LivingCondition)
	}

	if pet.Gender != "" {
		dtoPet.Gender = string(pet.Gender)
	}

	if pet.Type != "" {
		dtoPet.Type = string(pet.Type)
	}

	return dtoPet
}

// ToSimplifiedDTOs converts a slice of domain Pet models to simplified DTOs.
func (m *PetMapper) ToSimplifiedDTOs(pets []*model.Pet) []dto.Pet {
	if pets == nil {
		return nil
	}
	dtoPets := make([]dto.Pet, len(pets))
	for i, pet := range pets {
		dtoPets[i] = m.ToSimplifiedDTO(pet)
	}
	return dtoPets
}

// calculateBirthDateFromAge вычисляет дату рождения на основе возраста в годах и месяцах.
func calculateBirthDateFromAge(ageYears, ageMonths *int) *time.Time {
	if ageYears == nil && ageMonths == nil {
		return nil
	}

	years := 0
	if ageYears != nil {
		years = *ageYears
	}

	months := 0
	if ageMonths != nil {
		months = *ageMonths
	}

	if years > 0 || months > 0 {
		birthDate := time.Now().AddDate(-years, -months, 0)
		birthDate = time.Date(birthDate.Year(), birthDate.Month(), 1, 0, 0, 0, 0, time.UTC)
		return &birthDate
	}
	return nil
}

// nilable returns a pointer to the value if it's not the zero value for its type.
func nilable[T any](val T) *T {
	var zero T
	if any(val) == any(zero) {
		return nil
	}
	return &val
}
