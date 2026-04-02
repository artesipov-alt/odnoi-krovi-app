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
	coversNeededAmount := false
	avilableDonorAmount := pet.CalculateDonationAmount()
	halfVolume := (bloodreq.BloodVolumeNeeded - bloodreq.BloodVolumeReserved) / 2

	// (Группа-крови) Бизнес-логика, должна быть та же группа крови или любая если реципиент разрешил
	if pet.BloodGroupName != nil {
		sameBlood = slices.Contains(bloodreq.BloodGroupNames, *pet.BloodGroupName)
	} else if bloodreq.IncludeUnknownBloodGroup {
		sameBlood = true
	}

	// (Количество-крови) Бизнес-логика, донор должен покрывать весь объем или хотя бы половину от остатка если реципиент разрешил
	if bloodreq.BloodVolumeReserved+avilableDonorAmount >= bloodreq.BloodVolumeNeeded {
		coversNeededAmount = true
	} else if bloodreq.SmallPetsNotifyAllowed && avilableDonorAmount >= halfVolume {
		coversNeededAmount = true
	}

	if sameBlood && coversNeededAmount {
		donorBloodGroup := ""
		if pet.BloodGroupName != nil {
			donorBloodGroup = *pet.BloodGroupName
		}
		bloodreq.MatchingDonors = append(bloodreq.MatchingDonors, bloodreqmodel.MatchingDonorReadModel{
			PetName:         pet.Name,
			PetID:           pet.ID,
			DonorBloodGroup: donorBloodGroup,
			PhotoURLs:       pet.PhotoURLs,
			Amount:          pet.CalculateDonationAmount(),
		})
	}
}
