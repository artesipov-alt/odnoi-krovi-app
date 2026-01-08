package schema

import (
	"context"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/jaevor/go-nanoid"
)

// Prefix constants for ID generation
const (
	UserPrefix         = "USR"
	PetPrefix          = "PET"
	PetHealthPrefix    = "PHL"
	PetTreatmentPrefix = "PTR"
	PetAnalysisPrefix  = "PAN"
	PetBonusPrefix     = "PBN"
	BloodSearchPrefix  = "BLS"
)

// generateID generates a new ID with prefix and nanoID of 10 characters
func generateID(prefix string) string {
	gen, _ := nanoid.Standard(10)
	id := gen()
	return prefix + "-" + id
}

type softDeleteKey struct{}

// SkipSoftDelete returns a new context that skips the soft-delete interceptor.
func SkipSoftDelete(parent context.Context) context.Context {
	return context.WithValue(parent, softDeleteKey{}, true)
}

// StandardMixin implements the ID generation with prefix,
// time auditing, and soft delete fields/interceptors.
type StandardMixin struct {
	mixin.Schema
	Prefix string
}

// Fields of the StandardMixin.
func (m StandardMixin) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			Unique().
			Immutable().
			DefaultFunc(func() string { return generateID(m.Prefix) }).
			StructTag(`json:"id"`),
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			StructTag(`json:"createdAt"`),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			StructTag(`json:"updatedAt"`),
		field.Time("deleted_at").
			Optional().
			Nillable().
			StructTag(`json:"deletedAt"`),
	}
}

// Interceptors of the StandardMixin.
func (StandardMixin) Interceptors() []ent.Interceptor {
	return []ent.Interceptor{
		ent.InterceptFunc(func(next ent.Querier) ent.Querier {
			return ent.QuerierFunc(func(ctx context.Context, q ent.Query) (ent.Value, error) {
				// Skip soft-delete filter if the key is present in context.
				if skip, _ := ctx.Value(softDeleteKey{}).(bool); skip {
					return next.Query(ctx, q)
				}

				// Add "deleted_at IS NULL" predicate.
				type whereP interface {
					WhereP(...func(*sql.Selector))
				}
				if w, ok := q.(whereP); ok {
					w.WhereP(sql.FieldIsNull("deleted_at"))
				}

				return next.Query(ctx, q)
			})
		}),
	}
}
