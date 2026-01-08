package pg

import (
	"context"
	"fmt"
	"time"

	"github.com/artesipov-alt/odnoi-krovi-app/ent"
)

// SoftDeleteHook returns a hook that intercepts delete operations and transforms them into updates.
// This hook is placed in the pg repository package because it needs to import the generated 'ent' package,
// which would cause a circular dependency if placed inside 'ent/schema'.
func SoftDeleteHook() ent.Hook {
	return func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			// Check if the operation is a delete operation.
			if !m.Op().Is(ent.OpDelete | ent.OpDeleteOne) {
				return next.Mutate(ctx, m)
			}

			// Define an interface to access common mutation methods.
			// Since we are in the 'pg' package, we can safely use ent.Mutation and ent.Op.
			type SoftDeleteMutation interface {
				ent.Mutation
				SetOp(ent.Op)
				Client() *ent.Client
			}

			mx, ok := m.(SoftDeleteMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T for soft delete", m)
			}

			// Transform Delete to Update.
			mx.SetOp(ent.OpUpdate)

			// Set the deleted_at field to the current time.
			// We use SetField because the specific SetDeletedAt method belongs to concrete types,
			// but all ent.Mutations support SetField by string name.
			if err := mx.SetField("deleted_at", time.Now()); err != nil {
				return nil, fmt.Errorf("failed to set deleted_at field: %w", err)
			}

			// Execute the mutation using the client.
			return mx.Client().Mutate(ctx, m)
		})
	}
}
