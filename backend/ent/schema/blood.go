package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// BloodGroup holds the schema definition for the BloodGroup entity.
type BloodGroup struct {
	ent.Schema
}

// Fields of the BloodGroup.
func (BloodGroup) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("pet_type").
			Values("dog", "cat").
			StructTag(`json:"petType"`),
		field.String("blood_group").
			MaxLen(50).
			NotEmpty().
			StructTag(`json:"bloodGroup"`),
		field.String("description").
			Optional().
			StructTag(`json:"description"`),
	}
}

// Edges of the BloodGroup.
func (BloodGroup) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("search_requests", BloodSearchRequest.Type),
	}
}

// Mixins of the BloodGroup.
func (BloodGroup) Mixins() []ent.Mixin {
	return []ent.Mixin{
		AuditMixin{},
	}
}

// BloodComponent holds the schema definition for the BloodComponent entity.
type BloodComponent struct {
	ent.Schema
}

// Fields of the BloodComponent.
func (BloodComponent) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			MaxLen(255).
			NotEmpty().
			StructTag(`json:"name"`),
	}
}

// Edges of the BloodComponent.
func (BloodComponent) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("search_requests", BloodSearchRequest.Type).
			Ref("blood_components"),
	}
}

// Mixins of the BloodComponent.
func (BloodComponent) Mixins() []ent.Mixin {
	return []ent.Mixin{
		AuditMixin{},
	}
}
