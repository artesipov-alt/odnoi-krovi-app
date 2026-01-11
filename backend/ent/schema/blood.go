package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// BloodGroup holds the schema definition for the BloodGroup entity.
type BloodGroup struct {
	ent.Schema
}

// Fields of the BloodGroup.
func (BloodGroup) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").StructTag(`json:"id"`),
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
	return []ent.Edge{}
}

// Annotations of the BloodGroup.
func (BloodGroup) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table:  "blood_groups",
			Schema: "reference",
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
		field.Int("id").
			Unique().
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
		entsql.Schema("reference"),
		entsql.Table("blood_components"),
	}
}
