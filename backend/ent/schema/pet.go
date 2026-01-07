package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Pet holds the schema definition for the Pet entity.
type Pet struct {
	ent.Schema
}

// Fields of the Pet.
func (Pet) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			Immutable().
			DefaultFunc(func() string { return generateID(PetPrefix) }).
			StructTag(`json:"id"`),
		field.String("name").MaxLen(100).StructTag(`json:"name"`),
		field.Enum("type").Values("dog", "cat").StructTag(`json:"type"`),
		field.Enum("pet_status").Values("donor", "recipient").StructTag(`json:"petStatus"`),
		field.Float("weight_kg").Optional().StructTag(`json:"weightKg"`),
		field.String("blood_group").Optional().StructTag(`json:"bloodGroup"`),
		field.Enum("gender").Values("male", "female").Optional().StructTag(`json:"gender"`),
		field.Int("age_years").Optional().StructTag(`json:"ageYears"`),
		field.Int("age_months").Optional().StructTag(`json:"ageMonths"`),
		field.Time("birth_date").Optional().Nillable().StructTag(`json:"birthDate"`),
		field.String("chip_number").Optional().MaxLen(15).StructTag(`json:"chipNumber"`),
		field.String("photo_url").Optional().MaxLen(255).StructTag(`json:"photoUrl"`),
		field.Int("breed_id").Optional().StructTag(`json:"breedId"`),
		field.String("user_id").Optional().StructTag(`json:"userId"`),
		field.Enum("living_condition").Values("indoor", "leash_walking", "self_outdoor").Optional().StructTag(`json:"livingCondition"`),
	}
}

// Edges of the Pet.
func (Pet) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).
			Ref("pets").
			Unique().
			Field("user_id"),
		edge.To("health", PetHealth.Type).Unique(),
		edge.To("treatments", PetTreatment.Type).Unique(),
		edge.To("analyses", PetAnalysis.Type).Unique(),
		edge.To("bonuses", PetBonus.Type).Unique(),
		edge.From("breed_ref", Breed.Type).
			Ref("pets").
			Unique().
			Field("breed_id"),
		edge.To("blood_search_request", BloodSearchRequest.Type).Unique(),
	}
}

// Mixins of the Pet.
func (Pet) Mixins() []ent.Mixin {
	return []ent.Mixin{
		AuditMixin{},
	}
}

// PetHealth holds the schema definition for the PetHealth entity.
type PetHealth struct {
	ent.Schema
}

// Fields of the PetHealth.
func (PetHealth) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("reproductive_status").Values("pregnancy", "lactation", "estrus", "none").Optional().StructTag(`json:"reproductiveStatus"`),
		field.Enum("health_status").Values("healthy", "ill", "unknown").Optional().StructTag(`json:"healthStatus"`),
		field.Time("last_donation").Optional().Nillable().StructTag(`json:"lastDonation"`),
		field.Bool("transfused").Optional().StructTag(`json:"transfused"`),
		field.String("medications").Optional().StructTag(`json:"medications"`),
		field.String("surgical_interventions").Optional().StructTag(`json:"surgicalInterventions"`),
		field.String("pet_id").Optional().StructTag(`json:"petId"`),
	}
}

// Mixins of the PetHealth.
func (PetHealth) Mixins() []ent.Mixin {
	return []ent.Mixin{
		AuditMixin{},
	}
}

// Edges of the PetHealth.
func (PetHealth) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("pet", Pet.Type).
			Ref("health").
			Unique().
			Field("pet_id"),
	}
}

// PetTreatment holds the schema definition for the PetTreatment entity.
type PetTreatment struct {
	ent.Schema
}

// Fields of the PetTreatment.
func (PetTreatment) Fields() []ent.Field {
	return []ent.Field{
		field.Time("rabies_vaccination_date").Optional().Nillable().StructTag(`json:"rabiesVaccinationDate"`),
		field.Time("infection_vaccination_date").Optional().Nillable().StructTag(`json:"infectionVaccinationDate"`),
		field.Time("ectoparasite_treatment_date").Optional().Nillable().StructTag(`json:"ectoparasiteTreatmentDate"`),
		field.Time("deworming_date").Optional().Nillable().StructTag(`json:"dewormingDate"`),
		field.String("pet_id").Optional().StructTag(`json:"petId"`),
	}
}

