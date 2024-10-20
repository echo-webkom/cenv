package test

import (
	"fmt"
	"testing"

	"github.com/echo-webkom/cenv/cenv"
)

func TestCheck(t *testing.T) {
	numTestCases := 5

	for i := 0; i < numTestCases; i++ {
		env, err := cenv.ReadEnv(fmt.Sprintf("cases/check_cases/%d.env", i))
		if err != nil {
			t.Fatal(err)
		}

		schema, err := cenv.ReadSchema(fmt.Sprintf("cases/check_cases/%d.json", i))
		if err != nil {
			t.Fatal(err)
		}

		if cenv.ValidateSchema(env, schema, "") == nil {
			t.Fatalf("expected test case idx %d to fail", i)
		}
	}
}
