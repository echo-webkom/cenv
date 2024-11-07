package cenv

import (
	"fmt"
	"io/fs"
	"os"
	"strings"
	"time"
)

type CenvFile struct {
	LastUpdated time.Time            `json:"lastUpdated"`
	Fields      map[string]CenvField `json:"fields"`
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

// Fix inserts missing fields into your env and writes values if public.
func Fix(envPath, schemaPath string) error {
	schema, err := ReadSchema(schemaPath)
	if err != nil {
		return err
	}

	env, _ := ReadEnv(envPath)

	file := strings.Builder{}

	for _, f := range schema.Fields {
		s := ""

		if f.Required {
			s += "# @required\n"
		}
		if f.Public {
			s += "# @public\n"
		}
		if f.LengthRequired {
			s += fmt.Sprintf("# @length %d\n", f.Length)
		}

		if v, ok := env[f.Key]; ok && !f.Public {
			s += fmt.Sprintf("%s=%s\n", f.Key, v.value)
		} else {
			s += fmt.Sprintf("%s=%s\n", f.Key, f.Value)
		}

		fmt.Printf("cenv: added '%s'\n", f.Key)
		file.WriteString(s)
	}

	if err := os.WriteFile(envPath, []byte(file.String()), fs.FileMode(os.O_WRONLY)); err != nil {
		return err
	}

	return Check(envPath, schemaPath)
}
