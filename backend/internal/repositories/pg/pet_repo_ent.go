package pg

import (
	"context"
	"errors"
	"fmt"

	"github.com/artesipov-alt/odnoi-krovi-app/ent"
	"github.com/artesipov-alt/odnoi-krovi-app/ent/pet"
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

// nillable возвращает указатель на значение, если оно не является нулевым для своего типа, иначе nil.
func nillable[T comparable](v T) *T {
	var zero T
	if v == zero {
		return nil
	}
	return &v
}

// Create creates a new pet in the database along with its related entities in a transaction
func (r *EntPetRepository) Create(ctx context.Context, p *ent.Pet, health *ent.PetHealth, treatments *ent.PetTreatment, analyses []*ent.PetAnalysis, bonuses *ent.PetBonus) (*ent.Pet, error) {
	if p == nil {
		return nil, errors.New("pet cannot be nil")
	}

	// Start transaction
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}

	// 1. Create Pet
	petCreate := tx.Pet.Create().
		SetName(p.Name).
		SetType(p.Type).
		SetPetStatus(p.PetStatus).
		SetNillableWeightKg(nillable(p.WeightKg)).
		SetNillableBloodGroup(nillable(p.BloodGroup)).
		SetNillableGender(nillable(p.Gender)).
		SetNillableAgeYears(nillable(p.AgeYears)).
		SetNillableAgeMonths(nillable(p.AgeMonths)).
		SetNillableBirthDate(p.BirthDate).
		SetNillableChipNumber(nillable(p.ChipNumber)).
		SetNillablePhotoURL(nillable(p.PhotoURL)).
		SetNillableBreedID(nillable(p.BreedID)).
		SetNillableUserID(nillable(p.UserID)).
		SetNillableLivingCondition(nillable(p.LivingCondition))

	newPet, err := petCreate.Save(ctx)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create pet: %w", err)
	}

	// 2. Create related entities if provided
	if health != nil {
		_, err = tx.PetHealth.Create().
			SetOwner(newPet).
			SetNillableReproductiveStatus(nillable(health.ReproductiveStatus)).
			SetNillableHealthStatus(nillable(health.HealthStatus)).
			SetNillableLastDonation(health.LastDonation).
			SetTransfused(health.Transfused).
			SetMedications(health.Medications).
			SetSurgicalInterventions(health.SurgicalInterventions).
			Save(ctx)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create pet health: %w", err)
		}
	}

	if treatments != nil {
		_, err = tx.PetTreatment.Create().
			SetOwner(newPet).
			SetNillableRabiesVaccinationDate(treatments.RabiesVaccinationDate).
			SetNillableInfectionVaccinationDate(treatments.InfectionVaccinationDate).
			SetNillableEctoparasiteTreatmentDate(treatments.EctoparasiteTreatmentDate).
			SetNillableDewormingDate(treatments.DewormingDate).
			Save(ctx)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create pet treatment: %w", err)
		}
	}

	for _, a := range analyses {
		_, err = tx.PetAnalysis.Create().
			SetOwnerID(newPet.ID).
			SetAnalysisName(a.AnalysisName).
			SetAnalysisType(a.AnalysisType).
			SetAnalysisDate(*a.AnalysisDate).
			Save(ctx)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create pet analysis: %w", err)
		}
	}

	if bonuses != nil {
		_, err = tx.PetBonus.Create().
			SetOwner(newPet).
			SetIsArtist(bonuses.IsArtist).
			SetIsTherapist(bonuses.IsTherapist).
			SetIsFormerDonor(bonuses.IsFormerDonor).
			SetIsGuideDog(bonuses.IsGuideDog).
			Save(ctx)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create pet bonus: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return r.GetByID(ctx, newPet.ID, "Health", "Treatments", "Analysis", "Bonuses")
}

// GetByID retrieves a pet by their ID with related entities based on preloads
func (r *EntPetRepository) GetByID(ctx context.Context, id string, preloads ...string) (*ent.Pet, error) {
	if id == "" {
		return nil, errors.New("invalid pet ID")
	}

	query := r.client.Pet.Query().
		Where(pet.ID(id)).
		WithBreedRef().
		WithOwner()

	for _, preload := range preloads {
		switch preload {
		case "Health":
			query = query.WithHealth()
		case "Treatments":
			query = query.WithTreatments()
		case "Analysis":
			query = query.WithAnalyses()
		case "Bonuses":
			query = query.WithBonuses()
		}
	}

	p, err := query.Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("pet with id %s not found: %w", id, err)
		}
		return nil, fmt.Errorf("failed to get pet by id %s: %w", id, err)
	}

	return p, nil
}

// GetByUserID retrieves all pets for a specific user with related entities based on preloads
func (r *EntPetRepository) GetByUserID(ctx context.Context, userID string, preloads ...string) ([]*ent.Pet, error) {
	if userID == "" {
		return nil, errors.New("invalid user ID")
	}

	query := r.client.Pet.Query().
		Where(pet.UserID(userID)).
		WithBreedRef()

	for _, preload := range preloads {
		switch preload {
		case "Health":
			query = query.WithHealth()
		case "Treatments":
			query = query.WithTreatments()
		case "Analyses":
			query = query.WithAnalyses()
		case "Bonuses":
			query = query.WithBonuses()
		}
	}

	pets, err := query.All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get pets for user %s: %w", userID, err)
	}

	return pets, nil
}

