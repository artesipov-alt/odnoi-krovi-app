package pg

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
	entbloodreq "github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/bloodsearchrequest"
	entdonorpreference "github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/donorpreference"
	entpet "github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/petanalysis"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/pethealth"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/pettreatment"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/predicate"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/schema"
	entuser "github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/user"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/presistance/domainmapper"
)

// EntPetRepository реализует PetRepository с использованием ENT
type EntPetRepository struct {
	client *ent.Client
}

// NewEntPetRepository создает новый репозиторий питомцев ENT
func NewEntPetRepository(client *ent.Client) *EntPetRepository {
	return &EntPetRepository{
		client: client,
	}
}

// Create создает нового питомца в базе данных вместе с его связанными сущностями в транзакции
func (r *EntPetRepository) Create(ctx context.Context, petDomain *model.Pet) (*model.Pet, error) {
	if petDomain == nil {
		return nil, errors.New("входные данные питомца не могут быть nil")
	}

	// Начать транзакцию
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("не удалось начать транзакцию: %w", err)
	}

	// 1. Создать питомца
	builder := tx.Pet.Create().
		SetName(petDomain.Name).
		SetType(string(petDomain.Type)).
		SetWeightKg(petDomain.WeightKg).
		SetUserID(petDomain.OwnerID).
		SetBonuses(petDomain.Bonuses)

	if petDomain.Gender != "" {
		builder.SetGender(string(petDomain.Gender))
	}

	if petDomain.BirthDate != nil {
		builder.SetBirthDate(*petDomain.BirthDate)
	}
	if petDomain.ChipNumber != "" {
		builder.SetChipNumber(petDomain.ChipNumber)
	}
	if len(petDomain.PhotoURLs) > 0 {
		builder.SetPhotoUrls(petDomain.PhotoURLs)
	}
	if petDomain.LivingCondition != "" {
		builder.SetLivingCondition(string(petDomain.LivingCondition))
	}
	if petDomain.ReproductiveStatus != "" {
		builder.SetReproductiveStatus(string(petDomain.ReproductiveStatus))
	}
	if petDomain.BreedRefID != nil {
		builder.SetBreedRefID(*petDomain.BreedRefID)
	}

	if petDomain.BloodGroupName != "" {
		builder.SetBloodGroup(petDomain.BloodGroupName)
	}

	newPet, err := builder.Save(ctx)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("не удалось создать питомца: %w", err)
	}

	// 2. Создать связанные сущности, если они предоставлены
	if petDomain.Health != nil {
		healthBuilder := tx.PetHealth.Create().
			SetHealthStatus(pethealth.HealthStatus(petDomain.Health.HealthStatus)).
			SetOwner(newPet)

		if petDomain.Health.LastDonation != nil {
			healthBuilder.SetLastDonation(*petDomain.Health.LastDonation)
		}
		if petDomain.Health.Transfused != nil {
			healthBuilder.SetTransfused(*petDomain.Health.Transfused)
		}
		if petDomain.Health.Medications != nil {
			healthBuilder.SetMedications(*petDomain.Health.Medications)
		}
		if petDomain.Health.SurgicalInterventions != nil {
			healthBuilder.SetSurgicalInterventions(*petDomain.Health.SurgicalInterventions)
		}

		_, err = healthBuilder.Save(ctx)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("не удалось создать данные о здоровье питомца: %w", err)
		}
	}

	if petDomain.Treatments != nil {
		treatmentBuilder := tx.PetTreatment.Create().
			SetOwner(newPet)

		if petDomain.Treatments.RabiesVaccinationDate != nil {
			treatmentBuilder.SetRabiesVaccinationDate(*petDomain.Treatments.RabiesVaccinationDate)
		}
		if petDomain.Treatments.InfectionVaccinationDate != nil {
			treatmentBuilder.SetInfectionVaccinationDate(*petDomain.Treatments.InfectionVaccinationDate)
		}
		if petDomain.Treatments.EctoparasiteTreatmentDate != nil {
			treatmentBuilder.SetEctoparasiteTreatmentDate(*petDomain.Treatments.EctoparasiteTreatmentDate)
		}
		if petDomain.Treatments.DewormingDate != nil {
			treatmentBuilder.SetDewormingDate(*petDomain.Treatments.DewormingDate)
		}

		_, err = treatmentBuilder.Save(ctx)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("не удалось создать данные о лечении питомца: %w", err)
		}
	}

	for _, a := range petDomain.Analyses {
		_, err = tx.PetAnalysis.Create().
			SetAnalysisName(petanalysis.AnalysisName(a.AnalysisName)).
			SetAnalysisType(petanalysis.AnalysisType(a.AnalysisType)).
			SetAnalysisDate(*a.AnalysisDate).
			SetPetID(newPet.ID).
			Save(ctx)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("не удалось создать анализ питомца: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("не удалось зафиксировать транзакцию: %w", err)
	}

	// Re-fetch the created pet with all relations to return complete aggregate
	return r.GetByID(ctx, newPet.ID, pet.PetPreloadOptions{
		WithHealth:     petDomain.Health != nil,
		WithTreatments: petDomain.Treatments != nil,
		WithAnalyses:   len(petDomain.Analyses) > 0,
	})
}

