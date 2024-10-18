package test

import (
	"testing"

	"github.com/echo-webkom/cenv/cenv"
)

func TestUpdate(t *testing.T) {
	env, err := cenv.ReadEnv("cases/update_cases/.env")
	if err != nil {
		t.Fatal(err)
	}

	schema, err := cenv.ReadSchema("cases/update_cases/cenv.schema.json")
	if err != nil {
		t.Fatal(err)
	}

	if err := cenv.ValidateSchema(env, schema, "test/cases/update_cases/.env"); err != nil {
		t.Fatal(err)
	}
}
