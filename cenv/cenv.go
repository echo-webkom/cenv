package cenv

import "time"

type CenvFile struct {
	LastUpdated time.Time   `json:"lastUpdated"`
	Fields      []CenvField `json:"fields"`
}

type CenvField struct {
	Required       bool   `json:"required"`
	LengthRequired bool   `json:"lengthRequired"`
	Length         uint32 `json:"length"`
	Key            string `json:"key"`

	value string
}

// Update generates a cenv.schema.json file based on the given .env file
func Update(envPath, schemaPath string) error {
	fields, err := ReadEnv(envPath)
	if err != nil {
		return err
	}

	env := CenvFile{
		LastUpdated: time.Now(),
		Fields:      fields,
	}

	return writeShema(env, schemaPath)
}

// Check validates the .env file based on the schema
func Check(envPath, schemaPath string) error {
	fields, err := ReadEnv(envPath)
	if err != nil {
		return err
	}

	schema, err := ReadSchema(schemaPath)
	if err != nil {
		return err
	}

	return ValidateSchema(fields, schema)
}
