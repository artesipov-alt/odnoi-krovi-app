package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// UserIdentity holds the schema definition for the UserIdentity entity.
// This entity represents a user's identity from an external provider (Telegram, Google, etc.).
type UserIdentity struct {
	ent.Schema
}

// Fields of the UserIdentity.
func (UserIdentity) Fields() []ent.Field {
	return []ent.Field{
		// user_id - внешний ключ к пользователю.
		field.String("user_id"),
		// partner_id - внешний ключ к партнеру.
		field.String("partner_id").
			Optional(),
		// provider - поставщик идентификационных данных (telegram_bot, telegram_webapp, google и т.д.).
		field.Enum("provider").
			Values("telegram_bot", "max_bot", "service"),
		// provider_user_id - уникальный идентификатор от поставщика (например, Telegram ID).
		field.String("provider_user_id"),
		// metadata - JSONB для хранения дополнительных данных поставщика (имя пользователя, photo_url и т.д.).
		field.JSON("metadata", map[string]any{}).
			Optional(),
	}
}

// Edges of the UserIdentity.
func (UserIdentity) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("identities").
			Field("user_id").
			Required().
			Unique(),
		edge.From("partner", Partner.Type).
			Ref("partner_identities").
			Field("partner_id").
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

func (UserIdentity) Mixin() []ent.Mixin {
	return []ent.Mixin{
		StandardMixin{Prefix: IdentityPrefix},
	}
}

func (UserIdentity) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("provider_user_id", "provider").
			Unique(),
	}
}
