package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// BloodSearchRequest holds the schema definition for the BloodSearchRequest entity.
type BloodSearchRequest struct {
	ent.Schema
}

// Fields of the BloodSearchRequest.
func (BloodSearchRequest) Fields() []ent.Field {
	return []ent.Field{
		field.String("pet_id").
			StructTag(`json:"petId"`),
		field.Int32("blood_volume_needed").
			StructTag(`json:"bloodVolumeNeeded"`),
		field.Int32("blood_volume_reserved").
			Default(0).
			StructTag(`json:"bloodVolumeReserved"`),
		field.JSON("regions", []int32{}).
			Annotations(entgql.Type("Int")).
			StructTag(`json:"regions"`),
		field.Bool("small_pets_notify_allowed").
			Default(true).
			StructTag(`json:"smallPetsNotifyAllowed"`),
		field.Enum("status").
			Values("active", "closed", "draft").
			Default("active").
			StructTag(`json:"status"`),
		field.String("description").
			Optional().
			StructTag(`json:"description"`),
		field.JSON("photo_urls", []string{}).
			Annotations(entgql.Type("String")).
			Optional().
			StructTag(`json:"photoUrls"`),
		field.JSON("blood_group_names", []string{}).
			Annotations(entgql.Type("String")).
			Optional().
			StructTag(`json:"bloodGroupNames"`),
		field.JSON("blood_component_ids", []int{}).
			Annotations(entgql.Type("Int")).
			Optional().
			StructTag(`json:"bloodComponentIds"`),
	}
}

// Edges of the BloodSearchRequest.
func (BloodSearchRequest) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("pet", Pet.Type).
			Ref("blood_search_request").
			Field("pet_id").
			Unique().
			Required(),
	}
}

func (BloodSearchRequest) Mixin() []ent.Mixin {
	return []ent.Mixin{
		StandardMixin{Prefix: BloodSearchPrefix},
	}
}

// Annotations of the BloodSearchRequest.
func (BloodSearchRequest) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "blood_requests",
		},
	}

}
