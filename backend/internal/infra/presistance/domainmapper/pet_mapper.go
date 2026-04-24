package domainmapper

import (
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/common"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
)

// PetToDomain converts ent.Pet to domain model.Pet
func PetToDomain(e *ent.Pet) *model.Pet {
	if e == nil {
		return nil
	}

	pet := &model.Pet{
		ID:                 e.ID,
		Name:               e.Name,
		Type:               common.PetType(e.Type),
		PetStatus:          model.PetStatusNone,
		WeightKg:           e.WeightKg,
		Gender:             model.Gender(e.Gender),
		BirthDate:          e.BirthDate,
		ChipNumber:         e.ChipNumber,
		PhotoURLs:          e.PhotoUrls,
		LivingCondition:    model.LivingCondition(e.LivingCondition),
		ReproductiveStatus: model.ReproductiveStatus(e.ReproductiveStatus),
		OwnerID:            e.UserID,
		OwnerName:          e.Edges.Owner.FullName,
		BreedRefID:         e.BreedID,
		CreatedAt:          &e.CreatedAt,
		UpdatedAt:          &e.UpdatedAt,
		DeletedAt:          e.DeletedAt,
	}

	if e.Privilege != nil {
		pet.Privilege = common.Privilege(*e.Privilege)
	}

	// Map BloodGroupName directly from field
	pet.BloodGroupName = e.BloodGroup

	// Map BreedRefID from edge if available
	if e.Edges.BreedRef != nil {
		pet.BreedRefID = &e.Edges.BreedRef.ID
	}

	// Map Health
	if e.Edges.Health != nil {
		pet.Health = &model.PetHealth{
			HealthStatus:          model.HealthStatus(e.Edges.Health.HealthStatus),
			LastDonation:          e.Edges.Health.LastDonation,
			Transfused:            &e.Edges.Health.Transfused,
			Medications:           &e.Edges.Health.Medications,
			SurgicalInterventions: &e.Edges.Health.SurgicalInterventions,
		}
	}

	// Map Treatments
	if e.Edges.Treatments != nil {
		pet.Treatments = &model.PetTreatment{
			RabiesVaccinationDate:     e.Edges.Treatments.RabiesVaccinationDate,
			InfectionVaccinationDate:  e.Edges.Treatments.InfectionVaccinationDate,
			EctoparasiteTreatmentDate: e.Edges.Treatments.EctoparasiteTreatmentDate,
			DewormingDate:             e.Edges.Treatments.DewormingDate,
		}
	}

	// Map Analyses
	if len(e.Edges.Analyses) > 0 {
		pet.Analyses = make([]*model.PetAnalysis, len(e.Edges.Analyses))
		for i, a := range e.Edges.Analyses {
			pet.Analyses[i] = &model.PetAnalysis{
				ID:           a.ID,
				AnalysisName: string(a.AnalysisName),
				AnalysisType: string(a.AnalysisType),
				AnalysisDate: a.AnalysisDate,
			}
		}
	}

	return pet
}

// PetToDomainSlice converts slice of ent.Pet to slice of domain model.Pet
func PetToDomainSlice(pets []*ent.Pet) []*model.Pet {
	if pets == nil {
		return nil
	}
	result := make([]*model.Pet, len(pets))
	for i, p := range pets {
		result[i] = PetToDomain(p)
	}
	return result
}
