package pg

import (
	"context"
	"errors"
	"fmt"
	"time"

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
func (r *EntPetRepository) Create(ctx context.Context, p *ent.Pet) (*ent.Pet, error) {
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
		SetNillableWeightKg(&p.WeightKg).
		SetNillableBloodGroup(&p.BloodGroup).
		SetNillableGender(&p.Gender).
		SetNillableAgeYears(&p.AgeYears).
		SetNillableAgeMonths(&p.AgeMonths).
		SetNillableBirthDate(p.BirthDate).
		SetNillableChipNumber(&p.ChipNumber).
		SetNillablePhotoURL(&p.PhotoURL).
		SetNillableBreedID(&p.BreedID).
		SetNillableUserID(&p.UserID).
		SetNillableLivingCondition(&p.LivingCondition)

	newPet, err := petCreate.Save(ctx)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create pet: %w", err)
	}

	// 2. Create related entities if provided in Edges
	if p.Edges.Health != nil {
		h := p.Edges.Health
		_, err = tx.PetHealth.Create().
			SetPetID(newPet.ID).
			SetNillableReproductiveStatus(&h.ReproductiveStatus).
			SetNillableHealthStatus(&h.HealthStatus).
			SetNillableLastDonation(h.LastDonation).
			SetTransfused(h.Transfused).
			SetMedications(h.Medications).
			SetSurgicalInterventions(h.SurgicalInterventions).
			Save(ctx)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create pet health: %w", err)
		}
	}

	if p.Edges.Treatments != nil {
		t := p.Edges.Treatments
		_, err = tx.PetTreatment.Create().
			SetPetID(newPet.ID).
			SetNillableRabiesVaccinationDate(t.RabiesVaccinationDate).
			SetNillableInfectionVaccinationDate(t.InfectionVaccinationDate).
			SetNillableEctoparasiteTreatmentDate(t.EctoparasiteTreatmentDate).
			SetNillableDewormingDate(t.DewormingDate).
			Save(ctx)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create pet treatment: %w", err)
		}
	}

	if p.Edges.Analyses != nil {
		a := p.Edges.Analyses
		_, err = tx.PetAnalysis.Create().
			SetPetID(newPet.ID).
			SetNillableLeukemiaDate(a.LeukemiaDate).
			SetNillableLeukemiaType(&a.LeukemiaType).
			SetNillableImmunodeficiencyDate(a.ImmunodeficiencyDate).
			SetNillableImmunodeficiencyType(&a.ImmunodeficiencyType).
			SetNillableHemoplasmosisDate(a.HemoplasmosisDate).
			SetNillableHemoplasmosisType(&a.HemoplasmosisType).
			SetNillableBartonellosisDate(a.BartonellosisDate).
			SetNillableBartonellosisType(&a.BartonellosisType).
			SetNillableBabesiosisDate(a.BabesiosisDate).
			SetNillableBabesiosisType(&a.BabesiosisType).
			SetNillableDirofilariaDate(a.DirofilariaDate).
			SetNillableDirofilariaType(&a.DirofilariaType).
			SetNillableEhrlichiosisDate(a.EhrlichiosisDate).
			SetNillableEhrlichiosisType(&a.EhrlichiosisType).
			SetNillableAnaplasmosisDate(a.AnaplasmosisDate).
			SetNillableAnaplasmosisType(&a.AnaplasmosisType).
			Save(ctx)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create pet analysis: %w", err)
		}
	}

	if p.Edges.Bonuses != nil {
		b := p.Edges.Bonuses
		_, err = tx.PetBonus.Create().
			SetPetID(newPet.ID).
			SetIsArtist(b.IsArtist).
			SetIsTherapist(b.IsTherapist).
			SetIsFormerDonor(b.IsFormerDonor).
			SetIsGuideDog(b.IsGuideDog).
			Save(ctx)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create pet bonus: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return r.GetByID(ctx, newPet.ID)
}

