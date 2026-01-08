package entutil

import (
	"context"
	"fmt"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/schema/mixin"
)

type softDeleteKey struct{}

// SkipSoftDelete returns a new context that skips the soft-delete interceptor/mutators.
func SkipSoftDelete(parent context.Context) context.Context {
	return context.WithValue(parent, softDeleteKey{}, true)
}

// SoftDeleteMixin implements the soft delete pattern for schemas.
// It provides hooks and interceptors but expects the "deleted_at" field
// to be defined in the schema that embeds it.
type SoftDeleteMixin struct {
	mixin.Schema
}

// Interceptors of the SoftDeleteMixin.
func (d SoftDeleteMixin) Interceptors() []ent.Interceptor {
	return []ent.Interceptor{
		ent.InterceptFunc(func(next ent.Querier) ent.Querier {
			return ent.QuerierFunc(func(ctx context.Context, q ent.Query) (ent.Value, error) {
				if skip, _ := ctx.Value(softDeleteKey{}).(bool); skip {
					return next.Query(ctx, q)
				}
				d.addPredicate(q)
				return next.Query(ctx, q)
			})
		}),
	}
}

// Hooks of the SoftDeleteMixin.
func (d SoftDeleteMixin) Hooks() []ent.Hook {
	return []ent.Hook{
		func(next ent.Mutator) ent.Mutator {
			return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
				if skip, _ := ctx.Value(softDeleteKey{}).(bool); skip {
					return next.Mutate(ctx, m)
				}
				if !m.Op().Is(ent.OpDelete | ent.OpDeleteOne) {
					return next.Mutate(ctx, m)
				}

				// Define an interface that matches the generated mutation types.
				type SoftDeleteMutation interface {
					ent.Mutation
					SetOp(ent.Op)
					SetDeletedAt(time.Time)
					WhereP(...func(*sql.Selector))
					Client() any
				}

				mx, ok := m.(SoftDeleteMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}

				d.addPredicate(mx)
				mx.SetOp(ent.OpUpdate)
				mx.SetDeletedAt(time.Now())

				// Since we can't import 'ent' package here to get the concrete Client,
				// we use an interface to call Mutate on whatever Client() returns.
				type Mutator interface {
					Mutate(context.Context, ent.Mutation) (ent.Value, error)
				}
				mClient, ok := mx.Client().(Mutator)
				if !ok {
					return nil, fmt.Errorf("unexpected client type %T", mx.Client())
				}
				return mClient.Mutate(ctx, m)
			})
		},
	}
}

func (d SoftDeleteMixin) addPredicate(w any) {
	type whereP interface {
		WhereP(...func(*sql.Selector))
	}
	if wp, ok := w.(whereP); ok {
		wp.WhereP(
			sql.FieldIsNull("deleted_at"),
		)
	}
}
