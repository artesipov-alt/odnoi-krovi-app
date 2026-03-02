package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
// This entity represents a user in the system, linked to Telegram.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		// telegram_id is the unique identifier from Telegram.
		field.Int64("telegram_id").
			Unique().
			StructTag(`json:"telegramId"`),
		// full_name is the user's full name.
		field.String("full_name").
			Optional().
			MaxLen(255).
			StructTag(`json:"fullName"`),
		// phone is the user's phone number.
		field.String("phone").
			Optional().
			MaxLen(20).
			StructTag(`json:"phone"`),
		// email is the user's email address.
		field.String("email").
			Optional().
			MaxLen(255).
			StructTag(`json:"email"`),
		// organization_name is the name of the user's organization.
		field.String("organization_name").
			Optional().
			MaxLen(255).
			StructTag(`json:"organizationName"`),
		// consent_pd indicates if the user has consented to personal data processing.
		field.Bool("consent_pd").
			Default(false).
			StructTag(`json:"consentPd"`),
		// on_boarding is a list of onboarding steps completed by the user.
		field.JSON("on_boarding", []string{}).
			Optional().
			StructTag(`json:"onBoarding"`),
		// allow_geo indicates if the user allows geolocation.
		field.Bool("allow_geo").
			Default(false).
			StructTag(`json:"allowGeo"`),
		// location_id is the foreign key to the location.
		field.String("location_id").
			Optional().
			StructTag(`json:"locationId"`),
		// photo_urls is a list of URLs to the user's photos.
		field.JSON("photo_urls", []string{}).
			Optional().
			StructTag(`json:"photoUrls"`),
		// role is the user's role in the system.
		field.Enum("role").
			Values("user", "admin").
			Default("user").
			StructTag(`json:"role"`),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		// pets is the edge to the user's pets.
		edge.To("pets", Pet.Type),
		// location is the edge to the user's location.
		edge.From("location", Location.Type).
			Ref("users").
			Unique().
			Field("location_id"),
		// donor_preference is the edge to the user's donor preferences.
		edge.To("donor_preference", DonorPreference.Type).
			Unique(),
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
		entgql.Mutations(
			entgql.MutationCreate(),
			entgql.MutationUpdate(),
		),
	}
}
