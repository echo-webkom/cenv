package internal

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadAndCheck(envPath, schemaPath string) error {
	env, err := ReadSchema(schemaPath)
	if err != nil {
		return err
	}

	if err := godotenv.Load(envPath); err != nil {
		return err
	}

	for k, v := range env.Fields {
		v.value = os.Getenv(k)
		if err := validateEnvField(v); err != nil {
			return err
		}
		if err := compareFields(v, v); err.Error() != nil {
			return err.Error()
		}
	}

	return nil
}
