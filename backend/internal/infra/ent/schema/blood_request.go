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
		field.String("pet_id"),
		field.Int32("blood_volume_needed"),
		field.JSON("regions", []string{}).
			Annotations(entgql.Type("Int")),
		field.Bool("small_pets_notify_allowed").
			Default(true),
		field.Enum("status").
			Values("active", "closed", "draft", "reserved_full").
			Default("active"),
		field.String("description").
			Optional(),
		field.JSON("photo_urls", []string{}).
			Annotations(entgql.Type("String")).
			Optional(),
		field.JSON("blood_group_names", []string{}).
			Annotations(entgql.Type("String")).
			Optional(),
		field.JSON("blood_component_ids", []string{}).
			Annotations(entgql.Type("String")).
			Optional(),
		field.JSON("on_boarding", []string{}).
			Optional(),
		field.Bool("priority_search").
			Default(false),
		field.Bool("include_unknown_blood_group").
			Default(false),
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
		edge.To("responses", DonorResponse.Type),
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
