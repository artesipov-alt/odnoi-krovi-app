package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Location holds the schema definition for the Location entity.
type Location struct {
	ent.Schema
}

// Fields of the Location.
func (Location) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").StructTag(`json:"id"`),
		field.String("name").
			MaxLen(255).
			NotEmpty().
			StructTag(`json:"name"`),
	}
}

// Edges of the Location.
func (Location) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("users", User.Type),
	}
}

// Annotations of the Location.
func (Location) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Schema("reference"),
		entsql.Table("locations"),
	}
}
