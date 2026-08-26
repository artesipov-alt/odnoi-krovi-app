package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// User определяет схему для сущности User.
// Эта сущность представляет пользователя в системе, связанного с Telegram.
type User struct {
	ent.Schema
}

// Поля User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		// full_name - полное имя пользователя.
		field.String("full_name").
			Optional().
			MaxLen(255),
		// phone - номер телефона пользователя.
		field.String("phone").
			Optional().
			MaxLen(20),
		// verified - указывает, верифицирован ли пользователь.
		field.Bool("verified").
			Default(false),
		// email - адрес электронной почты пользователя.
		field.String("email").
			Optional().
			MaxLen(255),
		// organization_name - название организации пользователя.
		field.String("organization_name").
			Optional().
			MaxLen(255),
		// consent_pd указывает, согласился ли пользователь на обработку персональных данных.
		field.Bool("consent_pd").
			Default(false),
		// on_boarding - список шагов онбординга, выполненных пользователем.
		field.JSON("on_boarding", []string{}).
			Optional(),
		// allow_geo указывает, разрешает ли пользователь геолокацию.
		field.Bool("allow_geo").
			Default(false),
		// location_id - внешний ключ к локации.
		field.String("location_id").
			Optional(),
		// photo_urls - список URL-адресов фотографий пользователя.
		field.JSON("photo_urls", []string{}).
			Optional(),
		// role - роль пользователя в системе.
		field.Enum("role").
			Values("user", "admin", "clinic").
			Default("user"),
		// origin_source - UTM-метка, указывающая, откуда пришел пользователь.
		field.String("origin_source").
			Optional().
			MaxLen(255),
		// priority_search_count - количество приоритетных поисков, начисляемых за донации.
		field.Int("priority_search_count").
			Default(0),
		// last_seen_at - время последнего посещения пользователя.
		field.Time("last_seen_at").
			Optional().
			Nillable(),
	}
}

// Связи User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		// pets - связь с питомцами пользователя.
		edge.To("pets", Pet.Type),
		// location - связь с локацией пользователя.
		edge.From("location", Location.Type).
			Ref("users").
			Unique().
			Field("location_id"),
		// donor_preference - связь с предпочтениями донора пользователя.
		edge.To("donor_preference", DonorPreference.Type).
			Unique(),
		// identities - связь с идентичностями пользователя.
		edge.To("identities", UserIdentity.Type),
		// utm_histories - связь с историей UTM пользователя.
		edge.To("utm_histories", UtmHistory.Type),
		// bonuses - связь с бонусами пользователя.
		edge.To("bonuses", Bonus.Type),
	}
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		StandardMixin{Prefix: UserPrefix},
	}
}

// Indexes of the User.
// phone and email are unique only among non-deleted users, so that a soft-deleted
// record does not block re-registration with the same phone/email.
func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("phone").
			Unique().
			Annotations(entsql.IndexWhere("deleted_at IS NULL")),
		index.Fields("email").
			Unique().
			Annotations(entsql.IndexWhere("deleted_at IS NULL")),
	}
}

// Аннотации User.
func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "users",
		},
	}
}
