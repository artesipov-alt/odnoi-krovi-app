package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// DonorPreference содержит определение схемы для сущности DonorPreference.
// Эта сущность хранит предпочтения донора по умолчанию для ответов на донации крови.
type DonorPreference struct {
	ent.Schema
}

// Поля DonorPreference.
func (DonorPreference) Fields() []ent.Field {
	return []ent.Field{
		// user_id - внешний ключ к пользователю
		field.String("user_id"),
		// preferred_location_ids - список ID локаций, где донор хочет помочь
		field.JSON("preferred_location_ids", []string{}).
			Optional(),
		// recovery_period_months - период восстановления между донациями (2-6 месяцев)
		field.Int("recovery_period_months").
			Optional().
			Default(2),
		// compensation_type указывает предпочтение донора по компенсации
		field.Enum("compensation_type").
			Values("free", "paid", "food").
			Optional(),
		// taxi_compensation указывает, нужна ли донору компенсация за такси
		field.Bool("taxi_compensation").
			Default(false),
		// notification_frequency указывает, как часто донор хочет получать уведомления
		field.Enum("notification_frequency").
			Values("immediately", "daily", "weekly", "never").
			Default("immediately"),
		// open_for_contact - разрешает ли донор получать приглашения от реципиентов
		field.Bool("open_for_contact").
			Default(false).
			Comment("Разрешает реципиентам находить донора как потенциального и приглашать"),
	}
}

// Рёбра DonorPreference.
func (DonorPreference) Edges() []ent.Edge {
	return []ent.Edge{
		// user - ребро к пользователю, которому принадлежат эти предпочтения
		edge.From("user", User.Type).
			Ref("donor_preference").
			Unique().
			Required().
			Field("user_id"),
	}
}

func (DonorPreference) Mixin() []ent.Mixin {
	return []ent.Mixin{
		StandardMixin{Prefix: DonorPreferencePrefix},
	}
}

// Аннотации DonorPreference.
func (DonorPreference) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "donor_preferences",
		},
	}
}
