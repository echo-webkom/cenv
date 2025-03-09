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
	b, err := json.MarshalIndent(env, "", "\t")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath, b, 0666)
}

func ValidateSchema(env map[string]CenvField, schema CenvFile) error {
	errs := longError{}

	for _, f := range schema.Fields {
		if ff, ok := env[f.Key]; ok {
			errs.AddMany(validateField(f, ff))
		} else {
			errs.Add(fmt.Sprintf("missing field '%s'", f.Key))
		}
	}

	return errs.Error()
}

func assertBoolEqual(key, name string, schema, env bool) error {
	if schema && !env {
		return fmt.Errorf("'%s' is tagged with %s in schema, but not in env", key, name)
	}
	if !schema && env {
		return fmt.Errorf("'%s' is tagged with %s in env, but not in schema", key, name)
	}
	return nil
}

func validateField(sf CenvField, ef CenvField) (errs longError) {
	if err := assertBoolEqual(sf.Key, "required", sf.Required, ef.Required); err != nil {
		errs.Add(err.Error())
	}
	if err := assertBoolEqual(sf.Key, "public", sf.Public, ef.Public); err != nil {
		errs.Add(err.Error())
	}
	if err := assertBoolEqual(sf.Key, "a required length", sf.LengthRequired, ef.LengthRequired); err != nil {
		errs.Add(err.Error())
	}
	if sf.LengthRequired && ef.LengthRequired && sf.Length != ef.Length {
		errs.Add(fmt.Sprintf("'%s' is tagged with length %d in schema, but is %d in env", sf.Key, sf.Length, ef.Length))
	}
	if err := assertBoolEqual(sf.Key, "a required format", sf.Format != "", ef.Format != ""); err != nil {
		errs.Add(err.Error())
	}

	return errs
}
