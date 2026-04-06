//go:build ignore

package main

import (
	"log"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

func main() {

	// Запускаем генерацию с твоими фичами (intercept, snapshot)
	opts := []entc.Option{
		entc.FeatureNames("intercept", "schema/snapshot", "sql/upsert"),
	}

	if err := entc.Generate("./schema", &gen.Config{}, opts...); err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}
