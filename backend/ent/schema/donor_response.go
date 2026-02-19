package schema

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

const DonorResponsePrefix = "RES"

type DonorResponse struct {
	ent.Schema
}

// Fields of the BloodSearchRequest.
func (DonorResponse) Fields() []ent.Field {
	return []ent.Field{
		field.Int32("amount_ml"),
	}
}

// Edges of the BloodSearchRequest.
func (DonorResponse) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("request", BloodSearchRequest.Type).Ref("responses").Unique().Required(),
		edge.To("donor", Pet.Type).Unique().Required(),
	}
}

func (DonorResponse) Mixin() []ent.Mixin {
	return []ent.Mixin{
		StandardMixin{Prefix: DonorResponsePrefix},
	}
}

// Annotations of the BloodSearchRequest.
func (DonorResponse) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{
			Table: "donor_responses",
		},
		entgql.Mutations(
			entgql.MutationCreate(),
			entgql.MutationUpdate(),
		),
	}

}