// GetByID получает питомца по ID с опциями загрузки связанных данных
func (r *EntPetRepository) GetByID(ctx context.Context, id string, opts pet.PetPreloadOptions) (*model.Pet, error) {
	pquery := r.client.Pet.Query().Where(entpet.ID(id)).
		WithBreedRef().
		WithOwner(func(uq *ent.UserQuery) {
			uq.Select(entuser.FieldFullName)
		})

	// Apply preload options
	if opts.WithAll {
		pquery = pquery.WithHealth().WithTreatments().WithAnalyses()
	} else {
		if opts.WithHealth {
			pquery = pquery.WithHealth()
		}
		if opts.WithTreatments {
			pquery = pquery.WithTreatments()
		}
		if opts.WithAnalyses {
			pquery = pquery.WithAnalyses()
		}
	}

	entPet, err := pquery.Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperrors.ErrPetNotFound
		}
		return nil, apperrors.Internal(err, "failed to get pet")
	}

	return domainmapper.PetToDomain(entPet), nil
}

// GetByUserID получает всех питомцев пользователя
func (r *EntPetRepository) GetByUserID(ctx context.Context, userID string, opts pet.PetPreloadOptions) ([]*model.Pet, error) {
	pquery := r.client.Pet.Query().Where(entpet.UserID(userID)).
		WithBreedRef().
		WithOwner(func(uq *ent.UserQuery) {
			uq.Select(entuser.FieldFullName)
		})

	// Apply preload options
	if opts.WithAll {
		pquery = pquery.WithHealth().WithTreatments().WithAnalyses()
	} else {
		if opts.WithHealth {
			pquery = pquery.WithHealth()
		}
		if opts.WithTreatments {
			pquery = pquery.WithTreatments()
		}
		if opts.WithAnalyses {
			pquery = pquery.WithAnalyses()
		}
	}

	pets, err := pquery.All(ctx)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get pets")
	}

	return domainmapper.PetToDomainSlice(pets), nil
}

