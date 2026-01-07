package schema

import (
	"entgo.io/ent"
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
		field.Int("id").
			Unique().
			Immutable(),
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
