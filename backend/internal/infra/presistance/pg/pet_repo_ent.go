package pg

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqljson"

	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/domain"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/bloodgroup"
	entpet "github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/petanalysis"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/pethealth"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/pettreatment"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/infra/ent/schema"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/pet/model"

	"github.com/jinzhu/copier"
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

	var petInput ent.CreatePetInput
	if err := copier.Copy(&petInput, petDomain); err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("не удалось скопировать данные питомца: %w", err)
	}

	if petDomain.BloodGroupName != nil {
		bg, err := tx.BloodGroup.Query().Where(bloodgroup.BloodGroupEQ(*petDomain.BloodGroupName)).Only(ctx)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("не удалось найти группу крови: %w", err)
		}
		petInput.BloodGroupRefID = &bg.ID
	}

	// 1. Создать питомца
	newPet, err := tx.Pet.Create().
		SetInput(petInput).
		Save(ctx)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("не удалось создать питомца: %w", err)
	}

	// 2. Создать связанные сущности, если они предоставлены
	if petDomain.Health != nil {
		var healthInput ent.CreatePetHealthInput
		hs := pethealth.HealthStatus(petDomain.Health.HealthStatus)
		healthInput.HealthStatus = &hs
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
			SetOwner(newPet).
			Save(ctx)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("не удалось создать данные о здоровье питомца: %w", err)
		}
	}

	if petDomain.Treatments != nil {
		var treatmentsInput ent.CreatePetTreatmentInput
		if err := copier.Copy(&treatmentsInput, petDomain.Treatments); err != nil {
			return nil, fmt.Errorf("не удалось скопировать данные лечения питомца: %w", err)
		}
		_, err = tx.PetTreatment.Create().
			SetInput(treatmentsInput).
			SetOwner(newPet).
			Save(ctx)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("не удалось создать данные о лечении питомца: %w", err)
		}
	}

	for _, a := range petDomain.Analyses {
		var analysisInput ent.CreatePetAnalysisInput
		if err := copier.Copy(&analysisInput, a); err != nil {
			return nil, fmt.Errorf("не удалось скопировать данные анализа питомца: %w", err)
		}
		builder := tx.PetAnalysis.Create().
			SetInput(analysisInput).
			SetPetID(newPet.ID)

		_, err = builder.Save(ctx)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("не удалось создать анализ питомца: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("не удалось зафиксировать транзакцию: %w", err)
	}

	result := &model.Pet{
		ID:        newPet.ID,
		CreatedAt: &newPet.CreatedAt,
	}

	return result, nil
}

// GetPet возвращает питомца по его ID с возможностью предварительной загрузки связанных данных
func (r *EntPetRepository) GetPet(ctx context.Context, id string, opts domain.PetPreloadOptions) (*model.Pet, error) {
	pquery := r.client.Pet.Query().Where(entpet.ID(id)).WithBreedRef().WithBloodGroupRef()

	// Применяем опции предварительной загрузки
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
		return nil, apperrors.Internal(err, "не удалось получить питомца")
	}

	var result model.Pet
	if err := copier.Copy(&result, entPet); err != nil {
		return nil, fmt.Errorf("не удалось скопировать данные питомца: %w", err)
	}

	if entPet.Edges.BloodGroupRef != nil {
		result.BloodGroupName = &entPet.Edges.BloodGroupRef.BloodGroup
	}

	if entPet.Edges.BreedRef != nil {
		result.BreedRefID = &entPet.Edges.BreedRef.ID
	}

	if entPet.Edges.Health != nil {
		result.Health = &model.PetHealth{
			HealthStatus:          model.HealthStatus(entPet.Edges.Health.HealthStatus),
			LastDonation:          entPet.Edges.Health.LastDonation,
			Transfused:            &entPet.Edges.Health.Transfused,
			Medications:           &entPet.Edges.Health.Medications,
			SurgicalInterventions: &entPet.Edges.Health.SurgicalInterventions,
		}
	}
	if entPet.Edges.Treatments != nil {
		var treatments model.PetTreatment
		if err := copier.Copy(&treatments, entPet.Edges.Treatments); err == nil {
			result.Treatments = &treatments
		}
	}
	if len(entPet.Edges.Analyses) > 0 {
		var analyses []*model.PetAnalysis
		if err := copier.Copy(&analyses, entPet.Edges.Analyses); err == nil {
			result.Analyses = analyses
		}
	}

	return &result, nil
}

