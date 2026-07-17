package query

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/application/pet/enrich"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	petmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
)

// GetByPetIDResult содержит результат поиска заявки с подходящими и потенциальными донорами.
type GetByPetIDResult struct {
	BloodRequest    *model.BloodRequestWithApplications
	SuitableDonors  int
	PotentialDonors []*petmodel.PotentialDonor
}

type GetByPetIDHandler struct {
	bloodRepo    bloodsearch.Repository
	petRepo      pet.Repository
	petEnricher  enrich.PetEnricher
	bloodCounter *bloodsearch.BloodCounterService
}

func NewGetByPetIDHandler(
	bloodRepo bloodsearch.Repository,
	petRepo pet.Repository,
	petEnricher enrich.PetEnricher,
) *GetByPetIDHandler {
	return &GetByPetIDHandler{
		bloodRepo:    bloodRepo,
		petRepo:      petRepo,
		petEnricher:  petEnricher,
		bloodCounter: bloodsearch.NewBloodCounterService(),
	}
}

func (h *GetByPetIDHandler) Handle(ctx context.Context, petID string) (*GetByPetIDResult, error) {
	bloodReq, err := h.bloodRepo.GetByPetID(ctx, petID)
	if err != nil {
		return nil, err
	}

	donated, reserved := h.bloodCounter.RecalculateBloodAmount(bloodReq.BloodRequest, bloodReq.DonorApplications)
	bloodReq.BloodRequest.SetBloodVolume(donated, reserved)

	suitableDonors, err := h.petRepo.CountSuitableDonors(ctx, bloodReq.BloodRequest.BloodGroupNames)
	if err != nil {
		return nil, err
	}

	// Загружаем питомца-реципиента, чтобы получить его тип для поиска потенциальных доноров
	recipientPet, err := h.petRepo.GetByID(ctx, petID, pet.PetPreloadOptions{})
	if err != nil {
		return nil, err
	}

	potentialDonors, err := h.petRepo.FindPotentialDonors(ctx, pet.PotentialDonorsCriteria{
		PetType:          recipientPet.Type,
		BloodGroups:      bloodReq.SearchingBloodGroupNames(),
		Regions:          bloodReq.BloodRequest.Regions,
		ExcludeRequestID: bloodReq.BloodRequest.ID,
		ExcludePetID:     petID,
	})
	if err != nil {
		return nil, err
	}

	// Извлекаем питомцев для батчевого Fetch и пересчёта статуса/факторов
	pets := make([]*petmodel.Pet, len(potentialDonors))
	for i, pd := range potentialDonors {
		pets[i] = pd.Pet
	}

	// Батчевый Fetch для всех потенциальных доноров
	fc, err := h.petEnricher.Fetch(ctx, petIDsOf(pets))
	if err != nil {
		return nil, err
	}

	// Пересчёт статуса, факторов и recovery для каждого донора с его индивидуальным периодом
	donors := make([]*petmodel.PotentialDonor, 0, len(potentialDonors))
	for _, pd := range potentialDonors {
		h.petEnricher.Recalculate(pd.Pet, fc, enrich.Options{
			RecoveryPeriodMonths: pd.RecoveryPeriodMonths,
		})

		if pd.Pet.PetStatus == petmodel.PetStatusDonor {
			donors = append(donors, pd)
		}
	}

	return &GetByPetIDResult{
		BloodRequest:    bloodReq,
		SuitableDonors:  suitableDonors,
		PotentialDonors: donors,
	}, nil
}

func petIDsOf(pets []*petmodel.Pet) []string {
	ids := make([]string, len(pets))
	for i, p := range pets {
		ids[i] = p.ID
	}
	return ids
}
