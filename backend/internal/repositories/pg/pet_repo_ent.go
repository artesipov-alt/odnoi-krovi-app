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
)

// EntPetRepository implements PetRepository using ENT
type EntPetRepository struct {
	client *ent.Client
}

// NewEntPetRepository creates a new ENT pet repository
func NewEntPetRepository(client *ent.Client) *EntPetRepository {
	return &EntPetRepository{
		client: client,
	}
}

// Create creates a new pet in the database along with its related entities in a transaction
func (r *EntPetRepository) Create(ctx context.Context, input *ent.CreatePetInput, healthInput *ent.CreatePetHealthInput, treatmentsInput *ent.CreatePetTreatmentInput, analysesInput []*ent.CreatePetAnalysisInput, bonusesInput *ent.CreatePetBonusInput) (*ent.Pet, error) {
	if input == nil {
		return nil, errors.New("pet input cannot be nil")
	}

	// Start transaction
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}

	// 1. Create Pet
	petCreate := tx.Pet.Create().
		SetInput(*input)

	newPet, err := petCreate.Save(ctx)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create pet: %w", err)
	}

	// 2. Create related entities if provided
	if healthInput != nil {
		_, err = tx.PetHealth.Create().
			SetInput(*healthInput).
			SetOwner(newPet).
			Save(ctx)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create pet health: %w", err)
		}
	}

	if treatmentsInput != nil {
		_, err = tx.PetTreatment.Create().
			SetInput(*treatmentsInput).
			SetOwner(newPet).
			Save(ctx)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create pet treatment: %w", err)
		}
	}

	for _, a := range analysesInput {
		builder := tx.PetAnalysis.Create().
			SetInput(*a).
			SetOwnerID(newPet.ID)

		_, err = builder.Save(ctx)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create pet analysis: %w", err)
		}
	}

	if bonusesInput != nil {
		_, err = tx.PetBonus.Create().
			SetInput(*bonusesInput).
			SetOwner(newPet).
			Save(ctx)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create pet bonus: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	query := r.client.Pet.Query().Where(pet.ID(newPet.ID)).WithBreedRef().WithOwner()
	return query.Only(ctx)
}

// GetPetQuery returns a query for eager loading
func (r *EntPetRepository) GetPetQuery(ctx context.Context, id string) *ent.PetQuery {
	return r.client.Pet.Query().Where(pet.ID(id)).WithBreedRef().WithOwner()
}

// GetPetsQueryByUser returns a query for eager loading pets by user ID
func (r *EntPetRepository) GetPetsQueryByUser(ctx context.Context, userID string) *ent.PetQuery {
	return r.client.Pet.Query().Where(pet.UserID(userID)).WithBreedRef()
}

// Update updates an existing pet and its related entities in a transaction
func (r *EntPetRepository) Update(ctx context.Context, id string, petInput *ent.UpdatePetInput, healthInput *ent.UpdatePetHealthInput, treatmentsInput *ent.UpdatePetTreatmentInput, analysesInput []*ent.UpdatePetAnalysisInput, bonusesInput *ent.UpdatePetBonusInput) (*ent.Pet, error) {
	if id == "" {
		return nil, errors.New("invalid pet ID")
	}

	tx, err := r.client.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}

	// 1. Update Pet
	if petInput != nil {
		err = tx.Pet.UpdateOneID(id).
			SetInput(*petInput).
			Exec(ctx)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to update pet: %w", err)
		}
	}

	// 2. Update health
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
			return nil, fmt.Errorf("failed to update pet health: %w", err)
		}
	}

	// 3. Update treatments
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
			return nil, fmt.Errorf("failed to update pet treatment: %w", err)
		}
	}

	// For analyses, delete existing ones and create new ones if provided
	if analysesInput != nil {
		// Delete existing analyses for this pet (hard delete)
		_, err = tx.PetAnalysis.Delete().Where(petanalysis.HasOwnerWith(pet.ID(id))).Exec(schema.SkipSoftDelete(ctx))
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to delete existing pet analyses: %w", err)
		}

		// Create new analyses
		for _, a := range analysesInput {
			builder := tx.PetAnalysis.Create().
				SetOwnerID(id)

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
				return nil, fmt.Errorf("failed to create pet analysis: %w", err)
			}
		}
	}

	// 5. Update bonuses
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
			return nil, fmt.Errorf("failed to update pet bonus: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	query := r.client.Pet.Query().Where(pet.ID(id)).WithBreedRef().WithOwner()
	return query.Only(ctx)
}

