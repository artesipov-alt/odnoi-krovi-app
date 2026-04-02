package model

// import (
// 	"slices"

// 	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/common"
// 	petmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
// )

// // Recipient представляет модель чтения реципиент
// type Recipient struct {
// 	ID                       string
// 	PetID                    string
// 	PetName                  string
// 	PetType                  petmodel.PetType
// 	OwnerName                string
// 	PhotoURLs                []string
// 	SearchRegions            []string
// 	SearchingBloodNames      []string
// 	BloodVolumeNeeded        int32
// 	BloodVolumeReserved      int32
// 	BloodGroupName           string
// 	PrioritySearch           bool
// 	IncludeUnknownBloodGroup bool
// 	SmallPetsNotifyAllowed   bool
// 	Status                   string
// 	MatchingDonors           []MatchingDonorReadModel
// 	DefaultDonorPrefs        *DefaultDonorPrefs
// 	AdvancedInfo             *AdvancedInfo
// }

// type AdvancedInfo struct {
// 	Description string
// 	PhotoURLs   []string
// }

// type DefaultDonorPrefs struct {
// 	CompensationType common.CompensationType
// 	Bonuses          []string
// 	TaxiCompensation bool
// }

// // MatchingDonorReadModel представляет модель чтения для подходящего донора
// type MatchingDonorReadModel struct {
// 	PetID           string
// 	PetName         string
// 	Amount          int32
// 	DonorBloodGroup string
// 	PhotoURLs       []string
// }

// func (r *Recipient) MatchDonor(pet *petmodel.Pet) {

// 	sameBlood := false
// 	coversNeededAmount := false
// 	avilableDonorAmount := pet.CalculateDonationAmount()
// 	halfVolume := (r.BloodVolumeNeeded - r.BloodVolumeReserved) / 2

// 	// (Группа-крови) Бизнес-логика, должна быть та же группа крови или любая если реципиент разрешил
// 	if pet.BloodGroupName != nil {
// 		sameBlood = slices.Contains(r.SearchingBloodNames, *pet.BloodGroupName)
// 	} else if r.IncludeUnknownBloodGroup {
// 		sameBlood = true
// 	}

// 	// (Количество-крови) Бизнес-логика, донор должен покрывать весь объем или хотя бы половину от остатка если реципиент разрешил
// 	if r.BloodVolumeReserved+avilableDonorAmount >= r.BloodVolumeNeeded {
// 		coversNeededAmount = true
// 	} else if r.SmallPetsNotifyAllowed && avilableDonorAmount >= halfVolume {
// 		coversNeededAmount = true
// 	}

// 	if sameBlood && coversNeededAmount {
// 		donorBloodGroup := ""
// 		if pet.BloodGroupName != nil {
// 			donorBloodGroup = *pet.BloodGroupName
// 		}
// 		r.MatchingDonors = append(r.MatchingDonors, MatchingDonorReadModel{
// 			PetName:         pet.Name,
// 			PetID:           pet.ID,
// 			DonorBloodGroup: donorBloodGroup,
// 			PhotoURLs:       pet.PhotoURLs,
// 			Amount:          pet.CalculateDonationAmount(),
// 		})
// 	}
// }

// func (r *Recipient) SetDefaultPrefs(compensationType common.CompensationType, taxiCompensation bool) {
// 	r.DefaultDonorPrefs = &DefaultDonorPrefs{
// 		CompensationType: compensationType,
// 		TaxiCompensation: taxiCompensation,
// 		Bonuses:          []string{},
// 	}
// }