// Mixins of the PetTreatment.
func (PetTreatment) Mixins() []ent.Mixin {
	return []ent.Mixin{
		AuditMixin{},
	}
}

// Edges of the PetTreatment.
func (PetTreatment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("pet", Pet.Type).
			Ref("treatments").
			Unique().
			Field("pet_id"),
	}
}

// PetAnalysis holds the schema definition for the PetAnalysis entity.
type PetAnalysis struct {
	ent.Schema
}

// Fields of the PetAnalysis.
func (PetAnalysis) Fields() []ent.Field {
	return []ent.Field{
		field.Time("leukemia_date").Optional().Nillable().StructTag(`json:"leukemiaDate"`),
		field.Enum("leukemia_type").Values("PCR", "ELISA", "ICA", "Microscopy", "Express").Optional().StructTag(`json:"leukemiaType"`),
		field.Time("immunodeficiency_date").Optional().Nillable().StructTag(`json:"immunodeficiencyDate"`),
		field.Enum("immunodeficiency_type").Values("PCR", "ELISA", "ICA", "Microscopy", "Express").Optional().StructTag(`json:"immunodeficiencyType"`),
		field.Time("hemoplasmosis_date").Optional().Nillable().StructTag(`json:"hemoplasmosisDate"`),
		field.Enum("hemoplasmosis_type").Values("PCR", "ELISA", "ICA", "Microscopy", "Express").Optional().StructTag(`json:"hemoplasmosisType"`),
		field.Time("bartonellosis_date").Optional().Nillable().StructTag(`json:"bartonellosisDate"`),
		field.Enum("bartonellosis_type").Values("PCR", "ELISA", "ICA", "Microscopy", "Express").Optional().StructTag(`json:"bartonellosisType"`),
		field.Time("babesiosis_date").Optional().Nillable().StructTag(`json:"babesiosisDate"`),
		field.Enum("babesiosis_type").Values("PCR", "ELISA", "ICA", "Microscopy", "Express").Optional().StructTag(`json:"babesiosisType"`),
		field.Time("dirofilaria_date").Optional().Nillable().StructTag(`json:"dirofilariaDate"`),
		field.Enum("dirofilaria_type").Values("PCR", "ELISA", "ICA", "Microscopy", "Express").Optional().StructTag(`json:"dirofilariaType"`),
		field.Time("ehrlichiosis_date").Optional().Nillable().StructTag(`json:"ehrlichiosisDate"`),
		field.Enum("ehrlichiosis_type").Values("PCR", "ELISA", "ICA", "Microscopy", "Express").Optional().StructTag(`json:"ehrlichiosisType"`),
		field.Time("anaplasmosis_date").Optional().Nillable().StructTag(`json:"anaplasmosisDate"`),
		field.Enum("anaplasmosis_type").Values("PCR", "ELISA", "ICA", "Microscopy", "Express").Optional().StructTag(`json:"anaplasmosisType"`),
		field.String("pet_id").Optional().StructTag(`json:"petId"`),
	}
}

// Mixins of the PetAnalysis.
func (PetAnalysis) Mixins() []ent.Mixin {
	return []ent.Mixin{
		AuditMixin{},
	}
}

// Edges of the PetAnalysis.
func (PetAnalysis) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("pet", Pet.Type).
			Ref("analyses").
			Unique().
			Field("pet_id"),
	}
}

// PetBonus holds the schema definition for the PetBonus entity.
type PetBonus struct {
	ent.Schema
}

// Fields of the PetBonus.
func (PetBonus) Fields() []ent.Field {
	return []ent.Field{
		field.Bool("is_artist").StructTag(`json:"isArtist"`),
		field.Bool("is_therapist").StructTag(`json:"isTherapist"`),
		field.Bool("is_former_donor").StructTag(`json:"isFormerDonor"`),
		field.Bool("is_guide_dog").StructTag(`json:"isGuideDog"`),
		field.String("pet_id").Optional().StructTag(`json:"petId"`),
	}
}

// Mixins of the PetBonus.
func (PetBonus) Mixins() []ent.Mixin {
	return []ent.Mixin{
		AuditMixin{},
	}
}

// Edges of the PetBonus.
func (PetBonus) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("pet", Pet.Type).
			Ref("bonuses").
			Unique().
			Field("pet_id"),
	}
}
