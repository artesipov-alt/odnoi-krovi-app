package pg

import (
	"context"
	"errors"
	"fmt"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain/pet/model"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/bloodgroup"
	entbloodreq "github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/bloodsearchrequest"
	entpet "github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/petanalysis"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/pethealth"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/pettreatment"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/schema"
)

// petToDomain converts ent.Pet to domain model.Pet
func petToDomain(e *ent.Pet) *model.Pet {
	if e == nil {
		return nil
	}

	pet := &model.Pet{
		ID:                 e.ID,
		Name:               e.Name,
		Type:               model.PetType(e.Type),
		PetStatus:          model.PetStatusNone,
		WeightKg:           e.WeightKg,
		Gender:             model.Gender(e.Gender),
		BirthDate:          e.BirthDate,
		ChipNumber:         e.ChipNumber,
		PhotoURLs:          e.PhotoUrls,
		LivingCondition:    model.LivingCondition(e.LivingCondition),
		ReproductiveStatus: model.ReproductiveStatus(e.ReproductiveStatus),
		OwnerID:            e.UserID,
		BreedRefID:         e.BreedID,
		Bonuses:            e.Bonuses,
		CreatedAt:          &e.CreatedAt,
		UpdatedAt:          &e.UpdatedAt,
		DeletedAt:          e.DeletedAt,
	}

	// Map BloodGroupName from edge if available
	if e.Edges.BloodGroupRef != nil {
		pet.BloodGroupName = &e.Edges.BloodGroupRef.BloodGroup
	}

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

	if e.Edges.BloodSearchRequest != nil {
		pet.PetStatus = model.PetStatusRecipient
		if len(e.Edges.BloodSearchRequest.Edges.Responses) > 0 {
			pet.PetStatus = model.PetStatusBloodFound
		}
	}

	return pet
}

// petToDomainSlice converts slice of ent.Pet to slice of domain model.Pet
func petToDomainSlice(pets []*ent.Pet) []*model.Pet {
	if pets == nil {
		return nil
	}
	result := make([]*model.Pet, len(pets))
	for i, p := range pets {
		result[i] = petToDomain(p)
	}
	return result
}

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

	if petDomain.BloodGroupName != nil {
		bg, err := tx.BloodGroup.Query().Where(bloodgroup.BloodGroupEQ(*petDomain.BloodGroupName)).Only(ctx)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("не удалось найти группу крови: %w", err)
		}
		builder.SetBloodGroupRefID(bg.ID)
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

// GetPet возвращает питомца по его ID с возможностью предварительной загрузки связанных данных
// Deprecated: используйте GetByID
func (r *EntPetRepository) GetPet(ctx context.Context, id string, opts pet.PetPreloadOptions) (*model.Pet, error) {
	return r.GetByID(ctx, id, opts)
}

// GetPetsByUser возвращает запрос для предварительной загрузки питомцев по ID пользователя
// Deprecated: используйте GetByUserID
func (r *EntPetRepository) GetPetsByUser(ctx context.Context, userID string, opts pet.PetPreloadOptions) ([]*model.Pet, error) {
	return r.GetByUserID(ctx, userID, opts)
}

// GetByID получает питомца по ID с опциями загрузки связанных данных
func (r *EntPetRepository) GetByID(ctx context.Context, id string, opts pet.PetPreloadOptions) (*model.Pet, error) {
	pquery := r.client.Pet.Query().Where(entpet.ID(id)).WithBreedRef().WithBloodGroupRef()

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
	if opts.WithBloodReq {
		pquery = pquery.WithBloodSearchRequest(func(bsrq *ent.BloodSearchRequestQuery) {
			bsrq.Where(entbloodreq.StatusEQ(entbloodreq.DefaultStatus)).WithResponses()
		})
	}

	entPet, err := pquery.Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperrors.ErrPetNotFound
		}
		return nil, apperrors.Internal(err, "failed to get pet")
	}

	return petToDomain(entPet), nil
}

// GetByUserID получает всех питомцев пользователя
func (r *EntPetRepository) GetByUserID(ctx context.Context, userID string, opts pet.PetPreloadOptions) ([]*model.Pet, error) {
	pquery := r.client.Pet.Query().Where(entpet.UserID(userID)).WithBreedRef().WithBloodGroupRef()

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

	if opts.WithBloodReq {
		pquery = pquery.WithBloodSearchRequest(func(bsrq *ent.BloodSearchRequestQuery) {
			bsrq.Where(entbloodreq.StatusEQ(entbloodreq.DefaultStatus)).WithResponses()
		})
	}

	pets, err := pquery.All(ctx)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get pets")
	}

	return petToDomainSlice(pets), nil
}

// Update обновляет существующего питомца и его связанные сущности в транзакции
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
	if petDomain.BloodGroupName != nil {
		bg, err := tx.BloodGroup.Query().Where(bloodgroup.BloodGroupEQ(*petDomain.BloodGroupName)).Only(ctx)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("blood group not found: %w", err)
		}
		updater.SetBloodGroupRefID(bg.ID)
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
			var healthInput ent.CreatePetHealthInput
			healthInput.HealthStatus = new(pethealth.HealthStatus(petDomain.Health.HealthStatus))
			if petDomain.Health.LastDonation != nil {
				healthInput.LastDonation = petDomain.Health.LastDonation
			}
			if petDomain.Health.Transfused != nil {
				healthInput.Transfused = petDomain.Health.Transfused
			}
			if petDomain.Health.Medications != nil {
				healthInput.Medications = petDomain.Health.Medications
			}
			if petDomain.Health.SurgicalInterventions != nil {
				healthInput.SurgicalInterventions = petDomain.Health.SurgicalInterventions
			}
			_, err = tx.PetHealth.Create().
				SetInput(healthInput).
				SetOwnerID(id).
				Save(ctx)
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
func (r *EntPetRepository) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("неверный ID питомца")
	}

	// Мягкое удаление через хук SoftDeleteMixin
	err := r.client.Pet.DeleteOneID(id).Exec(ctx)

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

func (r *EntPetRepository) CountSuitableDonors(ctx context.Context, bloodGroups []string) (int, error) {
	// count, err := r.client.Pet.Query().
	// 	Where(
	// 		func(s *sql.Selector) {
	// 			s.Where(sqljson.LenEQ(entpet.FieldStopFactors, 0))
	// 		},
	// 		entpet.HasBloodGroupRefWith(bloodgroup.BloodGroupIn(bloodGroups...)),
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

// UpdateStatus обновляет статус питомца по его ID
// func (r *EntPetRepository) UpdateStatus(ctx context.Context, id string, status string) error {
// 	if id == "" {
// 		return errors.New("неверный ID питомца")
// 	}

// 	err := r.client.Pet.UpdateOneID(id).
// 		SetPetStatus(pet.PetStatus(status)).
// 		Exec(ctx)

// 	if err != nil {
// 		if ent.IsNotFound(err) {
// 			return fmt.Errorf("питомец с ID %s не найден", id)
// 		}
// 		return fmt.Errorf("не удалось обновить статус питомца: %w", err)
// 	}

// 	return nil
// }