func (r *EntPetRepository) Update(ctx context.Context, id string, petDomain *model.Pet) (*model.Pet, error) {
	if id == "" {
		return nil, errors.New("неверный ID питомца")
	}

	// Начать транзакцию
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("не удалось начать транзакцию: %w", err)
	}

	updater := tx.Pet.UpdateOneID(id)

	if petDomain.Name != "" {
		updater.SetName(petDomain.Name)
	}
	if petDomain.Type != "" {
		updater.SetType(string(petDomain.Type))
	}
	if petDomain.WeightKg != 0 {
		updater.SetWeightKg(petDomain.WeightKg)
	}
	if petDomain.Gender != "" {
		updater.SetGender(string(petDomain.Gender))
	}
	if petDomain.BirthDate != nil {
		updater.SetBirthDate(*petDomain.BirthDate)
	}
	if petDomain.ChipNumber != "" {
		updater.SetChipNumber(petDomain.ChipNumber)
	}
	if petDomain.LivingCondition != "" {
		updater.SetLivingCondition(string(petDomain.LivingCondition))
	}
	if petDomain.ReproductiveStatus != "" {
		updater.SetReproductiveStatus(string(petDomain.ReproductiveStatus))
	}
	if petDomain.BreedRefID != nil {
		updater.SetBreedRefID(*petDomain.BreedRefID)
	}
	if petDomain.BloodGroupName != "" {
		updater.SetBloodGroup(petDomain.BloodGroupName)
	}

	_, err = updater.Save(ctx)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("не удалось обновить питомца: %w", err)
	}

	// 2. Обновить данные о здоровье
	if petDomain.Health != nil {
		existingHealth, err := tx.PetHealth.Query().Where(pethealth.HasOwnerWith(entpet.ID(id))).Only(ctx)
		if err != nil && !ent.IsNotFound(err) {
			tx.Rollback()
			return nil, err
		}
		if existingHealth != nil {
			updater := tx.PetHealth.UpdateOne(existingHealth)
			updater.SetHealthStatus(pethealth.HealthStatus(petDomain.Health.HealthStatus))
			if petDomain.Health.LastDonation != nil {
				updater.SetLastDonation(*petDomain.Health.LastDonation)
			} else {
				updater.ClearLastDonation()
			}
			if petDomain.Health.Transfused != nil {
				updater.SetTransfused(*petDomain.Health.Transfused)
			} else {
				updater.ClearTransfused()
			}
			if petDomain.Health.Medications != nil {
				updater.SetMedications(*petDomain.Health.Medications)
			} else {
				updater.ClearMedications()
			}
			if petDomain.Health.SurgicalInterventions != nil {
				updater.SetSurgicalInterventions(*petDomain.Health.SurgicalInterventions)
			} else {
				updater.ClearSurgicalInterventions()
			}
			_, err = updater.Save(ctx)
		} else {
			healthBuilder := tx.PetHealth.Create().
				SetHealthStatus(pethealth.HealthStatus(petDomain.Health.HealthStatus)).
				SetOwnerID(id)
			if petDomain.Health.LastDonation != nil {
				healthBuilder.SetLastDonation(*petDomain.Health.LastDonation)
			}
			if petDomain.Health.Transfused != nil {
				healthBuilder.SetTransfused(*petDomain.Health.Transfused)
			}
			if petDomain.Health.Medications != nil {
				healthBuilder.SetMedications(*petDomain.Health.Medications)
			}
			if petDomain.Health.SurgicalInterventions != nil {
				healthBuilder.SetSurgicalInterventions(*petDomain.Health.SurgicalInterventions)
			}
			_, err = healthBuilder.Save(ctx)
		}
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("не удалось обновить данные о здоровье питомца: %w", err)
		}
	}

	// 3. Обновить данные о лечении
	if petDomain.Treatments != nil {
		existingTreatment, err := tx.PetTreatment.Query().Where(pettreatment.HasOwnerWith(entpet.ID(id))).Only(ctx)
		if err != nil && !ent.IsNotFound(err) {
			tx.Rollback()
			return nil, err
		}
		if existingTreatment != nil {
			updater := tx.PetTreatment.UpdateOne(existingTreatment)
			if petDomain.Treatments.RabiesVaccinationDate != nil {
				updater.SetRabiesVaccinationDate(*petDomain.Treatments.RabiesVaccinationDate)
			} else {
				updater.ClearRabiesVaccinationDate()
			}
			if petDomain.Treatments.InfectionVaccinationDate != nil {
				updater.SetInfectionVaccinationDate(*petDomain.Treatments.InfectionVaccinationDate)
			} else {
				updater.ClearInfectionVaccinationDate()
			}
			if petDomain.Treatments.EctoparasiteTreatmentDate != nil {
				updater.SetEctoparasiteTreatmentDate(*petDomain.Treatments.EctoparasiteTreatmentDate)
			} else {
				updater.ClearEctoparasiteTreatmentDate()
			}
			if petDomain.Treatments.DewormingDate != nil {
				updater.SetDewormingDate(*petDomain.Treatments.DewormingDate)
			} else {
				updater.ClearDewormingDate()
			}
			_, err = updater.Save(ctx)
		} else {
			treatmentBuilder := tx.PetTreatment.Create().
				SetOwnerID(id)
			if petDomain.Treatments.RabiesVaccinationDate != nil {
				treatmentBuilder.SetRabiesVaccinationDate(*petDomain.Treatments.RabiesVaccinationDate)
			}
			if petDomain.Treatments.InfectionVaccinationDate != nil {
				treatmentBuilder.SetInfectionVaccinationDate(*petDomain.Treatments.InfectionVaccinationDate)
			}
			if petDomain.Treatments.EctoparasiteTreatmentDate != nil {
				treatmentBuilder.SetEctoparasiteTreatmentDate(*petDomain.Treatments.EctoparasiteTreatmentDate)
			}
			if petDomain.Treatments.DewormingDate != nil {
				treatmentBuilder.SetDewormingDate(*petDomain.Treatments.DewormingDate)
			}
			_, err = treatmentBuilder.Save(ctx)
		}
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("не удалось обновить данные о лечении питомца: %w", err)
		}
	}

	// Для анализов, удалить существующие и создать новые, если предоставлены
	if petDomain.Analyses != nil {
		// Удалить существующие анализы для этого питомца (жесткое удаление)
		_, err = tx.PetAnalysis.Delete().Where(petanalysis.PetID(id)).Exec(schema.SkipSoftDelete(ctx))
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("не удалось удалить существующие анализы питомца: %w", err)
		}

		// Создать новые анализы
		for _, a := range petDomain.Analyses {
			_, err = tx.PetAnalysis.Create().
				SetAnalysisName(petanalysis.AnalysisName(a.AnalysisName)).
				SetAnalysisType(petanalysis.AnalysisType(a.AnalysisType)).
				SetAnalysisDate(*a.AnalysisDate).
				SetPetID(id).
				Save(ctx)
			if err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("не удалось создать анализ питомца: %w", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("не удалось зафиксировать транзакцию: %w", err)
	}

	// Re-fetch the updated pet with all relations to return complete aggregate
	return r.GetByID(ctx, id, pet.PetPreloadOptions{
		WithHealth:     petDomain.Health != nil,
		WithTreatments: petDomain.Treatments != nil,
		WithAnalyses:   len(petDomain.Analyses) > 0,
	})
}

