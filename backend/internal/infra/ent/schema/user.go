package schema

import (
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
			Optional().
			Unique(),
		// full_name is the user's full name.
		field.String("full_name").
			Optional().
			MaxLen(255),
		// phone is the user's phone number.
		field.String("phone").
			Optional().
			MaxLen(20),
		// email is the user's email address.
		field.String("email").
			Optional().
			MaxLen(255),
		// organization_name is the name of the user's organization.
		field.String("organization_name").
			Optional().
			MaxLen(255),
		// consent_pd indicates if the user has consented to personal data processing.
		field.Bool("consent_pd").
			Default(false),
		// on_boarding is a list of onboarding steps completed by the user.
		field.JSON("on_boarding", []string{}).
			Optional(),
		// allow_geo indicates if the user allows geolocation.
		field.Bool("allow_geo").
			Default(false),
		// location_id is the foreign key to the location.
		field.String("location_id").
			Optional(),
		// photo_urls is a list of URLs to the user's photos.
		field.JSON("photo_urls", []string{}).
			Optional(),
		// role is the user's role in the system.
		field.Enum("role").
			Values("user", "admin").
			Default("user"),
		// origin_source is the origin source string.
		field.String("origin_source").
			Optional().
			MaxLen(255),
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
		// identities is the edge to the user's identities.
		edge.To("identities", UserIdentity.Type),
		// utm_histories is the edge to the user's UTM history.
		edge.To("utm_histories", UtmHistory.Type),
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
