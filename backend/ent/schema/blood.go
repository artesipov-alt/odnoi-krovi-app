package schema

import (
	"context"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
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

// Edges of the BloodGroup.
func (BloodGroup) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("search_requests", BloodSearchRequest.Type),
	}
}

// Interceptors of the BloodGroup.
func (BloodGroup) Interceptors() []ent.Interceptor {
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

// Edges of the BloodComponent.
func (BloodComponent) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("search_requests", BloodSearchRequest.Type).
			Ref("blood_components"),
	}
}

// Interceptors of the BloodComponent.
func (BloodComponent) Interceptors() []ent.Interceptor {
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