// Delete удаляет питомца по его ID (мягкое удаление)
func (r *EntPetRepository) DeleteWithRelations(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("неверный ID питомца")
	}

	// Get the pet to access related IDs
	petEnt, err := r.client.Pet.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return fmt.Errorf("питомец с ID %s не найден", id)
		}
		return fmt.Errorf("не удалось получить питомца: %w", err)
	}

	// Delete health if exists
	if petEnt.HealthID != "" {
		err = r.client.PetHealth.DeleteOneID(petEnt.HealthID).Exec(ctx)
		if err != nil && !ent.IsNotFound(err) {
			return fmt.Errorf("не удалось удалить здоровье питомца: %w", err)
		}
	}

	// Delete treatment if exists
	if petEnt.TreatmentID != "" {
		err = r.client.PetTreatment.DeleteOneID(petEnt.TreatmentID).Exec(ctx)
		if err != nil && !ent.IsNotFound(err) {
			return fmt.Errorf("не удалось удалить лечение питомца: %w", err)
		}
	}

	// Delete analyses
	_, err = r.client.PetAnalysis.Delete().Where(petanalysis.PetID(id)).Exec(ctx)
	if err != nil {
		return fmt.Errorf("не удалось удалить анализы питомца: %w", err)
	}

	// Delete donor responses where this pet is the donor
	_, err = r.client.DonorResponse.Delete().Where(sql.FieldEQ("donor_response_donor", id)).Exec(ctx)
	if err != nil {
		return fmt.Errorf("не удалось удалить отклики донора: %w", err)
	}

	// Delete blood request and responses if exists
	bloodReq, err := r.client.BloodSearchRequest.Query().Where(entbloodreq.PetID(id)).Only(ctx)
	if err != nil && !ent.IsNotFound(err) {
		return fmt.Errorf("не удалось получить заявку на кровь: %w", err)
	}
	if bloodReq != nil {
		// Delete responses
		_, err = r.client.DonorResponse.Delete().Where(sql.FieldEQ("blood_search_request_responses", bloodReq.ID)).Exec(ctx)
		if err != nil {
			return fmt.Errorf("не удалось удалить отклики доноров: %w", err)
		}
		// Delete blood request
		err = r.client.BloodSearchRequest.DeleteOneID(bloodReq.ID).Exec(ctx)
		if err != nil {
			return fmt.Errorf("не удалось удалить заявку на кровь: %w", err)
		}
	}

	// Мягкое удаление через хук SoftDeleteMixin
	err = r.client.Pet.DeleteOneID(id).Exec(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return fmt.Errorf("питомец с ID %s не найден", id)
		}
		return fmt.Errorf("не удалось удалить питомца: %w", err)
	}

	return nil
}

