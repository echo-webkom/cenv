package cenv

import (
	"encoding/json"
	"fmt"
	"os"
)

func ReadSchema(filepath string) (schema CenvFile, err error) {
	file, err := os.ReadFile(filepath)
	if err != nil {
		return schema, err
	}

	schema = CenvFile{}
	if err = json.Unmarshal(file, &schema); err != nil {
		return schema, err
	}

	return schema, err
}

func writeShema(env CenvFile, filepath string) error {
	b, err := json.Marshal(env)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath, b, os.ModeAppend|os.ModePerm)
}

func ValidateSchema(env []CenvField, schema CenvFile) error {
	envMap := make(map[string]CenvField)
	for _, f := range env {
		envMap[f.Key] = f
	}

	for idx, f := range schema.Fields {
		if ff, ok := envMap[f.Key]; ok && idx <= len(env) {
			if err := validateField(f, ff); err != nil {
				return err
			}

			continue
		}

		return fmt.Errorf("missing field '%s'", f.Key)
	}

	return nil
}

func validateField(sf CenvField, ef CenvField) error {
	if sf.Required && !ef.Required {
		return fmt.Errorf("field '%s' is tagged as required in schema, but not in env", sf.Key)
	}

	if !sf.Required && ef.Required {
		return fmt.Errorf("field '%s' is marked as required in .env, but is not in schema", ef.Key)
	}

	if sf.LengthRequired && !ef.LengthRequired {
		return fmt.Errorf("field '%s' is tagged with length %d in schema, but is not in .env", sf.Key, sf.Length)
	}

	if !sf.LengthRequired && ef.LengthRequired {
		return fmt.Errorf("field '%s' is marked with length %d in .env, but not in schema", sf.Key, sf.Length)
	}

	if sf.Length != ef.Length {
		return fmt.Errorf("field '%s' is tagged with length %d in schema, but is %d in .env", sf.Key, sf.Length, ef.Length)
	}

	return nil
}