// GetPetsByUser возвращает запрос для предварительной загрузки питомцев по ID пользователя
func (r *EntPetRepository) GetPetsByUser(ctx context.Context, userID string, opts domain.PetPreloadOptions) ([]*model.Pet, error) {
	pquery := r.client.Pet.Query().Where(entpet.UserID(userID)).WithBreedRef().WithBloodGroupRef()

	// Применяем опции предварительной загрузки
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

	var result []*model.Pet
	for _, p := range pets {
		var petDomain model.Pet
		if err := copier.Copy(&petDomain, p); err != nil {
			return nil, fmt.Errorf("не удалось скопировать данные питомца: %w", err)
		}

		if p.Edges.BloodGroupRef != nil {
			petDomain.BloodGroupName = &p.Edges.BloodGroupRef.BloodGroup
		}

		if p.Edges.BreedRef != nil {
			petDomain.BreedRefID = &p.Edges.BreedRef.ID
		}

		if p.Edges.Health != nil {
			petDomain.Health = &model.PetHealth{
				HealthStatus:          model.HealthStatus(p.Edges.Health.HealthStatus),
				LastDonation:          p.Edges.Health.LastDonation,
				Transfused:            &p.Edges.Health.Transfused,
				Medications:           &p.Edges.Health.Medications,
				SurgicalInterventions: &p.Edges.Health.SurgicalInterventions,
			}
		}
		if p.Edges.Treatments != nil {
			var treatments model.PetTreatment
			if err := copier.Copy(&treatments, p.Edges.Treatments); err == nil {
				petDomain.Treatments = &treatments
			}
		}
		if len(p.Edges.Analyses) > 0 {
			var analyses []*model.PetAnalysis
			if err := copier.Copy(&analyses, p.Edges.Analyses); err == nil {
				petDomain.Analyses = analyses
			}
		}
		result = append(result, &petDomain)
	}
	return result, nil
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
	if petDomain.StopFactors != nil {
		updater.SetStopFactors(petDomain.StopFactors)
	}
	if petDomain.WarnFactors != nil {
		updater.SetWarnFactors(petDomain.WarnFactors)
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

	updatedPetEntity, err := updater.Save(ctx)
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
			var createInput ent.CreatePetTreatmentInput
			if err := copier.Copy(&createInput, petDomain.Treatments); err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("не удалось скопировать данные лечения питомца: %w", err)
			}
			_, err = tx.PetTreatment.Create().
				SetInput(createInput).
				SetOwnerID(id).
				Save(ctx)
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
			var analysisInput ent.CreatePetAnalysisInput
			if err := copier.Copy(&analysisInput, a); err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("не удалось скопировать данные анализа питомца: %w", err)
			}
			_, err = tx.PetAnalysis.Create().
				SetInput(analysisInput).
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

	return &model.Pet{ID: updatedPetEntity.ID, UpdatedAt: &updatedPetEntity.UpdatedAt}, nil
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
	count, err := r.client.Pet.Query().
		Where(
			func(s *sql.Selector) {
				s.Where(sqljson.LenEQ(entpet.FieldStopFactors, 0))
			},
			entpet.HasBloodGroupRefWith(bloodgroup.BloodGroupIn(bloodGroups...)),
		).
		Count(ctx)
	if err != nil {
		return 0, apperrors.Internal(err, "failed to count suitable donors")
	}
	return count, nil
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
