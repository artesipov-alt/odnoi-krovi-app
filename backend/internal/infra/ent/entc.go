//go:build ignore

package main

import (
	"log"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

func main() {
	// Настраиваем расширение для генерации Input-структур
	ex, err := entgql.NewExtension(
		// Сообщаем Ent, что нам нужны типы Input для Create и Update
		entgql.WithConfigPath("./gqlgen.yml"), // Файл может не существовать, это ок
		entgql.WithSchemaGenerator(),
		entgql.WithWhereInputs(true),
	)
	if err != nil {
		log.Fatalf("creating entgql extension: %v", err)
	}

	// Запускаем генерацию с твоими фичами (intercept, snapshot)
	opts := []entc.Option{
		entc.Extensions(ex),
		entc.FeatureNames("intercept", "schema/snapshot", "sql/upsert"),
	}

	if err := entc.Generate("./schema", &gen.Config{}, opts...); err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}
