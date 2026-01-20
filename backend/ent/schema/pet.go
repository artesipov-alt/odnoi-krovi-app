package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Pet holds the schema definition for the Pet entity.
type Pet struct {
	ent.Schema
}

// Mixin of the Pet.
func (Pet) Mixin() []ent.Mixin {
	return []ent.Mixin{
		StandardMixin{Prefix: PetPrefix},
	}
}

// Fields of the Pet.
func (Pet) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			MaxLen(100).StructTag(`json:"name"`),
		field.Enum("type").
			Values("dog", "cat").StructTag(`json:"type"`),
		field.Enum("pet_status").
			Values("donor", "recipient", "none").StructTag(`json:"petStatus"`),
		field.Float("weight_kg").
			Optional().StructTag(`json:"weightKg"`),
		field.String("blood_group").
			Optional().StructTag(`json:"bloodGroup"`),
		field.Enum("gender").
			Values("male", "female").Optional().StructTag(`json:"gender"`),
		field.Int("age_years").
			Optional().StructTag(`json:"ageYears"`),
		field.Int("age_months").
			Optional().StructTag(`json:"ageMonths"`),
		field.Time("birth_date").
			Optional().Nillable().StructTag(`json:"birthDate"`),
		field.String("chip_number").
			Optional().MaxLen(15).StructTag(`json:"chipNumber"`),
		field.JSON("photo_urls", []string{}).
			Optional().StructTag(`json:"photoUrls"`),
		field.Int("breed_id").
			Optional().StructTag(`json:"breedId"`),
		field.String("user_id").
			Optional().StructTag(`json:"userId"`),
		field.String("health_id").
			Optional().StructTag(`json:"healthId"`),
		field.String("treatment_id").
			Optional().StructTag(`json:"treatmentId"`),
		field.String("bonus_id").
			Optional().StructTag(`json:"bonusId"`),
		field.Enum("living_condition").
			Values("indoor", "leash_walking", "self_outdoor").Optional().StructTag(`json:"livingCondition"`),
	}
}

// Edges of the Pet.
func (Pet) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).
			Ref("pets").
			Unique().
			Field("user_id"),
		edge.From("health", PetHealth.Type).Ref("owner").Unique().Field("health_id"),
		edge.From("treatments", PetTreatment.Type).Ref("owner").Unique().Field("treatment_id"),
		edge.From("analyses", PetAnalysis.Type).Ref("owner"),
		edge.From("bonuses", PetBonus.Type).Ref("owner").Unique().Field("bonus_id"),
		edge.From("breed_ref", Breed.Type).
			Ref("pets").
			Unique().
			Field("breed_id"),
		edge.To("blood_search_request", BloodSearchRequest.Type).Unique(),
	}
}

// Annotations of the Pet.
func (Pet) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "pets",
		},
	}
}

// PetHealth holds the schema definition for the PetHealth entity.
type PetHealth struct {
	ent.Schema
}

// Mixin of the PetHealth.
func (PetHealth) Mixin() []ent.Mixin {
	return []ent.Mixin{
		StandardMixin{Prefix: PetHealthPrefix},
	}
}

// Fields of the PetHealth.
func (PetHealth) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("reproductive_status").Values("pregnancy", "lactation", "estrus").Optional().StructTag(`json:"reproductiveStatus"`),
		field.Enum("health_status").Values("healthy", "ill", "unknown").Optional().StructTag(`json:"healthStatus"`),
		field.Time("last_donation").Optional().Nillable().StructTag(`json:"lastDonation"`),
		field.Bool("transfused").Optional().StructTag(`json:"transfused"`),
		field.String("medications").Optional().StructTag(`json:"medications"`),
		field.String("surgical_interventions").Optional().StructTag(`json:"surgicalInterventions"`),
	}
}

// Edges of the PetHealth.
func (PetHealth) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("owner", Pet.Type).
			Unique().
			Required(),
	}
}

// Annotations of the PetHealth.
func (PetHealth) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "pet_healths",
		},
	}
}

// PetTreatment holds the schema definition for the PetTreatment entity.
type PetTreatment struct {
	ent.Schema
}

// Mixin of the PetTreatment.
func (PetTreatment) Mixin() []ent.Mixin {
	return []ent.Mixin{
		StandardMixin{Prefix: PetTreatmentPrefix},
	}
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

// Edges of the PetTreatment.
func (PetTreatment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("owner", Pet.Type).
			Unique().
			Required(),
	}
}

// Annotations of the PetTreatment.
func (PetTreatment) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "pet_treatments",
		},
	}
}

// PetAnalysis holds the schema definition for the PetAnalysis entity.
type PetAnalysis struct {
	ent.Schema
}

// Mixin of the PetAnalysis.
func (PetAnalysis) Mixin() []ent.Mixin {
	return []ent.Mixin{
		StandardMixin{Prefix: PetAnalysisPrefix},
	}
}

// Fields of the PetAnalysis.
func (PetAnalysis) Fields() []ent.Field {
	return []ent.Field{
		field.String("pet_id").
			StructTag(`json:"petId"`),
		field.Enum("analysis_name").Values("leukemia", "immunodeficiency", "hemoplasmosis", "bartonellosis", "babesiosis", "dirofilaria", "ehrlichiosis", "anaplasmosis").
			Optional().
			StructTag(`json:"analysisName"`),
		field.Enum("analysis_type").Values("PCR", "ELISA", "ICA", "Microscopy", "Express").
			Optional().
			StructTag(`json:"analysisType"`),
		field.Time("analysis_date").Nillable().StructTag(`json:"analysisDate"`),
	}
}

// Edges of the PetAnalysis.
func (PetAnalysis) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("owner", Pet.Type).
			Unique().
			Required().
			Field("pet_id"),
	}
}

// Annotations of the PetAnalysis.
func (PetAnalysis) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "pet_analyses",
		},
	}
}

// PetBonus holds the schema definition for the PetBonus entity.
type PetBonus struct {
	ent.Schema
}

// Mixin of the PetBonus.
func (PetBonus) Mixin() []ent.Mixin {
	return []ent.Mixin{
		StandardMixin{Prefix: PetBonusPrefix},
	}
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

// Edges of the PetBonus.
func (PetBonus) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("owner", Pet.Type).
			Unique().
			Required(),
	}
}

// Annotations of the PetBonus.
func (PetBonus) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "pet_bonuses",
		},
	}
}
