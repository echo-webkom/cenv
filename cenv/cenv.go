package cenv

import "time"

type CenvFile struct {
	LastUpdated time.Time   `json:"lastUpdated"`
	Fields      []CenvField `json:"fields"`
}

type CenvField struct {
	Required       bool   `json:"required"`       // This field has to be present and have a non-empty value
	Public         bool   `json:"public"`         // This has a publicly known but required value, stored in the schema
	LengthRequired bool   `json:"lengthRequired"` // The length of this field is specified in the schema
	Length         uint32 `json:"length"`         // The required length, if LengthRequired is true
	Key            string `json:"key"`            // Field name
	Value          string `json:"value"`          // Public only

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
