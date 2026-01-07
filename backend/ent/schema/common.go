package schema

import (
	"context"
	"time"

	"entgo.io/ent"
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

// BaseMixin provides common fields and hooks for all schemas
type BaseMixin struct {
	mixin.Schema
	prefix string
}

// NewBaseMixin creates a new BaseMixin with the specified prefix
func NewBaseMixin(prefix string) BaseMixin {
	return BaseMixin{prefix: prefix}
}

func (b BaseMixin) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique().Immutable().DefaultFunc(func() string { return generateID(b.prefix) }).StructTag(`json:"id"`),
		field.Time("created_at").Default(time.Now).StructTag(`json:"createdAt"`),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now).StructTag(`json:"updatedAt"`),
		field.Time("deleted_at").Optional().Nillable().StructTag(`json:"deletedAt"`),
	}
}
