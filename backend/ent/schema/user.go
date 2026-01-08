package schema

import (
	"entgo.io/ent"
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

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		AuditMixin{},
	}
}
