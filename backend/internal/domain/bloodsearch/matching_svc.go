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

func hasIntersection(a, b []string) bool {
	for _, x := range a {
		if slices.Contains(b, x) {
			return true
		}
	}
	return false
}

func (r *MatchingService) MatchDonor(bloodreq *bloodreqmodel.BloodRequestWithMatchingDonors, donorPet *petmodel.Pet, preferredLocations []string) {
	sameBlood := false
	sameType := false
	sameRegion := false
	sameOwner := false
	coversNeededAmount := false
	avilableDonorAmount := donorPet.CalculateDonationAmount()
	halfVolume := (bloodreq.BloodRequest.BloodVolumeNeeded - bloodreq.BloodRequest.BloodVolumeReserved) / 2
	bloodSearchRegions := bloodreq.BloodRequest.Regions

	// (Тип-питомца) Бизнес-логика, типы питомцев должны совпадать
	sameType = donorPet.Type == bloodreq.RecipientData.PetType

	// (Регионы) Бизнес-логика, Создавая поиск, рецепиин всегда указывает регионы. И донор может быть донором, только если указал регионы.
	sameRegion = hasIntersection(bloodSearchRegions, preferredLocations)

	// (Группа-крови) Бизнес-логика, должна быть та же группа крови или неизвестная если реципиент разрешил
	sameBlood = slices.Contains(bloodreq.BloodRequest.BloodGroupNames, donorPet.BloodGroupName) || (bloodreq.BloodRequest.IncludeUnknownBloodGroup && donorPet.BloodGroupName == "UNKNOWN")

	sameOwner = bloodreq.RecipientData.OwnerID == donorPet.OwnerID
	// (Количество-крови) Бизнес-логика, донор должен покрывать весь объем или хотя бы половину от остатка если реципиент разрешил
	if bloodreq.BloodRequest.BloodVolumeReserved+avilableDonorAmount >= bloodreq.BloodRequest.BloodVolumeNeeded {
		coversNeededAmount = true
	} else if bloodreq.BloodRequest.SmallPetsNotifyAllowed && avilableDonorAmount >= halfVolume {
		coversNeededAmount = true
	}

	if sameType && sameBlood && sameRegion && !sameOwner && coversNeededAmount {
		donorBloodGroup := donorPet.BloodGroupName
		bloodreq.MatchingDonors = append(bloodreq.MatchingDonors, bloodreqmodel.MatchingDonorReadModel{
			PetName:         donorPet.Name,
			PetID:           donorPet.ID,
			DonorBloodGroup: donorBloodGroup,
			PhotoURLs:       donorPet.PhotoURLs,
			Amount:          donorPet.CalculateDonationAmount(),
		})
	}
}
