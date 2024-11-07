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

func ValidateSchema(env map[string]CenvField, schema CenvFile) error {
	for _, f := range schema.Fields {
		if ff, ok := env[f.Key]; ok {
			if err := validateField(f, ff); err != nil {
				return err
			}

			continue
		}

		return fmt.Errorf("missing field '%s'", f.Key)
	}

	return nil
}

func assertBoolEqual(key, name string, schema, env bool) error {
	if schema && !env {
		return fmt.Errorf("field '%s' is tagged with %s in schema, but not in env", key, name)
	}
	if !schema && env {
		return fmt.Errorf("field '%s' is tagged with %s in env, but not in schema", key, name)
	}
	return nil
}

func validateField(sf CenvField, ef CenvField) error {
	if err := assertBoolEqual(sf.Key, "required", sf.Required, ef.Required); err != nil {
		return err
	}
	if err := assertBoolEqual(sf.Key, "public", sf.Public, ef.Public); err != nil {
		return err
	}
	if err := assertBoolEqual(sf.Key, "a required length", sf.LengthRequired, ef.LengthRequired); err != nil {
		return err
	}
	if sf.Length != ef.Length {
		return fmt.Errorf("field '%s' is tagged with length %d in schema, but is %d in .env", sf.Key, sf.Length, ef.Length)
	}

	return nil
}
