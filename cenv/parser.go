package cenv

import (
	"fmt"
	"os"
	"strings"
)

const prefix string = "@"

type field struct {
	isEmpty bool
	isTag   bool // Sets value to the tag name

	key   string
	value string
	err   error
}

func applyTag(fld *CenvField, tag string) error {
	if tag == "required" {
		fld.Required = true
	} else {
		return fmt.Errorf("unknown tag '%s'", tag)
	}

	return nil
}

func parseLine(line string) (f field, err error) {
	trimmed := strings.Trim(line, " ")
	if len(trimmed) == 0 {
		return field{isEmpty: true}, err
	}

	if trimmed[0] == '#' {
		t := strings.Trim(trimmed[1:], " ")

		if len(t) > 1 && string(t[0]) == prefix {
			return field{isTag: true, value: t[1:]}, err
		}

		return field{isEmpty: true}, err
	}

	split := strings.Split(trimmed, "=")
	if len(split) != 2 || split[0] == "" {
		return f, fmt.Errorf("syntax error: %s", line)
	}

	key := strings.Trim(split[0], " ")
	value := strings.ReplaceAll(strings.Trim(split[1], " "), "\"", "")
	return field{key: key, value: value, err: nil}, err
}

func ReadEnv(filepath string) (env CenvFile, err error) {
	file, err := os.ReadFile(filepath)
	if err != nil {
		return env, err
	}

	fld := CenvField{}
	for idx, line := range strings.Split(string(file), "\n") {
		f, err := parseLine(line)
		if err != nil {
			return env, fmt.Errorf("%s, line %d", err.Error(), idx+1)
		}

		if f.isEmpty {
			fld = CenvField{}
			continue
		}

		if f.isTag {
			if err = applyTag(&fld, f.value); err != nil {
				return env, fmt.Errorf("%s, line %d", err.Error(), idx+1)
			}
			continue
		}

		fld.Key = f.key
		fld.value = f.value
		env = append(env, fld)

		fld = CenvField{}
	}

	return env, nil
}
