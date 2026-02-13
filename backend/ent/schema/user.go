package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
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
		field.JSON("on_boarding", []string{}).
			Optional().
			StructTag(`json:"onBoarding"`),
		field.Bool("allow_geo").
			StructTag(`json:"allowGeo"`),
		field.Int("location_id").
			Optional().
			StructTag(`json:"locationId"`),
		field.JSON("photo_urls", []string{}).
			Optional().
			StructTag(`json:"photoUrls"`),
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
		StandardMixin{Prefix: UserPrefix},
	}
}

// Annotations of the User.
func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "users",
		},
	}
}
