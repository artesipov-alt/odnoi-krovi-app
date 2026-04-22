package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Bonus holds the schema definition for the Bonus entity.
type Bonus struct {
	ent.Schema
}

// Fields of the Bonus.
func (Bonus) Fields() []ent.Field {
	return []ent.Field{
		// user_id — Идентификатор пользователя, которому принадлежит бонус
		field.String("user_id").
			Optional(),

		// donor_response_id — Идентификатор отклика донора, из-за которого получен бонус
		field.String("donor_response_id").
			Optional(),

		// partner_name — Наименование партнера (юрлицо / бренд)
		field.String("partner_name"),

		// description — Описание бонуса
		field.Text("description"),

		// target — Для кого (кошка/собака/все)
		field.Enum("target").
			Values("cat", "dog", "all"),

		// recipient — Для кого (донор/реципиент/все)
		field.Enum("recipient").
			Values("donor", "recipient", "all"),

		// category — Категория (корма / препараты / другое)
		field.Enum("category").
			Values("food", "preparation", "other"),

		// subcategory — Подкатегория (например, для кормов: сухой/влажный)
		field.String("subcategory").
			Optional(),

		// promo_code — Промокод
		field.String("promo_code").
			Unique(),

		// expires_at — Дата окончания срока действия
		field.Time("expires_at"),

		// platform_name — Площадка применения (название магазина, маркетплейса)
		field.String("platform_name"),

		// platform_url — Ссылка на страницу площадки для активации промокода
		field.String("platform_url").
			Optional(),

		// stage — Статус бонуса
		field.Enum("stage").
			Values("stage", "used", "unused", "reserved").
			Default("unused"),
	}
}

// Edges of the Bonus.
func (Bonus) Edges() []ent.Edge {
	return []ent.Edge{
		// user — Связь с пользователем
		edge.From("user", User.Type).
			Ref("bonuses").
			Field("user_id").
			Unique(),

		// donor_response — Связь с откликом донора
		edge.From("donor_response", DonorResponse.Type).
			Ref("bonuses").
			Field("donor_response_id").
			Unique(),
	}
}

// Annotations of the Bonus.
func (Bonus) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "bonuses",
		},
	}
}

func (Bonus) Mixin() []ent.Mixin {
	return []ent.Mixin{
		StandardMixin{Prefix: BonusPrefix},
	}
}
