package schema

import (
	"context"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/artesipov-alt/odnoi-krovi-app/internal/entutil"
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

// SkipSoftDelete returns a new context that skips the soft-delete interceptor/mutators.
// It proxies the call to the internal entutil package.
func SkipSoftDelete(parent context.Context) context.Context {
	return entutil.SkipSoftDelete(parent)
}

// StandardMixin implements the ID generation with prefix,
// time auditing, and soft delete pattern.
type StandardMixin struct {
	mixin.Schema
	Prefix string
}

// Mixin returns the list of mixins for the StandardMixin.
// We include the SoftDeleteMixin from entutil to get the hooks and interceptors.
func (StandardMixin) Mixin() []ent.Mixin {
	return []ent.Mixin{
		entutil.SoftDeleteMixin{},
	}
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
		// We define the field here so Ent generates the necessary methods
		// (ClearDeletedAt, DeletedAtNotNil, etc.) for each entity.
		field.Time("deleted_at").
			Optional().
			Nillable().
			StructTag(`json:"deletedAt"`),
	}
}
