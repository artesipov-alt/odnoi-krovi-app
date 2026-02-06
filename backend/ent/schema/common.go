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
	"github.com/artesipov-alt/odnoi-krovi-app/ent/intercept"
	"github.com/jaevor/go-nanoid"
	sloghttp "github.com/samber/slog-http"
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
	gen, _ := nanoid.Custom("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz", 10)
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
		intercept.TraverseFunc(func(ctx context.Context, q intercept.Query) error {
			if IsSkipSoftDelete(ctx) {
				return nil
			}
			q.WhereP(sql.FieldIsNull("deleted_at"))
			return nil
		}),
	}
}

// SoftDeleteHook returns a hook that intercepts delete operations and transforms them into updates.
// This hook is placed in the pg repository package because it needs to import the generated 'ent' package,
// which would cause a circular dependency if placed inside 'ent/schema'.
func SoftDeleteHook() entgo.Hook {
	return func(next entgo.Mutator) entgo.Mutator {
		return entgo.MutateFunc(func(ctx context.Context, m entgo.Mutation) (entgo.Value, error) {
			// Check if the operation is a delete operation and if soft-delete should be skipped.
			if !m.Op().Is(entgo.OpDelete|entgo.OpDeleteOne) || IsSkipSoftDelete(ctx) {
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

// IsSkipSoftDelete checks if the soft-delete interceptor should be skipped.
func IsSkipSoftDelete(ctx context.Context) bool {
	skip, _ := ctx.Value(softDeleteKey{}).(bool)
	return skip
}

func DbInterceptor() ent.Interceptor {
	return ent.InterceptFunc(func(next ent.Querier) ent.Querier {
		return ent.QuerierFunc(func(ctx context.Context, q ent.Query) (ent.Value, error) {
			start := time.Now()
			res, err := next.Query(ctx, q)
			duration := time.Since(start)

			rid := sloghttp.GetRequestIDFromContext(ctx)

			// Порог медленного запроса (вынеси в конфиг потом)
			slowThreshold := 200 * time.Millisecond

			if duration > slowThreshold {
				// Логируем как предупреждение, если тормозит
				slog.WarnContext(ctx, "SLOW SQL Query",
					"rid", rid,
					"duration", duration,
					"threshold", slowThreshold,
				)
			} else {
				// Обычный лог в DEBUG
				slog.DebugContext(ctx, "SQL Query (Select)",
					"rid", rid,
					"duration", duration,
				)
			}

			return res, err
		})
	})
}
