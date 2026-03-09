package cenv

import (
	"errors"
	"fmt"
	"net"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/joho/godotenv"
)

type Config struct {
	EnvPath    string
	SchemaPath string
}

var defaultConfig = Config{
	EnvPath:    "./.env",
	SchemaPath: "./cenv.schema.toml",
}

// Check verifies that the environment variables match the cenv schema. [config] may be nil.
//
// Check tries to load a .env file using [config], or from the default path if not provided.
// If a .env file does not exist, the process env is used, and no error is returned.
//
// Check returns an error if the schema file is missing or malformed.
//
//	// Default config
//	Config{
//		EnvPath: "./.env",
//		SchemaPath: "./cenv.schema.toml",
//	}
func Check(config *Config) error {
	if config == nil {
		config = &defaultConfig
	}

	_ = godotenv.Load(config.EnvPath)

	schema, err := readSchema(config.SchemaPath)
	if err != nil {
		return fmt.Errorf("cenv: failed to read schema: %v", err)
	}

	if errs := validate(*schema); len(errs) > 0 {
		var err error
		for _, e := range errs {
			err = errors.Join(err, e)
		}
		return err
	}

	return nil
}

func readSchema(path string) (*Schema, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var schema Schema
	if err := toml.Unmarshal(data, &schema); err != nil {
		return nil, err
	}

	return &schema, nil
}

type validationError struct {
	Key     string
	Message string
	Hint    *string
}

func (e validationError) Error() string {
	if e.Hint != nil {
		return fmt.Sprintf("%s: %s\n\thint: %s", e.Key, e.Message, *e.Hint)
	}
	return fmt.Sprintf("%s: %s", e.Key, e.Message)
}

func validate(schema Schema) []validationError {
	var errors []validationError

	for _, entry := range schema.Entries {
		value, exists := os.LookupEnv(entry.Key)
		hint := entry.Hint

		// Check required
		if entry.Required {
			if !exists {
				errors = append(errors, validationError{
					Key:     entry.Key,
					Message: "required field is missing",
					Hint:    hint,
				})
				continue
			}
			if value == "" {
				errors = append(errors, validationError{
					Key:     entry.Key,
					Message: "required field is empty",
					Hint:    hint,
				})
				continue
			}
		}

		// Skip further validation if value is missing or empty (and not required)
		if !exists || value == "" {
			continue
		}

		// Check required_length
		if entry.RequiredLength != nil {
			requiredLen := *entry.RequiredLength
			if len(value) != requiredLen {
				errors = append(errors, validationError{
					Key:     entry.Key,
					Message: fmt.Sprintf("expected length %d, got %d", requiredLen, len(value)),
					Hint:    hint,
				})
			}
		}

		// Check legal_values
		if len(entry.LegalValues) > 0 {
			if !slices.Contains(entry.LegalValues, value) {
				errors = append(errors, validationError{
					Key:     entry.Key,
					Message: fmt.Sprintf("value '%s' is not one of: %s", value, strings.Join(entry.LegalValues, ", ")),
					Hint:    hint,
				})
			}
		}

		// Check regex_match
		if entry.RegexMatch != nil {
			re, err := regexp.Compile(*entry.RegexMatch)
			if err != nil {
				errors = append(errors, validationError{
					Key:     entry.Key,
					Message: fmt.Sprintf("invalid regex pattern: %s", err),
					Hint:    hint,
				})
			} else if !re.MatchString(value) {
				errors = append(errors, validationError{
					Key:     entry.Key,
					Message: fmt.Sprintf("value does not match pattern: %s", *entry.RegexMatch),
					Hint:    hint,
				})
			}
		}

		// Check kind-specific validation
		if entry.Kind != nil {
			if errMsg := validateKind(value, entry.Kind); errMsg != "" {
				errors = append(errors, validationError{
					Key:     entry.Key,
					Message: errMsg,
					Hint:    hint,
				})
			}
		}
	}

	return errors
}

func validateKind(value string, kind *EntryKind) string {
	kindType := strings.ToLower(strings.ReplaceAll(kind.Type, "_", ""))

	switch kindType {
	case "integer":
		n, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Sprintf("'%s' is not a valid integer", value)
		}
		if kind.MinInt != nil && n < *kind.MinInt {
			return fmt.Sprintf("value %d is less than minimum %d", n, *kind.MinInt)
		}
		if kind.MaxInt != nil && n > *kind.MaxInt {
			return fmt.Sprintf("value %d is greater than maximum %d", n, *kind.MaxInt)
		}

	case "float":
		n, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Sprintf("'%s' is not a valid float", value)
		}
		if kind.MinFloat != nil && n < *kind.MinFloat {
			return fmt.Sprintf("value %f is less than minimum %f", n, *kind.MinFloat)
		}
		if kind.MaxFloat != nil && n > *kind.MaxFloat {
			return fmt.Sprintf("value %f is greater than maximum %f", n, *kind.MaxFloat)
		}

	case "string":
		// No validation needed for string type
		return ""

	case "url":
		// URL validation using RFC 3986 compliant regex
		urlRegex := regexp.MustCompile(`(?i)^[a-z][a-z0-9+.-]*://(([a-z0-9._~%!$&'()*+,;=:-]*@)?([a-z0-9.-]+|\[[a-f0-9:]+\])(:[0-9]+)?)?(/[a-z0-9._~%!$&'()*+,;=:@/-]*)?(\?[a-z0-9._~%!$&'()*+,;=:@/?-]*)?(\#[a-z0-9._~%!$&'()*+,;=:@/?-]*)?$`)
		if !urlRegex.MatchString(value) {
			return fmt.Sprintf("'%s' is not a valid URL", value)
		}

	case "email":
		// Email validation using RFC 5321/5322 compliant regex
		emailRegex := regexp.MustCompile(`(?i)^[a-z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+@[a-z0-9]([a-z0-9-]*[a-z0-9])?(\.[a-z0-9]([a-z0-9-]*[a-z0-9])?)+$`)
		if !emailRegex.MatchString(value) {
			return fmt.Sprintf("'%s' is not a valid email address", value)
		}

	case "bool":
		lower := strings.ToLower(value)
		validBools := []string{"true", "false", "1", "0", "yes", "no"}
		if !slices.Contains(validBools, lower) {
			return fmt.Sprintf("'%s' is not a valid boolean", value)
		}

	case "ipaddress":
		if net.ParseIP(value) == nil {
			return fmt.Sprintf("'%s' is not a valid IP address", value)
		}

	case "path":
		if value == "" {
			return "path cannot be empty"
		}
		// Path validation - validates Unix and Windows path formats
		unixPathRegex := regexp.MustCompile(`^(/|\.{1,2}/)?([a-zA-Z0-9._-]+/?)*[a-zA-Z0-9._-]*$`)
		windowsPathRegex := regexp.MustCompile(`(?i)^([a-z]:[/\\]|\\\\[a-z0-9._-]+[/\\][a-z0-9._-]+)?([a-zA-Z0-9._-]+[/\\]?)*[a-zA-Z0-9._-]*$`)
		if !unixPathRegex.MatchString(value) && !windowsPathRegex.MatchString(value) {
			return fmt.Sprintf("'%s' is not a valid path", value)
		}
	}

	return ""
}
