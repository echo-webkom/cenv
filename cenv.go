package cenv

import "github.com/echo-webkom/cenv/internal"

const (
	defaultEnvPath    = ".env"
	defaultSchemaPath = "cenv.schema.json"
)

// Load loads the .env file into the process environment and verifies
// that all values adhere to the schema rules. Variables defined in .env will
// override any pre-existing ones.
//
// The default path for both .env and cenv.schema.json are used, which are
// both in the project root directory. To specify path use LoadEx()
func Load() error {
	return internal.LoadAndCheck(defaultEnvPath, defaultSchemaPath, false)
}

func LoadEx(envPath string, schemaPath string) error {
	return internal.LoadAndCheck(envPath, schemaPath, false)
}

// Verify that the values in the loaded environment mathces the schema.
func Verify() error {
	return internal.LoadAndCheck(defaultEnvPath, defaultSchemaPath, true)
}
