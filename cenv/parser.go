package cenv

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jesperkha/gokenizer"
)

const prefix string = "@"

// ReadEnv parses the env file and tags. It also checks that
// tags are defined correctly and will return an error otherwise.
func ReadEnv(filepath string) (env map[string]CenvField, err error) {
	env = make(map[string]CenvField)

	file, err := os.ReadFile(filepath)
	if err != nil {
		return env, err
	}

	formats := []string{
		// Dont change order
		fmt.Sprintf("#{ws}%s{word} {number}", prefix),
		fmt.Sprintf("#{ws}%s{word} {text}", prefix),
		fmt.Sprintf("#{ws}%s{word}", prefix),
	}

	tokr := gokenizer.New()
	tokr.Class("key", "{var}")
	tokr.Class("value", "{string}", "{text}")
	tokr.Class("keyValue", "{ws}{key}{ws}={ws}{value}", "{ws}{key}{ws}={ws}")
	tokr.Class("comment", "#{any}")
	tokr.Class("tag", formats...)
	tokr.Class("expression", "{tag}", "{comment}", "{keyValue}")

	fld := CenvField{}

	tokr.Pattern("{expression}", func(t gokenizer.Token) error {
		tag := t.Get("expression").Get("tag")
		if tag.Length != 0 {
			name := tag.Get("word").Lexeme
			num := tag.Get("number").Lexeme
			format := tag.Get("text").Lexeme

			usesNum := false // Tag uses number value
			n := uint64(0)

			if num != "" {
				n, err = strconv.ParseUint(num, 10, 32)
				if err != nil {
					return fmt.Errorf("expected unsigned int, got '%s'", num)
				}
			}

			switch name {
			case "required":
				fld.Required = true

			case "public":
				fld.Public = true

			case "length":
				fld.LengthRequired = true
				fld.Length = uint32(n)
				usesNum = true

			case "format":
				fld.Format = format

			default:
				return fmt.Errorf("unknown tag '%s'", name)
			}

			if !usesNum && num != "" {
				return fmt.Errorf("expected newline after tag name, got '%s'", num)
			}
		}

		keyval := t.Get("expression").Get("keyValue")

		if keyval.Length != 0 {
			fld.Key = keyval.Get("key").Lexeme
			fld.value = keyval.Get("value").Lexeme

			// Strip string quotes
			if fld.value != "" && fld.value[0] == '"' && fld.value[len(fld.value)-1] == '"' {
				fld.value = fld.value[1 : len(fld.value)-1]
			}

			if fld.Public {
				fld.Value = fld.value
			}

			env[fld.Key] = fld
			fld = CenvField{}
		}

		return nil
	})

	errs := longError{}

	// Run for each line
	for _, line := range strings.Split(string(file), "\n") {
		if err = tokr.Run(line); err != nil {
			errs.Add(err.Error())
		}
	}

	return env, validateEnv(env, errs)
}

func validateEnv(fs map[string]CenvField, err longError) error {
	for _, f := range fs {
		if f.Format != "" {
			tokr := gokenizer.New()
			if ok, e := tokr.Matches(f.value, f.Format); e != nil || !ok {
				err.Add(fmt.Sprintf("value for '%s' did not match the format '%s'", f.Key, f.Format))
			}
		}
		if f.Required && len(f.value) == 0 {
			err.Add(fmt.Sprintf("required field '%s' must have a value", f.Key))
		}
		if f.Public && len(f.value) == 0 {
			err.Add(fmt.Sprintf("public field '%s' must have a value", f.Key))
		}
		if f.LengthRequired && int(f.Length) != len(f.value) {
			err.Add(fmt.Sprintf("tag expects length of '%s' to be %d, is %d", f.Key, f.Length, len(f.value)))
		}
	}

	return err.Error()
}
