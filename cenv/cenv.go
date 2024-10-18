package cenv

type CenvFile []CenvField

type CenvField struct {
	Required bool   `json:"required"`
	Key      string `json:"key"`
	value    string
}

// Update generates a cenv.schema.json file based on the given .env file
func Update(envPath string, schemaPath string) error {
	env, err := ReadEnv(envPath)
	if err != nil {
		return err
	}

	return writeShema(env, schemaPath)
}

// Check validates the .env file based on the schema
func Check(envPath string, schemaPath string) error {
	env, err := ReadEnv(envPath)
	if err != nil {
		return err
	}

	schema, err := ReadSchema(schemaPath)
	if err != nil {
		return err
	}

	return ValidateSchema(env, schema, envPath)
}
