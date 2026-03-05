package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// UtmHistory holds the schema definition for the UtmHistory entity.
type UtmHistory struct {
	ent.Schema
}

// Fields of the UtmHistory.
func (UtmHistory) Fields() []ent.Field {
	return []ent.Field{
		field.String("user_id"),
		field.String("utm_source").
			Optional(),
		field.String("utm_medium").
			Optional(),
		field.String("utm_campaign").
			Optional(),
		field.String("utm_content").
			Optional(),
		field.String("utm_term").
			Optional(),
	}
}

// Edges of the UtmHistory.
func (UtmHistory) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("utm_histories").
			Unique().
			Field("user_id").
			Required(),
	}
}

func (UtmHistory) Mixin() []ent.Mixin {
	return []ent.Mixin{
		StandardMixin{Prefix: UtmHistoryPrefix},
	}
}

// Annotations of the UtmHistory.
func (UtmHistory) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "user_utm_history",
		},
	}
}

// Indexes of the UtmHistory.
func (UtmHistory) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("user_id", "utm_campaign", "utm_source").
			Unique(),
	}
}
