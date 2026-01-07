package schema

import (
	"context"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			Immutable().
			DefaultFunc(func() string { return generateID(UserPrefix) }).
			StructTag(`json:"id"`),
		field.Int64("telegram_id").
			Unique().
			StructTag(`json:"telegramId"`),
		field.String("full_name").
			Optional().
			MaxLen(255).
			StructTag(`json:"fullName"`),
		field.String("phone").
			Optional().
			MaxLen(20).
			StructTag(`json:"phone"`),
		field.String("email").
			Optional().
			MaxLen(255).
			StructTag(`json:"email"`),
		field.String("organization_name").
			Optional().
			MaxLen(255).
			StructTag(`json:"organizationName"`),
		field.Bool("consent_pd").
			StructTag(`json:"consentPd"`),
		field.Bool("on_boarding").
			StructTag(`json:"onBoarding"`),
		field.Bool("allow_geo").
			StructTag(`json:"allowGeo"`),
		field.Int("location_id").
			Optional().
			StructTag(`json:"locationId"`),
		field.Enum("role").
			Values("user", "admin").
			Default("user").
			StructTag(`json:"role"`),
		// Audit fields
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			StructTag(`json:"createdAt" swaggerignore:"true"`),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			StructTag(`json:"updatedAt" swaggerignore:"true"`),
		field.Time("deleted_at").
			Optional().
			Nillable().
			StructTag(`json:"deletedAt" swaggerignore:"true"`),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("pets", Pet.Type),
		edge.From("location", Location.Type).
			Ref("users").
			Unique().
			Field("location_id"),
	}
}

// Interceptors of the User.
func (User) Interceptors() []ent.Interceptor {
	return []ent.Interceptor{
		ent.TraverseFunc(func(ctx context.Context, q ent.Query) error {
			// If SkipSoftDelete is in context, show all records (including deleted).
			if skip, _ := ctx.Value(softDeleteKey{}).(bool); skip {
				return nil
			}

			// Add WHERE deleted_at IS NULL filter.
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