// Delete deletes a pet by their ID (soft delete)
func (r *EntPetRepository) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("invalid pet ID")
	}

	// Soft delete via SoftDeleteMixin hook
	err := r.client.Pet.DeleteOneID(id).Exec(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return fmt.Errorf("pet with id %s not found", id)
		}
		return fmt.Errorf("failed to delete pet: %w", err)
	}

	return nil
}

// ExistsByID checks if a pet with the given ID exists
func (r *EntPetRepository) ExistsByID(ctx context.Context, id string) (bool, error) {
	if id == "" {
		return false, errors.New("invalid pet ID")
	}

	exists, err := r.client.Pet.Query().
		Where(pet.ID(id)).
		Exist(ctx)

	if err != nil {
		return false, fmt.Errorf("failed to check pet existence by id %s: %w", id, err)
	}

	return exists, nil
}

// RestorePet restores a soft-deleted pet by setting deleted_at to NULL
func (r *EntPetRepository) RestorePet(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("invalid pet ID")
	}

	// Use SkipSoftDelete context to find the deleted record
	ctxWithSkip := schema.SkipSoftDelete(ctx)

	err := r.client.Pet.UpdateOneID(id).
		ClearDeletedAt().
		Exec(ctxWithSkip)

	if err != nil {
		if ent.IsNotFound(err) {
			return fmt.Errorf("pet with id %s not found", id)
		}
		return fmt.Errorf("failed to restore pet: %w", err)
	}

	return nil
}

// GetDeletedPets retrieves all soft-deleted pets
func (r *EntPetRepository) GetDeletedPets(ctx context.Context) ([]*ent.Pet, error) {
	// Use SkipSoftDelete context to see deleted records
	ctxWithSkip := schema.SkipSoftDelete(ctx)

	pets, err := r.client.Pet.Query().
		Where(pet.DeletedAtNotNil()).
		WithHealth().
		WithTreatments().
		WithAnalyses().
		WithBonuses().
		All(ctxWithSkip)

	if err != nil {
		return nil, fmt.Errorf("failed to get deleted pets: %w", err)
	}

	return pets, nil
}

// AddPhotoURLs adds new photo paths to the pet's PhotoUrls array
func (r *EntPetRepository) AddPhotoURLs(ctx context.Context, id string, paths []string) error {
	if id == "" {
		return errors.New("invalid pet ID")
	}

	// Fetch current photo URLs
	p, err := r.client.Pet.Get(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get pet for photo update: %w", err)
	}

	// Append new paths
	newPhotoUrls := append(p.PhotoUrls, paths...)

	// Update pet
	err = r.client.Pet.UpdateOneID(id).
		SetPhotoUrls(newPhotoUrls).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to update pet photo URLs: %w", err)
	}

	return nil
}

// UpdateStatus обновляет статус питомца по его ID
func (r *EntPetRepository) UpdateStatus(ctx context.Context, id string, status string) error {
	if id == "" {
		return errors.New("invalid pet ID")
	}

	err := r.client.Pet.UpdateOneID(id).
		SetPetStatus(pet.PetStatus(status)).
		Exec(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return fmt.Errorf("pet with id %s not found", id)
		}
		return fmt.Errorf("failed to update pet status: %w", err)
	}

	return nil
}

// UpdateStatusWithTx обновляет статус питомца по его ID в рамках транзакции
func (r *EntPetRepository) UpdateStatusWithTx(ctx context.Context, tx *ent.Tx, id string, status string) error {
	if id == "" {
		return errors.New("invalid pet ID")
	}

	err := tx.Pet.UpdateOneID(id).
		SetPetStatus(pet.PetStatus(status)).
		Exec(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return fmt.Errorf("pet with id %s not found", id)
		}
		return fmt.Errorf("failed to update pet status: %w", err)
	}

	return nil
}
