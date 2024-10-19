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

func ValidateSchema(env CenvFile, schema CenvFile, envPath string) error {
	envMap := make(map[string]CenvField)
	for _, f := range env {
		envMap[f.Key] = f
	}

	for idx, f := range schema {
		if ff, ok := envMap[f.Key]; ok && idx <= len(env) {
			if err := validateField(f, ff); err != nil {
				return err
			}

			continue
		}

		return fmt.Errorf("missing field '%s' in %s", f.Key, envPath)
	}

	return nil
}

func validateField(sf CenvField, ef CenvField) error {
	if sf.Required {
		if !ef.Required || ef.value == "" {
			return fmt.Errorf("field '%s' is required in .env", sf.Key)
		}
	}

	if !sf.Required && ef.Required {
		return fmt.Errorf("field '%s' is marked as required in .env, but is not in schema", ef.Key)
	}

	if sf.LengthRequired {
		if !ef.LengthRequired {
			return fmt.Errorf("field '%s' is tagged with length %d in schema, but is not in .env", sf.Key, sf.Length)
		}

		// Do not reorder, sf and ef length comparison must come first
		if sf.Length != ef.Length {
			return fmt.Errorf("field '%s' is tagged with length %d in schema, but is %d in .env", sf.Key, sf.Length, ef.Length)
		}

		if int(sf.Length) != len(ef.value) {
			return fmt.Errorf("value of '%s' in .env must be %d bytes, is %d", ef.Key, sf.Length, len(ef.value))
		}
	}

	if !sf.LengthRequired && ef.LengthRequired {
		return fmt.Errorf("field '%s' is marked with length %d in .env, but not in schema", sf.Key, sf.Length)
	}

	return nil
}
