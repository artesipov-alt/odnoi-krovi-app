package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// UserIdentity holds the schema definition for the UserIdentity entity.
// This entity represents a user's identity from an external provider (Telegram, Google, etc.).
type UserIdentity struct {
	ent.Schema
}

// Fields of the UserIdentity.
func (UserIdentity) Fields() []ent.Field {
	return []ent.Field{
		// user_id is the foreign key to the user.
		field.Int("user_id"),
		// provider is the identity provider (telegram_bot, telegram_webapp, google, etc.).
		field.Enum("provider").
			Values("telegram_bot", "max_bot"),
		// provider_user_id is the unique identifier from the provider (e.g., Telegram ID).
		field.String("provider_user_id"),
		// metadata is JSONB for storing additional provider data (username, photo_url, etc.).
		field.JSON("metadata", map[string]any{}).
			Optional(),
	}
}

// Edges of the UserIdentity.
func (UserIdentity) Edges() []ent.Edge {
	return []ent.Edge{
		// user is the edge to the user.
		edge.From("user", User.Type).
			Ref("identities").
			Field("user_id").
			Required().
			Unique(),
	}
}

// Annotations of the UserIdentity.
func (UserIdentity) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "user_identities",
		},
	}
}
