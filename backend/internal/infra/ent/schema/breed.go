package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Breed holds the schema definition for the Breed entity.
type Breed struct {
	ent.Schema
}

// Fields of the Breed.
func (Breed) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			Immutable().
			StructTag(`json:"id"`),
		field.String("name").
			MaxLen(100).
			NotEmpty().
			StructTag(`json:"name"`),
		field.Enum("type").
			Values("dog", "cat").
			StructTag(`json:"type"`),
	}
}

// Edges of the Breed.
func (Breed) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("pets", Pet.Type),
	}
}

// Annotations of the Breed.
func (Breed) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "ref_breeds",
		},
	}
}
