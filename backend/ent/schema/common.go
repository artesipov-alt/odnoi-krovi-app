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
)

// generateID generates a new ID with prefix and nanoID of 10 characters
func generateID(prefix string) string {
	gen, _ := nanoid.Standard(10)
	id := gen()
	return prefix + "-" + id
}

type softDeleteKey struct{}

// SkipSoftDelete возвращает контекст, который заставляет пропустить фильтрацию удаленных записей.
func SkipSoftDelete(parent context.Context) context.Context {
	return context.WithValue(parent, softDeleteKey{}, true)
}

// AuditMixin предоставляет поля времени создания, обновления и мягкого удаления.
// Используйте этот миксин для сущностей с автоинкрементным ID.
type AuditMixin struct {
	mixin.Schema
}

// Fields возвращает поля аудита.
func (AuditMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").Default(time.Now).Immutable().StructTag(`json:"createdAt"`),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now).StructTag(`json:"updatedAt"`),
		field.Time("deleted_at").Optional().Nillable().StructTag(`json:"deletedAt"`),
	}
}

// Interceptors автоматически добавляет фильтр "deleted_at IS NULL" ко всем запросам.
func (AuditMixin) Interceptors() []ent.Interceptor {
	return []ent.Interceptor{
		ent.TraverseFunc(func(ctx context.Context, q ent.Query) error {
			// Если в контексте есть SkipSoftDelete, показываем все записи (включая удаленные).
			if skip, _ := ctx.Value(softDeleteKey{}).(bool); skip {
				return nil
			}

			// Добавляем фильтр WHERE deleted_at IS NULL.
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

// BaseMixin предоставляет кастомную генерацию ID (NanoID с префиксом) вместе с полями аудита.
// Используйте NewBaseMixin для сущностей, которым нужен строковый уникальный ID.
type BaseMixin struct {
	AuditMixin
	prefix string
}

// NewBaseMixin создает новый BaseMixin с указанным префиксом для ID.
func NewBaseMixin(prefix string) BaseMixin {
	return BaseMixin{prefix: prefix}
}

// Fields возвращает поле ID и поля аудита.
func (b BaseMixin) Fields() []ent.Field {
	return append([]ent.Field{
		field.String("id").
			Unique().
			Immutable().
			DefaultFunc(func() string { return generateID(b.prefix) }).
			StructTag(`json:"id"`),
	}, b.AuditMixin.Fields()...)
}
