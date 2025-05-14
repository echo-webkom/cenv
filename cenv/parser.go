package cenv

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/jesperkha/gokenizer"
)

const PREFIX byte = '@'

// ReadEnv parses the env file and tags. It also checks that
// tags are defined correctly and will return an error otherwise.
func ReadEnv(filepath string) (env map[string]CenvField, err error) {
	file, err := os.ReadFile(filepath)
	if err != nil {
		return env, err
	}

	env = make(map[string]CenvField)
	fld := CenvField{}
	errs := longError{}

	for idx, line := range strings.Split(string(file), "\n") {
		line := strings.TrimSpace(line)
		if line == "" {
			continue
		}

		linenr := idx + 1

		// Parse tag
		// Tags are comment with any amount of whitespace followed by
		// the prefix '@' then the tag name, and finally the tag value if any.
		if strings.HasPrefix(line, "#") {
			tagLine := strings.TrimSpace(line[1:])
			if len(tagLine) == 0 || tagLine[0] != PREFIX {
				continue
			}

			split := strings.SplitN(tagLine, " ", 2)
			tag := strings.TrimPrefix(split[0], string(PREFIX))
			value := ""
			if len(split) > 1 {
				value = split[1]
			}

			switch tag {
			case "required": // @required
				fld.Required = true
			case "public": // @public
				fld.Public = true
			case "format": // @format <any>
				fld.Format = value

			case "enum": // @enum <name> | <name> | ...
				values := strings.Split(value, "|")
				for _, v := range values {
					v = strings.TrimSpace(v)
					if v != "" {
						fld.Enum = append(fld.Enum, v)
					}
				}

			case "length": // @length <number>
				n, err := strconv.Atoi(value)
				if err != nil {
					errs.Add(fmt.Errorf("invalid number literal '%s', line %d", value, linenr))
				}
				fld.Length = uint32(n)
				fld.LengthRequired = true

			default:
				errs.Add(fmt.Errorf("unknown tag '%s', line %d", tag, linenr))
			}

			continue
		}

		// Parse key-value pair
		// For the sake of consistency across code/config/.env we require that
		// keys are alphanumeric (in addition to $ and _).
		split := strings.SplitN(line, "=", 2)
		if len(split) == 0 {
			errs.Add(fmt.Errorf("syntax error, line %d", linenr))
		}

		key := strings.TrimSpace(split[0])
		value := ""
		if len(split) > 1 {
			value = strings.Trim(split[1], "\" ")
		}

		if !isAlnum(key) {
			errs.Add(fmt.Errorf("key must be alphanumeric, line %d", linenr))
		}

		fld.Key = key
		fld.value = value

		if fld.Public {
			fld.Value = value
		}

		env[fld.Key] = fld
		fld = CenvField{} // Reset
	}

	return env, validateEnv(env, errs)
}

func isAlnum(s string) bool {
	for _, c := range s {
		if strings.IndexRune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_$", c) == -1 {
			return false
		}
	}
	return true
}

// Validate the env file in-place; check that all tags actually match the
// values before comparing with schema.
func validateEnv(fs map[string]CenvField, longErr longError) error {
	for _, f := range fs {
		if err := validateEnvField(f); err != nil {
			longErr.Add(err)
		}
	}

	return longErr.Error()
}

func validateEnvField(field CenvField) error {
	if field.Format != "" {
		tokr := gokenizer.New()
		if ok, e := tokr.Matches(field.value, field.Format); e != nil || !ok {
			return fmt.Errorf("'%s': value did not match the format '%s'", field.Key, field.Format)
		}
	}
	if len(field.Enum) != 0 {
		if !slices.Contains(field.Enum, field.value) {
			list := strings.Join(field.Enum, ", ")
			return fmt.Errorf("'%s': field must have one of the specified enum values: [%s]", field.Key, list)
		}
	}
	if field.Required && len(field.value) == 0 {
		return fmt.Errorf("'%s': required field must have a value", field.Key)
	}
	if field.Public && len(field.value) == 0 {
		return fmt.Errorf("'%s': public field must have a value", field.Key)
	}
	if field.LengthRequired && int(field.Length) != len(field.value) {
		return fmt.Errorf("'%s': tag expects length to be %d, is %d", field.Key, field.Length, len(field.value))
	}
	return nil
}
