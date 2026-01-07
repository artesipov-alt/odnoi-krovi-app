package schema

import (
	"context"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// BloodSearchRequest holds the schema definition for the BloodSearchRequest entity.
type BloodSearchRequest struct {
	ent.Schema
}

// Fields of the BloodSearchRequest.
func (BloodSearchRequest) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			Immutable().
			DefaultFunc(func() string { return generateID(BloodSearchPrefix) }).
			StructTag(`json:"id"`),
		field.String("pet_id").
			StructTag(`json:"petId"`),
		field.Int32("blood_volume_needed").
			StructTag(`json:"bloodVolumeNeeded"`),
		field.Int32("blood_volume_reserved").
			Default(0).
			StructTag(`json:"bloodVolumeReserved"`),
		field.JSON("regions", []int32{}).
			StructTag(`json:"regions"`),
		field.Bool("small_pets_notify_allowed").
			Default(true).
			StructTag(`json:"smallPetsNotifyAllowed"`),
		field.Enum("status").
			Values("active", "closed", "draft").
			Default("active").
			StructTag(`json:"status"`),
		field.String("description").
			Optional().
			StructTag(`json:"description"`),
		field.JSON("photo_urls", []string{}).
			Optional().
			StructTag(`json:"photoUrls"`),
		field.JSON("blood_group_ids", []int{}).
			Optional().
			StructTag(`json:"bloodGroupIds"`),
		field.JSON("blood_component_ids", []string{}).
			Optional().
			StructTag(`json:"bloodComponentIds"`),
		// Audit fields
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			StructTag(`json:"createdAt" swaggerignore:"true"`),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			StructTag(`json:"updatedAt" swaggerignore:"true"`),
		field.Time("deleted_at").
			Optional().
			Nillable().
			StructTag(`json:"deletedAt" swaggerignore:"true"`),
	}
}

// Edges of the BloodSearchRequest.
func (BloodSearchRequest) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("pet", Pet.Type).
			Ref("blood_search_request").
			Field("pet_id").
			Unique().
			Required(),
	}
}

// Interceptors of the BloodSearchRequest.
func (BloodSearchRequest) Interceptors() []ent.Interceptor {
	return []ent.Interceptor{
		ent.TraverseFunc(func(ctx context.Context, q ent.Query) error {
			// If SkipSoftDelete is in context, show all records (including deleted).
			if skip, _ := ctx.Value(softDeleteKey{}).(bool); skip {
				return nil
			}

			// Add WHERE deleted_at IS NULL filter.
			type query interface {
				WhereP(...func(*sql.Selector))
			}
			if w, ok := q.(query); ok {
				w.WhereP(func(s *sql.Selector) {
					s.Where(sql.IsNull(s.C("deleted_at")))
				})
			}
			return nil
		}),
	}
}
