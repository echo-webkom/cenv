package cenv

import (
	"os"
	"testing"
)

func TestValidate_RequiredField(t *testing.T) {
	schema := Schema{
		Entries: []Entry{{
			Key:      "FOO",
			Required: true,
		}},
	}
	os.Unsetenv("FOO")
	errs := validate(schema)
	if len(errs) != 1 || errs[0].Key != "FOO" || errs[0].Message != "required field is missing" {
		t.Errorf("expected missing required field error, got: %+v", errs)
	}

	os.Setenv("FOO", "")
	errs = validate(schema)
	if len(errs) != 1 || errs[0].Key != "FOO" || errs[0].Message != "required field is empty" {
		t.Errorf("expected empty required field error, got: %+v", errs)
	}

	os.Setenv("FOO", "bar")
	errs = validate(schema)
	if len(errs) != 0 {
		t.Errorf("expected no error, got: %+v", errs)
	}
	os.Unsetenv("FOO")
}

func TestValidate_RequiredLength(t *testing.T) {
	schema := Schema{
		Entries: []Entry{{
			Key:            "FOO",
			RequiredLength: ptrInt(3),
		}},
	}
	os.Setenv("FOO", "ab")
	errs := validate(schema)
	if len(errs) != 1 || errs[0].Key != "FOO" {
		t.Errorf("expected length error, got: %+v", errs)
	}
	os.Setenv("FOO", "abc")
	errs = validate(schema)
	if len(errs) != 0 {
		t.Errorf("expected no error, got: %+v", errs)
	}
	os.Unsetenv("FOO")
}

func TestValidate_LegalValues(t *testing.T) {
	schema := Schema{
		Entries: []Entry{{
			Key:         "FOO",
			LegalValues: []string{"a", "b"},
		}},
	}
	os.Setenv("FOO", "c")
	errs := validate(schema)
	if len(errs) != 1 || errs[0].Key != "FOO" {
		t.Errorf("expected legal values error, got: %+v", errs)
	}
	os.Setenv("FOO", "a")
	errs = validate(schema)
	if len(errs) != 0 {
		t.Errorf("expected no error, got: %+v", errs)
	}
	os.Unsetenv("FOO")
}

func TestValidate_RegexMatch(t *testing.T) {
	schema := Schema{
		Entries: []Entry{{
			Key:        "FOO",
			RegexMatch: ptrString(`^foo[0-9]+$`),
		}},
	}
	os.Setenv("FOO", "bar")
	errs := validate(schema)
	if len(errs) != 1 || errs[0].Key != "FOO" {
		t.Errorf("expected regex error, got: %+v", errs)
	}
	os.Setenv("FOO", "foo123")
	errs = validate(schema)
	if len(errs) != 0 {
		t.Errorf("expected no error, got: %+v", errs)
	}
	os.Unsetenv("FOO")
}

func TestValidate_Kind_Integer(t *testing.T) {
	schema := Schema{
		Entries: []Entry{{
			Key:  "FOO",
			Kind: &EntryKind{Type: "integer", MinInt: ptrInt64(10), MaxInt: ptrInt64(20)},
		}},
	}
	os.Setenv("FOO", "5")
	errs := validate(schema)
	if len(errs) != 1 || errs[0].Key != "FOO" {
		t.Errorf("expected min int error, got: %+v", errs)
	}
	os.Setenv("FOO", "25")
	errs = validate(schema)
	if len(errs) != 1 || errs[0].Key != "FOO" {
		t.Errorf("expected max int error, got: %+v", errs)
	}
	os.Setenv("FOO", "15")
	errs = validate(schema)
	if len(errs) != 0 {
		t.Errorf("expected no error, got: %+v", errs)
	}
	os.Unsetenv("FOO")
}

func TestValidate_Kind_Bool(t *testing.T) {
	schema := Schema{
		Entries: []Entry{{
			Key:  "FOO",
			Kind: &EntryKind{Type: "bool"},
		}},
	}
	os.Setenv("FOO", "maybe")
	errs := validate(schema)
	if len(errs) != 1 || errs[0].Key != "FOO" {
		t.Errorf("expected bool error, got: %+v", errs)
	}
	os.Setenv("FOO", "true")
	errs = validate(schema)
	if len(errs) != 0 {
		t.Errorf("expected no error, got: %+v", errs)
	}
	os.Unsetenv("FOO")
}

func ptrInt(i int) *int          { return &i }
func ptrInt64(i int64) *int64    { return &i }
func ptrString(s string) *string { return &s }
