package schema

import (
	"context"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
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
		// Audit fields
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			StructTag(`json:"createdAt"`),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			StructTag(`json:"updatedAt"`),
		field.Time("deleted_at").
			Optional().
			Nillable().
			StructTag(`json:"deletedAt"`),
	}
}

// Edges of the Breed.
func (Breed) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("pets", Pet.Type),
	}
}

// Interceptors of the Breed.
func (Breed) Interceptors() []ent.Interceptor {
	return []ent.Interceptor{
		ent.TraverseFunc(func(ctx context.Context, q ent.Query) error {
			if skip, _ := ctx.Value(softDeleteKey{}).(bool); skip {
				return nil
			}
			type query interface {
				WhereP(...func(*sql.Selector))
			}
			if w, ok := q.(query); ok {
				w.WhereP(func(s *sql.Selector) {
					s.Where(sql.IsNull(s.C("deleted_at")))
				})
			}
			return nil
		}),
	}
}
