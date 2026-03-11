package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Partner holds the schema definition for the Partner entity.
// This entity represents a partner in the system, such as a clinic or service.
type Partner struct {
	ent.Schema
}

// Fields of the Partner.
func (Partner) Fields() []ent.Field {
	return []ent.Field{
		// name is a human-readable name (e.g., "Клиника на Пушкина" or "Бот Макс").
		field.String("name"),
		// provider_name is the name of the provider associated with this partner.
		// api_key is the secret key, stored as plain text for now with a unique index.
		field.String("api_key").
			Unique().
			Sensitive(),
		// role defines the role of the key holder (CLINIC, ADMIN, SERVICE), to be included in JWT.
		field.Enum("role").
			Values("clinic", "admin", "service"),
		// status allows instant banning of the partner (active, disabled, expired).
		field.Enum("status").
			Values("active", "disabled", "expired"),
		// description is for notes on who and why the key was issued.
		field.String("description").
			Optional(),
		// last_used_at tracks when the partner was last active.
		field.Time("last_used_at").
			Optional(),
	}
}

// Edges of the Partner.
func (Partner) Edges() []ent.Edge {
	return []ent.Edge{
		// partner_identities are the edges to the user identities associated with this partner.
		edge.To("partner_identities", UserIdentity.Type),
	}
}

// Annotations of the Partner.
func (Partner) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "partners",
		},
	}
}

func (Partner) Mixin() []ent.Mixin {
	return []ent.Mixin{
		StandardMixin{Prefix: PartnerPrefix},
	}
}
