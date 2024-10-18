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
			if err := validateField(f, ff, envPath); err != nil {
				return err
			}

			continue
		}

		return fmt.Errorf("missing field '%s' in %s", f.Key, envPath)
	}

	return nil
}

func validateField(sf CenvField, ef CenvField, envPath string) error {
	if sf.Required {
		if !ef.Required || ef.value == "" {
			return fmt.Errorf("field '%s' is required in %s", sf.Key, envPath)
		}
	}

	if !sf.Required && ef.Required {
		return fmt.Errorf("field '%s' is marked as required %s, but is not", ef.Key, envPath)
	}

	if sf.Length != 0 && sf.Length != ef.Length {
		return fmt.Errorf("value of '%s' in %s must be %d bytes", ef.Key, envPath, sf.Length)
	}

	return nil
}
