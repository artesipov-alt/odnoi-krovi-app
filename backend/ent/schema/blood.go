package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// BloodGroupPrefix is the prefix for BloodGroup IDs
const BloodGroupPrefix = "BLG"

// BloodComponentPrefix is the prefix for BloodComponent IDs
const BloodComponentPrefix = "BLC"

// BloodGroup holds the schema definition for the BloodGroup entity.
type BloodGroup struct {
	ent.Schema
}

// Fields of the BloodGroup.
func (BloodGroup) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			Immutable().
			StructTag(`json:"id"`),
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
		edge.To("pets", Pet.Type),
	}
}

// Annotations of the BloodGroup.
func (BloodGroup) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "ref_bloodg",
		},
	}
}

// BloodComponent holds the schema definition for the BloodComponent entity.
type BloodComponent struct {
	ent.Schema
}

// Fields of the BloodComponent.
func (BloodComponent) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			Immutable().
			StructTag(`json:"id"`),
		field.String("name").
			MaxLen(255).
			NotEmpty().
			StructTag(`json:"name"`),
	}
}

// Edges of the BloodComponent.
func (BloodComponent) Edges() []ent.Edge {
	return []ent.Edge{}
}

// Annotations of the BloodComponent.
func (BloodComponent) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "ref_bloodc",
		},
	}
}
