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
func ReadEnv(filepath string) (env []CenvField, err error) {
	file, err := os.ReadFile(filepath)
	if err != nil {
		return env, err
	}

	formats := []string{
		// Dont change order
		fmt.Sprintf("#{ws}%s{word} {number}", prefix),
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

			switch name {
			case "required":
				fld.Required = true
				if num != "" {
					return fmt.Errorf("expected newline after tag name, got '%s'", num)
				}

			case "length":
				fld.LengthRequired = true
				if n, err := strconv.ParseUint(num, 10, 32); err == nil {
					fld.Length = uint32(n)
				} else {
					return fmt.Errorf("expected unsigned int, got '%s'", num)
				}

			default:
				return fmt.Errorf("unknown tag '%s'", name)
			}
		}

		keyval := t.Get("expression").Get("keyValue")

		if keyval.Length != 0 {
			fld.Key = keyval.Get("key").Lexeme
			fld.value = keyval.Get("value").Lexeme

			env = append(env, fld)
			fld = CenvField{}
		}

		return nil
	})

	// Run for each line
	for _, line := range strings.Split(string(file), "\n") {
		if err = tokr.Run(line); err != nil {
			return env, err
		}
	}

	return env, validateEnv(env)
}

func validateEnv(fs []CenvField) error {
	for _, f := range fs {
		if f.Required && len(f.value) == 0 {
			return fmt.Errorf("required field '%s' must have a value", f.Key)
		}
		if f.LengthRequired && int(f.Length) != len(f.value) {
			return fmt.Errorf("tag expects length of '%s' to be %d, is %d", f.Key, f.Length, len(f.value))
		}
	}

	return nil
}