// GetByID retrieves a pet by their ID with all related entities
func (r *EntPetRepository) GetByID(ctx context.Context, id string) (*ent.Pet, error) {
	if id == "" {
		return nil, errors.New("invalid pet ID")
	}

	p, err := r.client.Pet.Query().
		Where(pet.ID(id)).
		WithHealth().
		WithTreatments().
		WithAnalyses().
		WithBonuses().
		WithBreedRef().
		WithOwner().
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("pet with id %s not found: %w", id, err)
		}
		return nil, fmt.Errorf("failed to get pet by id %s: %w", id, err)
	}

	return p, nil
}

// GetByUserID retrieves all pets for a specific user
func (r *EntPetRepository) GetByUserID(ctx context.Context, userID string) ([]*ent.Pet, error) {
	if userID == "" {
		return nil, errors.New("invalid user ID")
	}

	pets, err := r.client.Pet.Query().
		Where(pet.UserID(userID)).
		WithHealth().
		WithTreatments().
		WithAnalyses().
		WithBonuses().
		WithBreedRef().
		All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get pets for user %s: %w", userID, err)
	}

	return pets, nil
}

// Update updates an existing pet and its related entities in a transaction
func (r *EntPetRepository) Update(ctx context.Context, p *ent.Pet) (*ent.Pet, error) {
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
		SetWeightKg(p.WeightKg).
		SetBloodGroup(p.BloodGroup).
		SetGender(p.Gender).
		SetAgeYears(p.AgeYears).
		SetAgeMonths(p.AgeMonths).
		SetNillableBirthDate(p.BirthDate).
		SetChipNumber(p.ChipNumber).
		SetPhotoURL(p.PhotoURL).
		SetBreedID(p.BreedID).
		SetUserID(p.UserID).
		SetLivingCondition(p.LivingCondition).
		Exec(ctx)

	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to update pet: %w", err)
	}

	// 2. Update related entities
	if p.Edges.Health != nil {
		h := p.Edges.Health
		exists, err := tx.PetHealth.Query().Where(pethealth.PetID(p.ID)).Exist(ctx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		if exists {
			err = tx.PetHealth.Update().
				Where(pethealth.PetID(p.ID)).
				SetReproductiveStatus(h.ReproductiveStatus).
				SetHealthStatus(h.HealthStatus).
				SetNillableLastDonation(h.LastDonation).
				SetTransfused(h.Transfused).
				SetMedications(h.Medications).
				SetSurgicalInterventions(h.SurgicalInterventions).
				Exec(ctx)
		} else {
			_, err = tx.PetHealth.Create().
				SetPetID(p.ID).
				SetReproductiveStatus(h.ReproductiveStatus).
				SetHealthStatus(h.HealthStatus).
				SetNillableLastDonation(h.LastDonation).
				SetTransfused(h.Transfused).
				SetMedications(h.Medications).
				SetSurgicalInterventions(h.SurgicalInterventions).
				Save(ctx)
		}
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to update pet health: %w", err)
		}
	}

	if p.Edges.Treatments != nil {
		t := p.Edges.Treatments
		exists, err := tx.PetTreatment.Query().Where(pettreatment.PetID(p.ID)).Exist(ctx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		if exists {
			err = tx.PetTreatment.Update().
				Where(pettreatment.PetID(p.ID)).
				SetNillableRabiesVaccinationDate(t.RabiesVaccinationDate).
				SetNillableInfectionVaccinationDate(t.InfectionVaccinationDate).
				SetNillableEctoparasiteTreatmentDate(t.EctoparasiteTreatmentDate).
				SetNillableDewormingDate(t.DewormingDate).
				Exec(ctx)
		} else {
			_, err = tx.PetTreatment.Create().
				SetPetID(p.ID).
				SetNillableRabiesVaccinationDate(t.RabiesVaccinationDate).
				SetNillableInfectionVaccinationDate(t.InfectionVaccinationDate).
				SetNillableEctoparasiteTreatmentDate(t.EctoparasiteTreatmentDate).
				SetNillableDewormingDate(t.DewormingDate).
				Save(ctx)
		}
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to update pet treatment: %w", err)
		}
	}

	if p.Edges.Analyses != nil {
		a := p.Edges.Analyses
		exists, err := tx.PetAnalysis.Query().Where(petanalysis.PetID(p.ID)).Exist(ctx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		if exists {
			err = tx.PetAnalysis.Update().
				Where(petanalysis.PetID(p.ID)).
				SetNillableLeukemiaDate(a.LeukemiaDate).
				SetLeukemiaType(a.LeukemiaType).
				SetNillableImmunodeficiencyDate(a.ImmunodeficiencyDate).
				SetImmunodeficiencyType(a.ImmunodeficiencyType).
				SetNillableHemoplasmosisDate(a.HemoplasmosisDate).
				SetHemoplasmosisType(a.HemoplasmosisType).
				SetNillableBartonellosisDate(a.BartonellosisDate).
				SetBartonellosisType(a.BartonellosisType).
				SetNillableBabesiosisDate(a.BabesiosisDate).
				SetBabesiosisType(a.BabesiosisType).
				SetNillableDirofilariaDate(a.DirofilariaDate).
				SetDirofilariaType(a.DirofilariaType).
				SetNillableEhrlichiosisDate(a.EhrlichiosisDate).
				SetEhrlichiosisType(a.EhrlichiosisType).
				SetNillableAnaplasmosisDate(a.AnaplasmosisDate).
				SetAnaplasmosisType(a.AnaplasmosisType).
				Exec(ctx)
		} else {
			_, err = tx.PetAnalysis.Create().
				SetPetID(p.ID).
				SetNillableLeukemiaDate(a.LeukemiaDate).
				SetLeukemiaType(a.LeukemiaType).
				SetNillableImmunodeficiencyDate(a.ImmunodeficiencyDate).
				SetImmunodeficiencyType(a.ImmunodeficiencyType).
				SetNillableHemoplasmosisDate(a.HemoplasmosisDate).
				SetHemoplasmosisType(a.HemoplasmosisType).
				SetNillableBartonellosisDate(a.BartonellosisDate).
				SetBartonellosisType(a.BartonellosisType).
				SetNillableBabesiosisDate(a.BabesiosisDate).
				SetBabesiosisType(a.BabesiosisType).
				SetNillableDirofilariaDate(a.DirofilariaDate).
				SetDirofilariaType(a.DirofilariaType).
				SetNillableEhrlichiosisDate(a.EhrlichiosisDate).
				SetEhrlichiosisType(a.EhrlichiosisType).
				SetNillableAnaplasmosisDate(a.AnaplasmosisDate).
				SetAnaplasmosisType(a.AnaplasmosisType).
				Save(ctx)
		}
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to update pet analysis: %w", err)
		}
	}

	if p.Edges.Bonuses != nil {
		b := p.Edges.Bonuses
		exists, err := tx.PetBonus.Query().Where(petbonus.PetID(p.ID)).Exist(ctx)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		if exists {
			err = tx.PetBonus.Update().
				Where(petbonus.PetID(p.ID)).
				SetIsArtist(b.IsArtist).
				SetIsTherapist(b.IsTherapist).
				SetIsFormerDonor(b.IsFormerDonor).
				SetIsGuideDog(b.IsGuideDog).
				Exec(ctx)
		} else {
			_, err = tx.PetBonus.Create().
				SetPetID(p.ID).
				SetIsArtist(b.IsArtist).
				SetIsTherapist(b.IsTherapist).
				SetIsFormerDonor(b.IsFormerDonor).
				SetIsGuideDog(b.IsGuideDog).
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

	return r.GetByID(ctx, p.ID)
}

// Delete deletes a pet by their ID (soft delete)
func (r *EntPetRepository) Delete(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("invalid pet ID")
	}

	err := r.client.Pet.UpdateOneID(id).
		SetDeletedAt(time.Now()).
		Exec(ctx)

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
