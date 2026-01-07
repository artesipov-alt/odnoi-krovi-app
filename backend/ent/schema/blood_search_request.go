package schema

import (
	"entgo.io/ent"
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
		field.String("pet_id").StructTag(`json:"petId"`),
		field.Int32("blood_volume_needed").StructTag(`json:"bloodVolumeNeeded"`),
		field.Int32("blood_volume_reserved").Default(0).StructTag(`json:"bloodVolumeReserved"`),
		field.JSON("regions", []int32{}).StructTag(`json:"regions"`),
		field.Bool("small_pets_notify_allowed").Default(true).StructTag(`json:"smallPetsNotifyAllowed"`),
		field.Enum("status").
			Values("active", "closed", "draft").
			Default("active").
			StructTag(`json:"status"`),
		field.String("description").Optional().StructTag(`json:"description"`),
		field.JSON("photo_urls", []string{}).Optional().StructTag(`json:"photoUrls"`),
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
		edge.To("blood_components", BloodComponent.Type).
			StructTag(`json:"bloodComponents"`),
		edge.From("blood_group", BloodGroup.Type).
			Ref("search_requests").
			Unique(),
	}
}

// Mixins of the BloodSearchRequest.
func (BloodSearchRequest) Mixins() []ent.Mixin {
	return []ent.Mixin{
		NewBaseMixin(BloodSearchPrefix),
	}
}
