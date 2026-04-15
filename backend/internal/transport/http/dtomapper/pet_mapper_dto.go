// Package mapper provides conversion functions between domain models and DTOs.
package mapper

import (
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/common"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/filestorage"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/transport/http/dto"
)

// PetMapper handles conversions between domain Pet model and DTOs.
type PetMapper struct {
	storage filestorage.Repository
}

// NewPetMapper creates a new PetMapper instance.
func NewPetMapper(storage filestorage.Repository) *PetMapper {
	return &PetMapper{
		storage: storage,
	}
}

// buildPhotoURLs helper builds full URLs from paths using storage
func (m *PetMapper) buildPhotoURLs(paths []string, updatedAt *time.Time) []string {
	if m.storage == nil || updatedAt == nil {
		return paths
	}
	return m.storage.BuildPhotoURLs(paths, *updatedAt)
}

// ToResponse converts a domain Pet model to a PetDetail DTO.
func (m *PetMapper) ToResponse(petmodel model.Pet) dto.PetDetail {
	petDTO := dto.PetDetail{
		ID:                   petmodel.ID,
		Name:                 petmodel.Name,
		ChipNumber:           petmodel.ChipNumber,
		OwnerName:            petmodel.OwnerName,
		OwnerID:              petmodel.OwnerID,
		AvailableBloodAmount: petmodel.CalculateDonationAmount(),
		PhotoURLs:            m.buildPhotoURLs(petmodel.PhotoURLs, petmodel.UpdatedAt),
		WeightKg:             petmodel.WeightKg,
		BirthDate:            petmodel.BirthDate,
		PetStatus:            string(petmodel.PetStatus),
		RecoveryDays:         petmodel.RecoveryDays,
		LivingCondition:      string(petmodel.LivingCondition),
		Gender:               string(petmodel.Gender),
		Type:                 string(petmodel.Type),
		ReproductiveStatus:   string(petmodel.ReproductiveStatus),
		Bonuses:              petmodel.Bonuses,
		CreatedAt:            petmodel.CreatedAt,
		UpdatedAt:            petmodel.UpdatedAt,
		DeletedAt:            petmodel.DeletedAt,
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
	petDTO.BloodGroup = petmodel.BloodGroupName

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

// ToResponseSlice converts a slice of domain Pet models to DTOs.
func (m *PetMapper) ToResponseSlice(pets []*model.Pet) []dto.PetDetail {
	if pets == nil {
		return nil
	}
	petDTOs := make([]dto.PetDetail, len(pets))
	for i, p := range pets {
		if p != nil {
			petDTOs[i] = m.ToResponse(*p)
		}
	}
	return petDTOs
}

// FromCreate converts a CreatePetBody DTO to a domain Pet model using the constructor.
func (m *PetMapper) FromCreate(petDto dto.CreatePetBody) (*model.Pet, error) {
	var birthDate *time.Time
	if petDto.BirthDate != nil {
		birthDate = petDto.BirthDate
	} else if petDto.AgeMonths != 0 || petDto.AgeYears != 0 {
		birthDate = calculateBirthDateFromAge(&petDto.AgeYears, &petDto.AgeMonths)
	}

	var breedRefID *string
	if petDto.BreedID != "" {
		breedRefID = &petDto.BreedID
	}

	var bloodGroupName string
	bloodGroup := petDto.BloodGroup
	if bloodGroup == "" {
		bloodGroup = "unknown"
	}
	bloodGroupName = bloodGroup

	var reproductiveStatus model.ReproductiveStatus
	if petDto.ReproductiveStatus != "" {
		reproductiveStatus = model.ReproductiveStatus(petDto.ReproductiveStatus)
	}

	// Handle PetHealth
	var health *model.PetHealth = &model.PetHealth{
		HealthStatus: model.HealthStatusUnknown,
	}
	if petDto.Health != nil {
		if petDto.Health.HealthStatus != nil {
			health.HealthStatus = model.HealthStatus(*petDto.Health.HealthStatus)
		}
		if petDto.Health.LastDonation != nil {
			health.LastDonation = petDto.Health.LastDonation
		}
		if petDto.Health.Transfused != nil {
			health.Transfused = petDto.Health.Transfused
		}
		if petDto.Health.Medications != nil {
			health.Medications = petDto.Health.Medications
		}
		if petDto.Health.SurgicalInterventions != nil {
			health.SurgicalInterventions = petDto.Health.SurgicalInterventions
		}
	}

	// Handle PetTreatment
	var treatments *model.PetTreatment
	if petDto.Treatments != nil {
		treatments = &model.PetTreatment{}
		if petDto.Treatments.RabiesVaccinationDate != nil {
			treatments.RabiesVaccinationDate = petDto.Treatments.RabiesVaccinationDate
		}
		if petDto.Treatments.InfectionVaccinationDate != nil {
			treatments.InfectionVaccinationDate = petDto.Treatments.InfectionVaccinationDate
		}
		if petDto.Treatments.EctoparasiteTreatmentDate != nil {
			treatments.EctoparasiteTreatmentDate = petDto.Treatments.EctoparasiteTreatmentDate
		}
		if petDto.Treatments.DewormingDate != nil {
			treatments.DewormingDate = petDto.Treatments.DewormingDate
		}
	}

	// Handle PetAnalysis
	var analyses []*model.PetAnalysis
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
					analyses = append(analyses, analysis)
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
	}

	return model.NewPet(
		petDto.Name,
		common.PetType(petDto.Type),
		petDto.WeightKg,
		model.Gender(petDto.Gender),
		"", // ownerID will be set by command handler
		birthDate,
		petDto.ChipNumber,
		model.LivingCondition(petDto.LivingCondition),
		reproductiveStatus,
		breedRefID,
		bloodGroupName,
		health,
		treatments,
		analyses,
		petDto.Bonuses,
	)
}

// ToUpdateModel converts UpdatePetBody DTO to a domain Pet model containing only changes.
// This is used to apply updates to an existing pet through the aggregate's UpdateFrom method.
func (m *PetMapper) ToUpdateModel(petDto dto.UpdatePetBody) *model.Pet {
	petUpdate := &model.Pet{}

	if petDto.Name != nil {
		petUpdate.Name = *petDto.Name
	}
	if petDto.Type != nil {
		petUpdate.Type = common.PetType(*petDto.Type)
	}
	if petDto.WeightKg != nil {
		petUpdate.WeightKg = *petDto.WeightKg
	}
	if petDto.Gender != nil {
		petUpdate.Gender = model.Gender(*petDto.Gender)
	}
	if petDto.ChipNumber != nil {
		petUpdate.ChipNumber = *petDto.ChipNumber
	}
	if petDto.LivingCondition != nil {
		petUpdate.LivingCondition = model.LivingCondition(*petDto.LivingCondition)
	}
	if petDto.BreedID != nil {
		petUpdate.BreedRefID = petDto.BreedID
	}
	if petDto.BloodGroup != nil {
		bloodGroup := *petDto.BloodGroup
		if bloodGroup == "" {
			bloodGroup = "unknown"
		}
		petUpdate.BloodGroupName = bloodGroup
	}
	if petDto.BirthDate != nil {
		petUpdate.BirthDate = petDto.BirthDate
	}
	if petDto.ReproductiveStatus != nil {
		petUpdate.ReproductiveStatus = model.ReproductiveStatus(*petDto.ReproductiveStatus)
	}
	if petDto.Bonuses != nil {
		petUpdate.Bonuses = *petDto.Bonuses
	}
	if petDto.AgeMonths != nil || petDto.AgeYears != nil {
		petUpdate.BirthDate = calculateBirthDateFromAge(petDto.AgeYears, petDto.AgeMonths)
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
		petUpdate.Health = healthmodel
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
		petUpdate.Treatments = treatmentmodel
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
		petUpdate.Analyses = analysesDomain
	}

	return petUpdate
}

// ToSimplifiedResponse converts a domain Pet model to a simplified DTO for nested usage.
// This is used when embedding pet info in other DTOs (e.g., UserDTO).
func (m *PetMapper) ToSimplifiedResponse(pet *model.Pet) dto.PetDetail {
	if pet == nil {
		return dto.PetDetail{}
	}

	dtoPet := dto.PetDetail{
		ID:         pet.ID,
		Name:       pet.Name,
		ChipNumber: pet.ChipNumber,
		PhotoURLs:  m.buildPhotoURLs(pet.PhotoURLs, pet.UpdatedAt),
		WeightKg:   pet.WeightKg,
		BirthDate:  pet.BirthDate,
		CreatedAt:  pet.CreatedAt,
		UpdatedAt:  pet.UpdatedAt,
	}

	if pet.BreedRefID != nil {
		dtoPet.BreedID = *pet.BreedRefID
	}

	dtoPet.BloodGroup = pet.BloodGroupName

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

// ToSimplifiedResponseSlice converts a slice of domain Pet models to simplified DTOs.
func (m *PetMapper) ToSimplifiedResponseSlice(pets []*model.Pet) []dto.PetDetail {
	if pets == nil {
		return nil
	}
	dtoPets := make([]dto.PetDetail, len(pets))
	for i, pet := range pets {
		dtoPets[i] = m.ToSimplifiedResponse(pet)
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
