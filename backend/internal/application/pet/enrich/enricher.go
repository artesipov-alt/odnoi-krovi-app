package enrich

import (
	"context"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch"
	bloodreqmodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/bloodsearch/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor"
	donormodel "github.com/artesipov-alt/odnoi-krovi-app/internal/domain/donor/model"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"
)

// Options задаёт параметры пересчёта.
type Options struct {
	// RecoveryPeriodMonths — период восстановления в месяцах.
	// 0 = не считать, RecalculateRecoveryDays сам это обработает.
	RecoveryPeriodMonths int
}

// FetchContext — данные, добытые батчем для пересчёта статусов питомцев.
type FetchContext struct {
	// LatestApplications — последний по created_at отклик на pet, не обязательно активный —
	// решение об актуальности принимает Pet.RecalculateStatus через IsActiveForDonation().
	LatestApplications map[string]*donormodel.DonorResponse

	// BloodReqs — максимум один актуальный BloodRequest на pet — редукция уже сделана в репозитории.
	BloodReqs map[string]*bloodreqmodel.BloodRequestWithApplications
}

// PetEnricher пересчитывает статус/факторы/recovery питомцев по данным
// из репозиториев откликов доноров и заявок на поиск крови.
type PetEnricher interface {
	// Fetch батчево достаёт последние отклики и текущие заявки для списка pet ID.
	Fetch(ctx context.Context, petIDs []string) (*FetchContext, error)

	// Recalculate пересчитывает статус/факторы/recovery одного питомца по уже
	// добытому FetchContext. Возвращает отклик, использованный для пересчёта
	// (последний по дате, не обязательно активный — проверяйте IsActiveForDonation()
	// на стороне вызывающего кода, если нужен именно активный).
	Recalculate(p *model.Pet, fc *FetchContext, opts Options) *donormodel.DonorResponse

	// RecalculateAll — Fetch + Recalculate для всего среза pets одним вызовом.
	RecalculateAll(ctx context.Context, pets []*model.Pet, opts Options) (*FetchContext, error)

	// RecalculateOne — для случая, когда отклик и заявка уже получены не батчем
	// (по конкретному ID, а не "последний по дате").
	RecalculateOne(p *model.Pet, app *donormodel.DonorResponse, bloodReq *bloodreqmodel.BloodRequestWithApplications, opts Options)

	// CountFullyCompletedDonations — количество подтверждённых донаций владельца,
	// включая мягко удалённых питомцев. Тонкий проброс на репозиторий.
	CountFullyCompletedDonations(ctx context.Context, ownerID string) (int, error)
}

// Enricher — реализация PetEnricher поверх Ent-репозиториев.
type Enricher struct {
	donorRespRepo donor.Repository
	bloodReqRepo  bloodsearch.Repository
}

// New создаёт новый Enricher.
func New(donorRespRepo donor.Repository, bloodReqRepo bloodsearch.Repository) *Enricher {
	return &Enricher{donorRespRepo: donorRespRepo, bloodReqRepo: bloodReqRepo}
}

// compile-time проверка, что *Enricher реализует PetEnricher.
var _ PetEnricher = (*Enricher)(nil)

func (e *Enricher) Fetch(ctx context.Context, petIDs []string) (*FetchContext, error) {
	apps, err := e.donorRespRepo.GetLatestByPetIDs(ctx, petIDs)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get latest donor applications")
	}
	reqs, err := e.bloodReqRepo.GetByPetIDs(ctx, petIDs)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get blood requests")
	}
	return &FetchContext{LatestApplications: apps, BloodReqs: reqs}, nil
}

func (e *Enricher) Recalculate(p *model.Pet, fc *FetchContext, opts Options) *donormodel.DonorResponse {
	app := fc.LatestApplications[p.ID]
	bloodReq := fc.BloodReqs[p.ID]
	now := time.Now()
	p.RecalculateStatus(now, pet.BuildDonationContext(app, bloodReq))
	p.RecalculateRecoveryDays(opts.RecoveryPeriodMonths, now)
	return app
}

func (e *Enricher) RecalculateAll(ctx context.Context, pets []*model.Pet, opts Options) (*FetchContext, error) {
	fc, err := e.Fetch(ctx, model.CollectIDs(pets))
	if err != nil {
		return nil, err
	}
	for _, p := range pets {
		e.Recalculate(p, fc, opts)
	}
	return fc, nil
}

func (e *Enricher) RecalculateOne(p *model.Pet, app *donormodel.DonorResponse, bloodReq *bloodreqmodel.BloodRequestWithApplications, opts Options) {
	now := time.Now()
	p.RecalculateStatus(now, pet.BuildDonationContext(app, bloodReq))
	p.RecalculateRecoveryDays(opts.RecoveryPeriodMonths, now)
}

func (e *Enricher) CountFullyCompletedDonations(ctx context.Context, ownerID string) (int, error) {
	return e.donorRespRepo.CountFullyCompletedByOwnerID(ctx, ownerID)
}
