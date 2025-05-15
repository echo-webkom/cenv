package internal

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type CenvFile struct {
	LastUpdated time.Time            `json:"lastUpdated"`
	Fields      map[string]CenvField `json:"fields"`
}

type CenvField struct {
	Required       bool     `json:"required,omitempty"`       // This field has to be present and have a non-empty value
	Public         bool     `json:"public,omitempty"`         // This has a publicly known but required value, stored in the schema
	LengthRequired bool     `json:"lengthRequired,omitempty"` // The length of this field is specified in the schema
	Length         uint32   `json:"length,omitempty"`         // The required length, if LengthRequired is true
	Format         string   `json:"format,omitempty"`         // Require a specified format for the value
	Key            string   `json:"key,omitempty"`            // Field name
	Value          string   `json:"value,omitempty"`          // Public only
	Enum           []string `json:"enum,omitempty"`           // List of legal enum values
	Default        string   `json:"default,omitempty"`        // Default value if none is given

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
		if f.Format != "" {
			s += fmt.Sprintf("# @format %s\n", f.Format)
		}
		if len(f.Enum) != 0 {
			s += fmt.Sprintf("# @enum %s\n", strings.Join(f.Enum, " | "))
		}
		if f.Default != "" {
			s += fmt.Sprintf("# @default %s\n", f.Default)
		}

		// Public > Private > Default

		if f.Public {
			// The value is public (and static)
			s += fmt.Sprintf("%s=%s\n", f.Key, f.Value)
		} else if v, ok := env[f.Key]; ok && v.value != "" {
			// There is already a value here
			s += fmt.Sprintf("%s=%s\n", f.Key, v.value)
		} else {
			// The value is either empty or set to default (if any)
			s += fmt.Sprintf("%s=%s\n", f.Key, f.Default)
		}

		fmt.Printf("cenv: added '%s'\n", f.Key)
		file.WriteString(s)
		file.WriteByte('\n')
	}

	if err := os.WriteFile(envPath, []byte(file.String()), 0666); err != nil {
		return err
	}

	return Check(envPath, schemaPath)
}
