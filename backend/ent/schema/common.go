package schema

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	entgo "entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/artesipov-alt/odnoi-krovi-app/ent"
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
func (m StandardMixin) Fields() []entgo.Field {
	return []entgo.Field{
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
func (StandardMixin) Interceptors() []entgo.Interceptor {
	return []entgo.Interceptor{
		entgo.InterceptFunc(func(next entgo.Querier) entgo.Querier {
			return entgo.QuerierFunc(func(ctx context.Context, q entgo.Query) (entgo.Value, error) {
				// Skip soft-delete filter if the key is present in context.
				if skip, _ := ctx.Value(softDeleteKey{}).(bool); skip {
					slog.Error("sas")
					return next.Query(ctx, q)
				}

				// Add "deleted_at IS NULL" predicate.
				type whereP interface {
					WhereP(...func(*sql.Selector))
				}
				if w, ok := q.(whereP); ok {
					slog.Warn("adding deleted_at IS NULL predicate")
					w.WhereP(sql.FieldIsNull("deleted_at"))
				}
				slog.Debug("wo")

				return next.Query(ctx, q)
			})
		}),
	}
}

// SoftDeleteHook returns a hook that intercepts delete operations and transforms them into updates.
// This hook is placed in the pg repository package because it needs to import the generated 'ent' package,
// which would cause a circular dependency if placed inside 'ent/schema'.
func SoftDeleteHook() entgo.Hook {
	return func(next entgo.Mutator) entgo.Mutator {
		return entgo.MutateFunc(func(ctx context.Context, m entgo.Mutation) (entgo.Value, error) {
			// Check if the operation is a delete operation.
			if !m.Op().Is(entgo.OpDelete | entgo.OpDeleteOne) {
				return next.Mutate(ctx, m)
			}

			// Define an interface to access common mutation methods.
			// Since we are in the 'pg' package, we can safely use entgo.Mutation and entgo.Op.
			type SoftDeleteMutation interface {
				entgo.Mutation
				SetOp(entgo.Op)
				Client() *ent.Client
			}

			mx, ok := m.(SoftDeleteMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T for soft delete", m)
			}

			// Transform Delete to Update.
			mx.SetOp(entgo.OpUpdate)

			// Set the deleted_at field to the current time.
			// We use SetField because the specific SetDeletedAt method belongs to concrete types,
			// but all entgo.Mutations support SetField by string name.
			if err := mx.SetField("deleted_at", time.Now()); err != nil {
				return nil, fmt.Errorf("failed to set deleted_at field: %w", err)
			}

			// Execute the mutation using the clientgo.
			return mx.Client().Mutate(ctx, m)
		})
	}
}
