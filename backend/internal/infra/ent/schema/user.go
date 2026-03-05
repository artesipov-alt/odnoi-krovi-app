package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// User определяет схему для сущности User.
// Эта сущность представляет пользователя в системе, связанного с Telegram.
type User struct {
	ent.Schema
}

// Поля User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		// telegram_id - уникальный идентификатор из Telegram.
		// DEPRECATED
		field.Int64("telegram_id").
			Optional().
			Unique(),
		// full_name - полное имя пользователя.
		field.String("full_name").
			Optional().
			MaxLen(255),
		// phone - номер телефона пользователя.
		field.String("phone").
			Optional().
			Unique().
			MaxLen(20),
		// email - адрес электронной почты пользователя.
		field.String("email").
			Optional().
			Unique().
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
			Values("user", "admin").
			Default("user"),
		// origin_source - UTM-метка, указывающая, откуда пришел пользователь.
		field.String("origin_source").
			Optional().
			MaxLen(255),
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
	}
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		StandardMixin{Prefix: UserPrefix},
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
