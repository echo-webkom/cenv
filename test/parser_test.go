package test

import (
	"cenv/cenv"
	"fmt"
	"testing"
)

func TestParser(t *testing.T) {
	numTestFiles := 5

	for i := 0; i < numTestFiles; i++ {
		_, err := cenv.ReadEnv(fmt.Sprintf("cases/parser_cases/%d.env", i))
		if err != nil {
			t.Fatal(err)
		}
	}
}
