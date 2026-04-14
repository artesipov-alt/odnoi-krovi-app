package bloodsearch

import (
	"slices"

	bloodreqmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	petmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
)

type MatchingService struct{}

func NewMatchingService() *MatchingService {
	return &MatchingService{}
}

func (r *MatchingService) MatchDonor(bloodreq *bloodreqmodel.BloodRequestWithMatchingDonors, pet *petmodel.Pet) {
	sameBlood := false
	sameType := false
	coversNeededAmount := false
	avilableDonorAmount := pet.CalculateDonationAmount()
	halfVolume := (bloodreq.BloodVolumeNeeded - bloodreq.BloodVolumeReserved) / 2
	// bloodSearchRegions := bloodreq.Regions

	// (Тип-питомца) Бизнес-логика, типы питомцев должны совпадать
	sameType = pet.Type == bloodreq.RecipientData.PetType

	// (Группа-крови) Бизнес-логика, должна быть та же группа крови или неизвестная если реципиент разрешил
	sameBlood = slices.Contains(bloodreq.BloodGroupNames, pet.BloodGroupName) || (bloodreq.IncludeUnknownBloodGroup && pet.BloodGroupName == "UNKNOWN")

	// (Количество-крови) Бизнес-логика, донор должен покрывать весь объем или хотя бы половину от остатка если реципиент разрешил
	if bloodreq.BloodVolumeReserved+avilableDonorAmount >= bloodreq.BloodVolumeNeeded {
		coversNeededAmount = true
	} else if bloodreq.SmallPetsNotifyAllowed && avilableDonorAmount >= halfVolume {
		coversNeededAmount = true
	}

	if sameType && sameBlood && coversNeededAmount {
		donorBloodGroup := pet.BloodGroupName
		bloodreq.MatchingDonors = append(bloodreq.MatchingDonors, bloodreqmodel.MatchingDonorReadModel{
			PetName:         pet.Name,
			PetID:           pet.ID,
			DonorBloodGroup: donorBloodGroup,
			PhotoURLs:       pet.PhotoURLs,
			Amount:          pet.CalculateDonationAmount(),
		})
	}
}
