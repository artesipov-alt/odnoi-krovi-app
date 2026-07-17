package query

import (
	"context"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/application/pet/enrich"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	bloodsearchmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	petmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
)

// GetByPetIDResult содержит результат поиска заявки с подходящими и потенциальными донорами.
type GetByPetIDResult struct {
	BloodRequest    *bloodsearchmodel.BloodRequestWithApplications
	SuitableDonors  int
	PotentialDonors []*bloodsearchmodel.PotentialDonor
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

func (h *GetByPetIDHandler) Handle(ctx context.Context, callerUserID, petID string) (*GetByPetIDResult, error) {
	// Строгая авторизация: без userID потенциальные доноры не возвращаются.
	// HTTP-layer всегда передаёт userID (иначе middleware Auth блокирует раньше),
	// а прямые вызовы должны идти через HTTP — здесь просто защита от случайного вызова.
	if callerUserID == "" {
		return nil, apperrors.Unauthorized("missing user context")
	}
	bloodReq, err := h.bloodRepo.GetByPetID(ctx, petID)
	if err != nil {
		return nil, err
	}

	donated, reserved := h.bloodCounter.RecalculateBloodAmount(bloodReq.BloodRequest, bloodReq.DonorApplications)
	bloodReq.BloodRequest.SetBloodVolume(donated, reserved)
	bloodReq.BloodRequest.RecalculateStatus()

	suitableDonors, err := h.petRepo.CountSuitableDonors(ctx, bloodReq.BloodRequest.BloodGroupNames)
	if err != nil {
		return nil, err
	}

	// Если заявка полностью зарезервирована, нет смысла искать новых
	// потенциальных доноров — нужный объём крови уже обеспечен откликнувшимися.
	if bloodReq.BloodRequest.IsReservedFull() {
		return &GetByPetIDResult{
			BloodRequest:   bloodReq,
			SuitableDonors: suitableDonors,
		}, nil
	}

	// Загружаем питомца-реципиента, чтобы получить его тип для поиска потенциальных доноров
	// и убедиться, что caller — владелец питомца (потенциальные доноры видны только ему).
	recipientPet, err := h.petRepo.GetByID(ctx, petID, pet.PetPreloadOptions{})
	if err != nil {
		return nil, apperrors.NotFound("recipient pet not found")
	}
	if recipientPet.OwnerID != callerUserID {
		return nil, apperrors.Forbidden("only the recipient owner can view potential donors")
	}

	potentialDonors, err := h.petRepo.FindPotentialDonors(ctx, pet.PotentialDonorsCriteria{
		PetType:          recipientPet.Type,
		BloodGroups:      bloodReq.SearchingBloodGroupNames(),
		Regions:          bloodReq.BloodRequest.Regions,
		ExcludeRequestID: bloodReq.BloodRequest.ID,
		ExcludeOwnerID:   callerUserID,
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
	donors := make([]*bloodsearchmodel.PotentialDonor, 0, len(potentialDonors))
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
