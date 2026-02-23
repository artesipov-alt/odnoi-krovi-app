package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Примечание: структурные теги убраны для упрощения схемы.
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
		field.String("name"),
		field.String("type"),
		field.Float("weight_kg").Optional(),
		field.String("gender").Optional(),
		field.Time("birth_date").Optional().Nillable(),
		field.String("chip_number").Optional().MaxLen(15),
		field.JSON("photo_urls", []string{}).Optional(),
		field.String("breed_id").Optional().Nillable(),
		field.String("user_id").Optional(),
		field.String("health_id").Optional(),
		field.String("treatment_id").Optional(),
		field.String("living_condition").Optional(),
		field.String("reproductive_status").Optional(),
		field.JSON("donor_restrictions", []string{}).Optional(),
		field.String("blood_group_id").Optional().Nillable(),
		field.JSON("bonuses", []string{}).Optional(),
	}
}

// Edges of the Pet.
func (Pet) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).
			Ref("pets").
			Unique().
			Field("user_id"),
		edge.From("health", PetHealth.Type).
			Ref("owner").
			Unique().
			Field("health_id"),
		edge.From("treatments", PetTreatment.Type).
			Ref("owner").
			Unique().
			Field("treatment_id"),
		edge.From("analyses", PetAnalysis.Type).
			Ref("owner"),
		edge.From("breed_ref", Breed.Type).
			Ref("pets").
			Unique().
			Field("breed_id"),
		edge.From("blood_group_ref", BloodGroup.Type).
			Ref("pets").
			Unique().
			Field("blood_group_id"),
		edge.From("donations", DonorResponse.Type).
			Ref("donor"),
		edge.To("blood_search_request", BloodSearchRequest.Type).
			Unique(),
	}
}

// Annotations of the Pet.
func (Pet) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "pets",
		},
		entgql.Mutations(
			entgql.MutationCreate(),
			entgql.MutationUpdate(),
		),
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
		field.Enum("health_status").Values("healthy", "ill", "unknown").Optional(),
		field.Time("last_donation").Optional().Nillable(),
		field.Bool("transfused").Optional(),
		field.String("medications").Optional(),
		field.String("surgical_interventions").Optional(),
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
		entgql.Mutations(
			entgql.MutationCreate(),
			entgql.MutationUpdate(),
		),
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
		field.Time("rabies_vaccination_date").Optional().Nillable(),
		field.Time("infection_vaccination_date").Optional().Nillable(),
		field.Time("ectoparasite_treatment_date").Optional().Nillable(),
		field.Time("deworming_date").Optional().Nillable(),
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
		entgql.Mutations(
			entgql.MutationCreate(),
			entgql.MutationUpdate(),
		),
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
		field.String("pet_id"),
		field.Enum("analysis_name").Values("leukemia", "immunodeficiency", "hemoplasmosis", "bartonellosis", "babesiosis", "dirofilaria", "ehrlichiosis", "anaplasmosis").
			Optional(),
		field.Enum("analysis_type").Values("PCR", "ELISA", "ICA", "Microscopy", "Express").
			Optional(),
		field.Time("analysis_date").Nillable(),
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
		entgql.Mutations(
			entgql.MutationCreate(),
			entgql.MutationUpdate(),
		),
	}
}
