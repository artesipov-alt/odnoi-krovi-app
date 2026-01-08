package schema

import (
	"context"
	"fmt"
	"reflect"
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

// SkipSoftDelete returns a new context that skips the soft-delete interceptor/mutators.
func SkipSoftDelete(parent context.Context) context.Context {
	return context.WithValue(parent, softDeleteKey{}, true)
}

// StandardMixin implements the ID generation with prefix,
// time auditing, and soft delete pattern.
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

// Hooks of the StandardMixin.
func (StandardMixin) Hooks() []ent.Hook {
	return []ent.Hook{
		softDeleteHook(),
	}
}

// Interceptors of the StandardMixin.
func (StandardMixin) Interceptors() []ent.Interceptor {
	return []ent.Interceptor{
		softDeleteInterceptor(),
	}
}

// --- Shared Soft Delete Logic ---

func softDeleteInterceptor() ent.Interceptor {
	return interceptFunc(func(ctx context.Context, q ent.Query) error {
		if skip, _ := ctx.Value(softDeleteKey{}).(bool); skip {
			return nil
		}
		addSoftDeletePredicate(q)
		return nil
	})
}

func softDeleteHook() ent.Hook {
	return func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			if skip, _ := ctx.Value(softDeleteKey{}).(bool); skip {
				return next.Mutate(ctx, m)
			}
			if !m.Op().Is(ent.OpDelete | ent.OpDeleteOne) {
				return next.Mutate(ctx, m)
			}
			type SoftDeleteMutation interface {
				ent.Mutation
				SetOp(ent.Op)
				SetDeletedAt(time.Time)
				WhereP(...func(*sql.Selector))
			}
			mx, ok := m.(SoftDeleteMutation)
			if !ok {
				return nil, fmt.Errorf("unexpected mutation type %T", m)
			}

			// Use reflection to call Client() to avoid circular dependency and return type mismatch.
			rv := reflect.ValueOf(m)
			method := rv.MethodByName("Client")
			if !method.IsValid() {
				return nil, fmt.Errorf("mutation type %T does not implement Client()", m)
			}
			client := method.Call(nil)[0].Interface()

			addSoftDeletePredicate(mx)
			mx.SetOp(ent.OpUpdate)
			mx.SetDeletedAt(time.Now())

			type Mutator interface {
				Mutate(context.Context, ent.Mutation) (ent.Value, error)
			}
			mClient, ok := client.(Mutator)
			if !ok {
				return nil, fmt.Errorf("unexpected client type %T", client)
			}
			return mClient.Mutate(ctx, m)
		})
	}
}

func addSoftDeletePredicate(w any) {
	type whereP interface {
		WhereP(...func(*sql.Selector))
	}
	if wp, ok := w.(whereP); ok {
		wp.WhereP(
			sql.FieldIsNull("deleted_at"),
		)
	}
}

// interceptFunc is an adapter to allow the use of ordinary functions as interceptors.
type interceptFunc func(context.Context, ent.Query) error

// Intercept implements the ent.Interceptor interface.
func (f interceptFunc) Intercept(next ent.Querier) ent.Querier {
	return ent.QuerierFunc(func(ctx context.Context, q ent.Query) (ent.Value, error) {
		if err := f(ctx, q); err != nil {
			return nil, err
		}
		return next.Query(ctx, q)
	})
}

// Traverse implements the ent.Traverser interface.
func (f interceptFunc) Traverse(ctx context.Context, q ent.Query) error {
	return f(ctx, q)
}