// ExistsByID проверяет, существует ли питомец с заданным ID
func (r *EntPetRepository) ExistsByID(ctx context.Context, id string) (bool, error) {
	if id == "" {
		return false, errors.New("неверный ID питомца")
	}

	exists, err := r.client.Pet.Query().
		Where(entpet.ID(id)).
		Exist(ctx)

	if err != nil {
		return false, fmt.Errorf("не удалось проверить существование питомца по ID %s: %w", id, err)
	}

	return exists, nil
}

// RestorePet восстанавливает мягко удаленного питомца, устанавливая deleted_at в NULL
func (r *EntPetRepository) RestorePet(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("неверный ID питомца")
	}

	// Используйте контекст SkipSoftDelete для поиска удаленной записи
	ctxWithSkip := schema.SkipSoftDelete(ctx)

	err := r.client.Pet.UpdateOneID(id).
		ClearDeletedAt().
		Exec(ctxWithSkip)

	if err != nil {
		if ent.IsNotFound(err) {
			return fmt.Errorf("питомец с ID %s не найден", id)
		}
		return fmt.Errorf("не удалось восстановить питомца: %w", err)
	}

	return nil
}

// GetDeletedPets извлекает всех мягко удаленных питомцев
func (r *EntPetRepository) GetDeletedPets(ctx context.Context) ([]*ent.Pet, error) {
	// Используйте контекст SkipSoftDelete для просмотра удаленных записей
	ctxWithSkip := schema.SkipSoftDelete(ctx)

	pets, err := r.client.Pet.Query().
		Where(entpet.DeletedAtNotNil()).
		WithHealth().
		WithTreatments().
		WithAnalyses().
		All(ctxWithSkip)

	if err != nil {
		return nil, fmt.Errorf("не удалось получить удаленных питомцев: %w", err)
	}

	return pets, nil
}

