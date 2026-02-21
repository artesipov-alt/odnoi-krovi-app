package pg

import (
	"context"
	"errors"
	"fmt"

	"github.com/artesipov-alt/odnoi-krovi-app/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/pet"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/petanalysis"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/petbonus"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/pethealth"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/pettreatment"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/schema"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/apperrors"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/services"
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
func (r *EntPetRepository) Create(ctx context.Context, input *ent.CreatePetInput, healthInput *ent.CreatePetHealthInput, treatmentsInput *ent.CreatePetTreatmentInput, analysesInput []*ent.CreatePetAnalysisInput, bonusesInput *ent.CreatePetBonusInput) (*ent.Pet, error) {
	if input == nil {
		return nil, errors.New("входные данные питомца не могут быть nil")
	}

	// Начать транзакцию
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("не удалось начать транзакцию: %w", err)
	}

	// 1. Создать питомца
	petCreate := tx.Pet.Create().
		SetInput(*input)

	newPet, err := petCreate.Save(ctx)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("не удалось создать питомца: %w", err)
	}

	// 2. Создать связанные сущности, если они предоставлены
	if healthInput != nil {
		_, err = tx.PetHealth.Create().
			SetInput(*healthInput).
			SetOwner(newPet).
			Save(ctx)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("не удалось создать данные о здоровье питомца: %w", err)
		}
	}

	if treatmentsInput != nil {
		_, err = tx.PetTreatment.Create().
			SetInput(*treatmentsInput).
			SetOwner(newPet).
			Save(ctx)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("не удалось создать данные о лечении питомца: %w", err)
		}
	}

	for _, a := range analysesInput {
		builder := tx.PetAnalysis.Create().
			SetInput(*a).
			SetPetID(newPet.ID)

		_, err = builder.Save(ctx)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("не удалось создать анализ питомца: %w", err)
		}
	}

	if bonusesInput != nil {
		_, err = tx.PetBonus.Create().
			SetInput(*bonusesInput).
			SetOwner(newPet).
			Save(ctx)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("не удалось создать бонусы питомца: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("не удалось зафиксировать транзакцию: %w", err)
	}

	query := r.client.Pet.Query().Where(pet.ID(newPet.ID)).WithBreedRef().WithOwner()
	return query.Only(ctx)
}

// GetPet возвращает питомца по его ID с возможностью предварительной загрузки связанных данных
func (r *EntPetRepository) GetPet(ctx context.Context, id string, opts services.PetPreloadOptions) (*ent.Pet, error) {
	pquery := r.client.Pet.Query().Where(pet.ID(id)).WithBreedRef().WithBloodGroupRef()

	// Применяем опции предварительной загрузки
	if opts.WithAll {
		pquery = pquery.WithHealth().WithTreatments().WithAnalyses().WithBonuses()
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
		if opts.WithBonuses {
			pquery = pquery.WithBonuses()
		}
	}

	pet, err := pquery.Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, apperrors.ErrPetNotFound
		}
		return nil, apperrors.Internal(err, "не удалось получить питомца")
	}
	return pet, nil
}

// GetPetsByUser возвращает запрос для предварительной загрузки питомцев по ID пользователя
func (r *EntPetRepository) GetPetsByUser(ctx context.Context, userID string, opts services.PetPreloadOptions) ([]*ent.Pet, error) {
	pquery := r.client.Pet.Query().Where(pet.UserID(userID)).WithBreedRef().WithBloodGroupRef()

	// Применяем опции предварительной загрузки
	if opts.WithAll {
		pquery = pquery.WithHealth().WithTreatments().WithAnalyses().WithBonuses()
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
		if opts.WithBonuses {
			pquery = pquery.WithBonuses()
		}
	}
	pets, err := pquery.All(ctx)
	if err != nil {
		return nil, apperrors.Internal(err, "failed to get pets")
	}
	return pets, nil
}

