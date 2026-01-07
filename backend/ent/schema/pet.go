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
		field.String("breed").Optional().MaxLen(100).StructTag(`json:"breed"`),
		field.Enum("living_condition").Values("indoor", "leash_walking", "self_outdoor").Optional().StructTag(`json:"livingCondition"`),
	}
}

// Edges of the Pet.
func (Pet) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).Ref("pets").Unique(),
		edge.To("health", PetHealth.Type),
		edge.To("treatments", PetTreatment.Type),
		edge.To("analyses", PetAnalysis.Type),
		edge.To("bonuses", PetBonus.Type),
	}
}

// Mixins of the Pet.
func (Pet) Mixins() []ent.Mixin {
	return []ent.Mixin{
		NewBaseMixin(PetPrefix),
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
		edge.From("pet", Pet.Type).Ref("health").Unique(),
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
		edge.From("pet", Pet.Type).Ref("treatments").Unique(),
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
		edge.From("pet", Pet.Type).Ref("analyses"),
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
		edge.From("pet", Pet.Type).Ref("bonuses").Unique(),
	}
}
