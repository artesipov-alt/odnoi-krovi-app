package mocks

// import (
// 	"time"

// 	"github.com/artesipov-alt/odnoi-krovi-app/internal/dto"
// )

// var Pet1 = dto.Pet{
// 	ID:              "PET-001",
// 	Name:            "Рекс",
// 	ChipNumber:      "123456789012345",
// 	PhotoURLs:       []string{"https://example.com/rex.jpg"},
// 	BreedID:         "LABRADOR",
// 	WeightKg:        30.5,
// 	BirthDate:       func() *time.Time { t, _ := time.Parse(time.RFC3339, "2019-03-10T00:00:00Z"); return &t }(),
// 	LivingCondition: "indoor",
// 	Gender:          "male",
// 	Type:            "dog",
// 	BloodGroup:      "DEA 1+",
// 	PetStatus:       "donor",
// 	Health: &dto.PetHealth{
// 		HealthStatus: func() *string { s := "healthy"; return &s }(),
// 		LastDonation: func() *time.Time { t, _ := time.Parse(time.RFC3339, "2023-09-15T10:00:00Z"); return &t }(),
// 	},
// 	Treatments: &dto.PetTreatment{
// 		RabiesVaccinationDate: func() *time.Time { t, _ := time.Parse(time.RFC3339, "2023-08-01T12:00:00Z"); return &t }(),
// 	},
// 	Analyses: &dto.PetAnalysisGroup{
// 		Leukemia: []*dto.PetAnalysis{
// 			{
// 				ID:           func() *string { s := "PAN-LEU-001"; return &s }(),
// 				AnalysisName: func() *string { s := "leukemia"; return &s }(),
// 				AnalysisType: func() *string { s := "PCR"; return &s }(),
// 				AnalysisDate: func() *time.Time { t, _ := time.Parse(time.RFC3339, "2023-07-20T12:00:00Z"); return &t }(),
// 			},
// 		},
// 	},
// 	Bonuses: &dto.PetBonus{
// 		IsFormerDonor: true,
// 	},
// 	CreatedAt: func() *time.Time { t, _ := time.Parse(time.RFC3339, "2023-01-01T12:00:00Z"); return &t }(),
// 	UpdatedAt: func() *time.Time { t, _ := time.Parse(time.RFC3339, "2023-10-01T12:00:00Z"); return &t }(),
// }

// var Pet2 = dto.Pet{
// 	ID:              "PET-002",
// 	Name:            "Мурка",
// 	ChipNumber:      "987654321098765",
// 	PhotoURLs:       []string{"https://example.com/murka.jpg"},
// 	BreedID:         "SIAMESE",
// 	WeightKg:        4.2,
// 	BirthDate:       func() *time.Time { t, _ := time.Parse(time.RFC3339, "2021-01-20T00:00:00Z"); return &t }(),
// 	LivingCondition: "indoor",
// 	Gender:          "female",
// 	Type:            "cat",
// 	BloodGroup:      "A",
// 	PetStatus:       "donor",
// 	Health: &dto.PetHealth{
// 		HealthStatus: func() *string { s := "healthy"; return &s }(),
// 		LastDonation: func() *time.Time { t, _ := time.Parse(time.RFC3339, "2023-10-05T11:00:00Z"); return &t }(),
// 	},
// 	Treatments: &dto.PetTreatment{
// 		InfectionVaccinationDate: func() *time.Time { t, _ := time.Parse(time.RFC3339, "2023-09-01T12:00:00Z"); return &t }(),
// 	},
// 	Analyses: &dto.PetAnalysisGroup{
// 		Immunodeficiency: []*dto.PetAnalysis{
// 			{
// 				ID:           func() *string { s := "PAN-IMM-001"; return &s }(),
// 				AnalysisName: func() *string { s := "immunodeficiency"; return &s }(),
// 				AnalysisType: func() *string { s := "ELISA"; return &s }(),
// 				AnalysisDate: func() *time.Time { t, _ := time.Parse(time.RFC3339, "2023-08-25T12:00:00Z"); return &t }(),
// 			},
// 		},
// 	},
// 	Bonuses: &dto.PetBonus{
// 		IsArtist: true,
// 	},
// 	CreatedAt: func() *time.Time { t, _ := time.Parse(time.RFC3339, "2023-02-01T12:00:00Z"); return &t }(),
// 	UpdatedAt: func() *time.Time { t, _ := time.Parse(time.RFC3339, "2023-10-05T12:00:00Z"); return &t }(),
// }

// var Pet3 = dto.Pet{
// 	ID:              "PET-003",
// 	Name:            "Барон",
// 	ChipNumber:      "112233445566778",
// 	PhotoURLs:       []string{"https://example.com/baron.jpg"},
// 	BreedID:         "GERMAN_SHEPHERD",
// 	WeightKg:        40.0,
// 	BirthDate:       func() *time.Time { t, _ := time.Parse(time.RFC3339, "2018-07-01T00:00:00Z"); return &t }(),
// 	LivingCondition: "leash_walking",
// 	Gender:          "male",
// 	Type:            "dog",
// 	BloodGroup:      "DEA 1-",
// 	PetStatus:       "donor",
// 	Health: &dto.PetHealth{
// 		HealthStatus: func() *string { s := "healthy"; return &s }(),
// 		LastDonation: func() *time.Time { t, _ := time.Parse(time.RFC3339, "2023-09-20T09:00:00Z"); return &t }(),
// 	},
// 	Treatments: &dto.PetTreatment{
// 		DewormingDate: func() *time.Time { t, _ := time.Parse(time.RFC3339, "2023-08-10T12:00:00Z"); return &t }(),
// 	},
// 	Analyses: &dto.PetAnalysisGroup{
// 		Dirofilaria: []*dto.PetAnalysis{
// 			{
// 				ID:           func() *string { s := "PAN-DIR-001"; return &s }(),
// 				AnalysisName: func() *string { s := "dirofilaria"; return &s }(),
// 				AnalysisType: func() *string { s := "Microscopy"; return &s }(),
// 				AnalysisDate: func() *time.Time { t, _ := time.Parse(time.RFC3339, "2023-07-15T12:00:00Z"); return &t }(),
// 			},
// 		},
// 	},
// 	Bonuses: &dto.PetBonus{
// 		IsGuideDog: true,
// 	},
// 	CreatedAt: func() *time.Time { t, _ := time.Parse(time.RFC3339, "2023-03-01T12:00:00Z"); return &t }(),
// 	UpdatedAt: func() *time.Time { t, _ := time.Parse(time.RFC3339, "2023-10-02T12:00:00Z"); return &t }(),
// }