// Update обновляет существующего питомца и его связанные сущности в транзакции
func (r *EntPetRepository) Update(ctx context.Context, id string, petInput *ent.UpdatePetInput, healthInput *ent.UpdatePetHealthInput, treatmentsInput *ent.UpdatePetTreatmentInput, analysesInput []*ent.UpdatePetAnalysisInput, bonusesInput *ent.UpdatePetBonusInput) (*ent.Pet, error) {
	if id == "" {
		return nil, errors.New("неверный ID питомца")
	}

	tx, err := r.client.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("не удалось начать транзакцию: %w", err)
	}

	// 1. Обновить питомца
	if petInput != nil {
		err = tx.Pet.UpdateOneID(id).
			SetInput(*petInput).
			Exec(ctx)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("не удалось обновить питомца: %w", err)
		}
	}

	// 2. Обновить данные о здоровье
	if healthInput != nil {
		existingHealth, err := tx.PetHealth.Query().Where(pethealth.HasOwnerWith(pet.ID(id))).Only(ctx)
		if err != nil && !ent.IsNotFound(err) {
			tx.Rollback()
			return nil, err
		}
		if existingHealth != nil {
			update := tx.PetHealth.UpdateOne(existingHealth)
			if healthInput.HealthStatus != nil {
				update.SetHealthStatus(*healthInput.HealthStatus)
			} else {
				update.ClearHealthStatus()
			}
			if healthInput.Transfused != nil {
				update.SetTransfused(*healthInput.Transfused)
			} else {
				update.ClearTransfused()
			}
			if healthInput.Medications != nil {
				update.SetMedications(*healthInput.Medications)
			} else {
				update.ClearMedications()
			}
			if healthInput.SurgicalInterventions != nil {
				update.SetSurgicalInterventions(*healthInput.SurgicalInterventions)
			} else {
				update.ClearSurgicalInterventions()
			}
			_, err = update.Save(ctx)
		} else {
			builder := tx.PetHealth.Create().SetOwnerID(id)
			if healthInput.HealthStatus != nil {
				builder.SetHealthStatus(*healthInput.HealthStatus)
			}
			if healthInput.Transfused != nil {
				builder.SetTransfused(*healthInput.Transfused)
			}
			if healthInput.Medications != nil {
				builder.SetMedications(*healthInput.Medications)
			}
			if healthInput.SurgicalInterventions != nil {
				builder.SetSurgicalInterventions(*healthInput.SurgicalInterventions)
			}
			_, err = builder.Save(ctx)
		}
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("не удалось обновить данные о здоровье питомца: %w", err)
		}
	}

	// 3. Обновить данные о лечении
	if treatmentsInput != nil {
		existingTreatment, err := tx.PetTreatment.Query().Where(pettreatment.HasOwnerWith(pet.ID(id))).Only(ctx)
		if err != nil && !ent.IsNotFound(err) {
			tx.Rollback()
			return nil, err
		}
		if existingTreatment != nil {
			update := tx.PetTreatment.UpdateOne(existingTreatment)
			if treatmentsInput.RabiesVaccinationDate != nil {
				update.SetRabiesVaccinationDate(*treatmentsInput.RabiesVaccinationDate)
			} else {
				update.ClearRabiesVaccinationDate()
			}
			if treatmentsInput.InfectionVaccinationDate != nil {
				update.SetInfectionVaccinationDate(*treatmentsInput.InfectionVaccinationDate)
			} else {
				update.ClearInfectionVaccinationDate()
			}
			if treatmentsInput.EctoparasiteTreatmentDate != nil {
				update.SetEctoparasiteTreatmentDate(*treatmentsInput.EctoparasiteTreatmentDate)
			} else {
				update.ClearEctoparasiteTreatmentDate()
			}
			if treatmentsInput.DewormingDate != nil {
				update.SetDewormingDate(*treatmentsInput.DewormingDate)
			} else {
				update.ClearDewormingDate()
			}
			_, err = update.Save(ctx)
		} else {
			builder := tx.PetTreatment.Create().SetOwnerID(id)
			if treatmentsInput.RabiesVaccinationDate != nil {
				builder.SetRabiesVaccinationDate(*treatmentsInput.RabiesVaccinationDate)
			}
			if treatmentsInput.InfectionVaccinationDate != nil {
				builder.SetInfectionVaccinationDate(*treatmentsInput.InfectionVaccinationDate)
			}
			if treatmentsInput.EctoparasiteTreatmentDate != nil {
				builder.SetEctoparasiteTreatmentDate(*treatmentsInput.EctoparasiteTreatmentDate)
			}
			if treatmentsInput.DewormingDate != nil {
				builder.SetDewormingDate(*treatmentsInput.DewormingDate)
			}
			_, err = builder.Save(ctx)
		}
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("не удалось обновить данные о лечении питомца: %w", err)
		}
	}

	// Для анализов, удалить существующие и создать новые, если предоставлены
	if analysesInput != nil {
		// Удалить существующие анализы для этого питомца (жесткое удаление)
		_, err = tx.PetAnalysis.Delete().Where(petanalysis.PetID(id)).Exec(schema.SkipSoftDelete(ctx))
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("не удалось удалить существующие анализы питомца: %w", err)
		}

		// Создать новые анализы
		for _, a := range analysesInput {
			builder := tx.PetAnalysis.Create().
				SetPetID(id)

			if a.AnalysisName != nil {
				builder.SetAnalysisName(*a.AnalysisName)
			}
			if a.AnalysisType != nil {
				builder.SetAnalysisType(*a.AnalysisType)
			}
			if a.AnalysisDate != nil {
				builder.SetAnalysisDate(*a.AnalysisDate)
			}

			_, err = builder.Save(ctx)
			if err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("не удалось создать анализ питомца: %w", err)
			}
		}
	}

	// 5. Обновить бонусы
	if bonusesInput != nil {
		existingBonus, err := tx.PetBonus.Query().Where(petbonus.HasOwnerWith(pet.ID(id))).Only(ctx)
		if err != nil && !ent.IsNotFound(err) {
			tx.Rollback()
			return nil, err
		}
		if existingBonus != nil {
			update := tx.PetBonus.UpdateOne(existingBonus)
			if bonusesInput.IsArtist != nil {
				update.SetIsArtist(*bonusesInput.IsArtist)
			} else {
				update.SetIsArtist(false)
			}
			if bonusesInput.IsTherapist != nil {
				update.SetIsTherapist(*bonusesInput.IsTherapist)
			} else {
				update.SetIsTherapist(false)
			}
			if bonusesInput.IsFormerDonor != nil {
				update.SetIsFormerDonor(*bonusesInput.IsFormerDonor)
			} else {
				update.SetIsFormerDonor(false)
			}
			if bonusesInput.IsGuideDog != nil {
				update.SetIsGuideDog(*bonusesInput.IsGuideDog)
			} else {
				update.SetIsGuideDog(false)
			}
			_, err = update.Save(ctx)
		} else {
			builder := tx.PetBonus.Create().SetOwnerID(id)
			if bonusesInput.IsArtist != nil {
				builder.SetIsArtist(*bonusesInput.IsArtist)
			}
			if bonusesInput.IsTherapist != nil {
				builder.SetIsTherapist(*bonusesInput.IsTherapist)
			}
			if bonusesInput.IsFormerDonor != nil {
				builder.SetIsFormerDonor(*bonusesInput.IsFormerDonor)
			}
			if bonusesInput.IsGuideDog != nil {
				builder.SetIsGuideDog(*bonusesInput.IsGuideDog)
			}
			_, err = builder.Save(ctx)
		}
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("не удалось обновить бонусы питомца: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("не удалось зафиксировать транзакцию: %w", err)
	}

	query := r.client.Pet.Query().Where(pet.ID(id)).WithBreedRef().WithOwner()
	return query.Only(ctx)
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
		Where(pet.ID(id)).
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
		Where(pet.DeletedAtNotNil()).
		WithHealth().
		WithTreatments().
		WithAnalyses().
		WithBonuses().
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

	// Получить текущие URL-адреса фотографий
	p, err := r.client.Pet.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("не удалось получить питомца для обновления фото: %w", err)
	}

	// Добавить новые пути
	newPhotoUrls := append(p.PhotoUrls, paths...)

	// Обновить питомца
	err = r.client.Pet.UpdateOneID(id).
		SetPhotoUrls(newPhotoUrls).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("не удалось обновить URL-адреса фотографий питомца: %w", err)
	}

	return nil
}

// UpdateStatus обновляет статус питомца по его ID
func (r *EntPetRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	if id == "" {
		return errors.New("неверный ID питомца")
	}

	err := r.client.Pet.UpdateOneID(id).
		SetPetStatus(pet.PetStatus(status)).
		Exec(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return fmt.Errorf("питомец с ID %s не найден", id)
		}
		return fmt.Errorf("не удалось обновить статус питомца: %w", err)
	}

	return nil
}