// Update updates an existing pet and its related entities in a transaction
func (r *EntPetRepository) Update(ctx context.Context, p *ent.Pet, health *ent.PetHealth, treatments *ent.PetTreatment, analyses []*ent.PetAnalysis, bonuses *ent.PetBonus) (*ent.Pet, error) {
	if p == nil {
		return nil, errors.New("pet cannot be nil")
	}

	if p.ID == "" {
		return nil, errors.New("invalid pet ID")
	}

	tx, err := r.client.Tx(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to start transaction: %w", err)
	}

	// 1. Update Pet
	err = tx.Pet.UpdateOneID(p.ID).
		SetName(p.Name).
		SetType(p.Type).
		SetPetStatus(p.PetStatus).
		SetNillableWeightKg(nillable(p.WeightKg)).
		SetNillableBloodGroup(nillable(p.BloodGroup)).
		SetNillableGender(nillable(p.Gender)).
		SetNillableAgeYears(nillable(p.AgeYears)).
		SetNillableAgeMonths(nillable(p.AgeMonths)).
		SetNillableBirthDate(p.BirthDate).
		SetNillableChipNumber(nillable(p.ChipNumber)).
		SetNillablePhotoURL(nillable(p.PhotoURL)).
		SetNillableBreedID(nillable(p.BreedID)).
		SetNillableUserID(nillable(p.UserID)).
		SetNillableLivingCondition(nillable(p.LivingCondition)).
		Exec(ctx)

	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update pet: %w", err)
	}

	// 2. Update related entities
	if health != nil {
		existingHealth, err := tx.PetHealth.Query().Where(pethealth.HasOwnerWith(pet.ID(p.ID))).Only(ctx)
		if err != nil && !ent.IsNotFound(err) {
			tx.Rollback()
			return nil, err
		}
		if existingHealth != nil {
			err = tx.PetHealth.UpdateOne(existingHealth).
				SetNillableReproductiveStatus(nillable(health.ReproductiveStatus)).
				SetNillableHealthStatus(nillable(health.HealthStatus)).
				SetNillableLastDonation(health.LastDonation).
				SetTransfused(health.Transfused).
				SetMedications(health.Medications).
				SetSurgicalInterventions(health.SurgicalInterventions).
				Exec(ctx)
		} else {
			_, err = tx.PetHealth.Create().
				SetOwnerID(p.ID).
				SetNillableReproductiveStatus(nillable(health.ReproductiveStatus)).
				SetNillableHealthStatus(nillable(health.HealthStatus)).
				SetNillableLastDonation(health.LastDonation).
				SetTransfused(health.Transfused).
				SetMedications(health.Medications).
				SetSurgicalInterventions(health.SurgicalInterventions).
				Save(ctx)
		}
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to update pet health: %w", err)
		}
	}

	if treatments != nil {
		existingTreatment, err := tx.PetTreatment.Query().Where(pettreatment.HasOwnerWith(pet.ID(p.ID))).Only(ctx)
		if err != nil && !ent.IsNotFound(err) {
			tx.Rollback()
			return nil, err
		}
		if existingTreatment != nil {
			err = tx.PetTreatment.UpdateOne(existingTreatment).
				SetNillableRabiesVaccinationDate(treatments.RabiesVaccinationDate).
				SetNillableInfectionVaccinationDate(treatments.InfectionVaccinationDate).
				SetNillableEctoparasiteTreatmentDate(treatments.EctoparasiteTreatmentDate).
				SetNillableDewormingDate(treatments.DewormingDate).
				Exec(ctx)
		} else {
			_, err = tx.PetTreatment.Create().
				SetOwnerID(p.ID).
				SetNillableRabiesVaccinationDate(treatments.RabiesVaccinationDate).
				SetNillableInfectionVaccinationDate(treatments.InfectionVaccinationDate).
				SetNillableEctoparasiteTreatmentDate(treatments.EctoparasiteTreatmentDate).
				SetNillableDewormingDate(treatments.DewormingDate).
				Save(ctx)
		}
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to update pet treatment: %w", err)
		}
	}

	// For analyses, add new ones as history (do not delete existing)
	for _, a := range analyses {
		_, err = tx.PetAnalysis.Create().
			SetOwnerID(p.ID).
			SetAnalysisName(a.AnalysisName).
			SetAnalysisType(a.AnalysisType).
			SetAnalysisDate(*a.AnalysisDate).
			Save(ctx)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create pet analysis: %w", err)
		}
	}

	if bonuses != nil {
		existingBonus, err := tx.PetBonus.Query().Where(petbonus.HasOwnerWith(pet.ID(p.ID))).Only(ctx)
		if err != nil && !ent.IsNotFound(err) {
			tx.Rollback()
			return nil, err
		}
		if existingBonus != nil {
			err = tx.PetBonus.UpdateOne(existingBonus).
				SetIsArtist(bonuses.IsArtist).
				SetIsTherapist(bonuses.IsTherapist).
				SetIsFormerDonor(bonuses.IsFormerDonor).
				SetIsGuideDog(bonuses.IsGuideDog).
				Exec(ctx)
		} else {
			_, err = tx.PetBonus.Create().
				SetOwnerID(p.ID).
				SetIsArtist(bonuses.IsArtist).
				SetIsTherapist(bonuses.IsTherapist).
				SetIsFormerDonor(bonuses.IsFormerDonor).
				SetIsGuideDog(bonuses.IsGuideDog).
				Save(ctx)
		}
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to update pet bonus: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return r.GetByID(ctx, p.ID, "Health", "Treatments", "Analysis", "Bonuses")
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
