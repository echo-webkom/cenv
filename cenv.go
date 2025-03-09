package cenv

import "github.com/echo-webkom/cenv/cenv"

const (
	defaultEnvPath    = ".env"
	defaultSchemaPath = "cenv.schema.json"
)

// CheckAndLoad loads the .env file into the process environment and verifies
// that all values adhere to the schema rules. Variables defined in .env will
// override any pre-existing ones.
//
// The default path for both .env and cenv.schema.json are used, which are
// both in the project root directory. To specify path use CheckAndLoadEx()
func CheckAndLoad() error {
	return cenv.CheckAndLoad(defaultEnvPath, defaultSchemaPath)
}

func CheckAndLoadEx(envPath string, schemaPath string) error {
	return cenv.CheckAndLoad(envPath, schemaPath)
}