// AddPhotoURLs добавляет новые пути к фотографиям в массив PhotoUrls питомца
func (r *EntPetRepository) AddPhotoURLs(ctx context.Context, id string, paths []string) error {
	if id == "" {
		return errors.New("неверный ID питомца")
	}

	// Заменить URL-адреса фотографий на новые
	newPhotoUrls := paths

	// Обновить питомца
	err := r.client.Pet.UpdateOneID(id).
		SetPhotoUrls(newPhotoUrls).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("не удалось обновить URL-адреса фотографий питомца: %w", err)
	}

	return nil
}

func (r *EntPetRepository) GetPetsByBloodGroupAndRegion(ctx context.Context, bloodGroups, regions []string) ([]*model.Pet, error) {
	query := r.client.Pet.Query().
		WithOwner(func(uq *ent.UserQuery) {
			uq.Select(entuser.FieldFullName)
		}).
		WithBreedRef().
		WithHealth().
		WithTreatments().
		WithAnalyses().
		Where(
			entpet.BloodGroupIn(bloodGroups...),
		)
	if len(regions) > 0 {
		// Build predicates: any region contained in JSON array OR array length zero (including null)
		predicates := make([]predicate.DonorPreference, 0, len(regions)+1)
		for _, region := range regions {
			region := region // capture loop variable
			predicates = append(predicates, func(s *sql.Selector) {
				s.Where(sqljson.ValueContains(entdonorpreference.FieldPreferredLocationIds, region))
			})
		}
		// Add predicate for empty array or null
		predicates = append(predicates, func(s *sql.Selector) {
			s.Where(sql.Or(
				sqljson.LenEQ(entdonorpreference.FieldPreferredLocationIds, 0),
				sqljson.ValueIsNull(entdonorpreference.FieldPreferredLocationIds),
			))
		})
		query = query.Where(
			entpet.HasOwnerWith(entuser.HasDonorPreferenceWith(
				entdonorpreference.Or(predicates...),
			)),
		)
	}
	pets, err := query.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить питомцев: %w", err)
	}
	return domainmapper.PetToDomainSlice(pets), nil
}

func (r *EntPetRepository) CountSuitableDonors(ctx context.Context, bloodGroups []string) (int, error) {
	// count, err := r.client.Pet.Query().
	// 	Where(
	// 		func(s *sql.Selector) {
	// 			s.Where(sqljson.LenEQ(entpet.FieldStopFactors, 0))
	// 		},
	// 		entpet.BloodGroupIn(bloodGroups...),
	// 	).
	// 	Count(ctx)
	// if err != nil {
	// 	return 0, apperrors.Internal(err, "failed to count suitable donors")
	// }
	return 0, nil
}

// Exists проверяет существование питомца (алиас для ExistsByID для совместимости с PetReadRepository)
func (r *EntPetRepository) Exists(ctx context.Context, id string) (bool, error) {
	return r.ExistsByID(ctx, id)
}

// SetLastDonation обновляет дату последнего донорства питомца
func (r *EntPetRepository) SetLastDonation(ctx context.Context, petID string, lastDonationDate *time.Time) error {
	if petID == "" {
		return errors.New("неверный ID питомца")
	}

	updater := r.client.PetHealth.Update().Where(pethealth.HasOwnerWith(entpet.ID(petID)))

	if lastDonationDate != nil {
		updater.SetLastDonation(*lastDonationDate)
	} else {
		updater.ClearLastDonation()
	}

	_, err := updater.Save(ctx)
	if err != nil {
		return fmt.Errorf("не удалось обновить дату последнего донорства питомца: %w", err)
	}

	return nil
}

// SetTransfused обновляет флаг переливания крови питомца
func (r *EntPetRepository) SetTransfused(ctx context.Context, petID string, transfused bool) error {
	if petID == "" {
		return errors.New("неверный ID питомца")
	}

	updater := r.client.PetHealth.Update().Where(pethealth.HasOwnerWith(entpet.ID(petID)))

	updater.SetTransfused(transfused)

	_, err := updater.Save(ctx)
	if err != nil {
		return fmt.Errorf("не удалось обновить флаг переливания питомца: %w", err)
	}

	return nil
}
