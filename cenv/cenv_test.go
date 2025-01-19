package cenv_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/echo-webkom/cenv/cenv"
)

func TestUpdate(t *testing.T) {
	tests := []struct {
		name          string
		envContent    string
		schemaContent string
		expectError   bool
	}{
		{
			name:          "simple",
			envContent:    `FOO=bar`,
			schemaContent: `{"lastUpdated":"1970-01-01T00:00:00.000000+00:00","fields":{"FOO":{"required":true,"public":false,"lengthRequired":false,"length":0,"format":"","key":"FOO","value":""}}}`,
			expectError:   false,
		},
		{
			name:          "missing field",
			envContent:    `FOO=`,
			schemaContent: `{"lastUpdated":"1970-01-01T00:00:00.000000+00:00","fields":{}}`,
			expectError:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tempDir := t.TempDir()

			envPath := filepath.Join(tempDir, ".env")
			schemaPath := filepath.Join(tempDir, "cenv.schema.json")

			if err := os.WriteFile(envPath, []byte(test.envContent), 0644); err != nil {
				t.Fatalf("failed to write .env file: %v", err)
			}

			if err := os.WriteFile(schemaPath, []byte(test.schemaContent), 0644); err != nil {
				t.Fatalf("failed to write cenv.schema.json file: %v", err)
			}

			err := cenv.Update(envPath, schemaPath)
			if (err != nil) && !test.expectError {
				t.Fatalf("Update failed: %v", err)
			}

			data, err := os.ReadFile(schemaPath)
			if err != nil {
				t.Fatalf("failed to read schema file: %v", err)
			}

			if !strings.Contains(string(data), `"FOO":{`) {
				t.Errorf("schema file is missing FOO")
			}

			if !strings.Contains(string(data), `"lastUpdated"`) {
				t.Errorf("schema file is missing the lastUpdated field")
			}
		})
	}
}

func TestCheck(t *testing.T) {
	tests := []struct {
		name          string
		envContent    string
		schemaContent string
		expectError   bool
	}{
		{
			name: "simple",
			envContent: `# @required
FOO=bar`,
			schemaContent: `{"lastUpdated":"1970-01-01T00:00:00.000000+00:00","fields":{"FOO":{"required":true,"public":false,"lengthRequired":false,"length":0,"format":"","key":"FOO","value":""}}}`,
			expectError:   false,
		},
		{
			name:          "missing field",
			envContent:    `BAR=foo`,
			schemaContent: `{"lastUpdated":"1970-01-01T00:00:00.000000+00:00","fields":{"FOO":{"required":true,"public":false,"lengthRequired":false,"length":0,"format":"","key":"FOO","value":""}}}`,
			expectError:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tempDir := t.TempDir()

			envPath := filepath.Join(tempDir, ".env")
			schemaPath := filepath.Join(tempDir, "cenv.schema.json")

			if err := os.WriteFile(envPath, []byte(test.envContent), 0644); err != nil {
				t.Fatalf("failed to write .env file: %v", err)
			}

			if err := os.WriteFile(schemaPath, []byte(test.schemaContent), 0644); err != nil {
				t.Fatalf("failed to write cenv.schema.json file: %v", err)
			}

			err := cenv.Check(envPath, schemaPath)
			if (err != nil) && !test.expectError {
				t.Fatalf("Check failed: %v", err)
			}
		})
	}
}

func TestFix(t *testing.T) {
	tests := []struct {
		name          string
		envContent    string
		schemaContent string
	}{
		{
			name:          "simple",
			envContent:    `FOO=bar`,
			schemaContent: `{"lastUpdated":"1970-01-01T00:00:00.000000+00:00","fields":{"FOO":{"required":true,"public":false,"lengthRequired":false,"length":0,"format":"","key":"FOO","value":""}}}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tempDir := t.TempDir()

			envPath := filepath.Join(tempDir, ".env")
			schemaPath := filepath.Join(tempDir, "cenv.schema.json")

			if err := os.WriteFile(envPath, []byte(test.envContent), 0644); err != nil {
				t.Fatalf("failed to write .env file: %v", err)
			}

			if err := os.WriteFile(schemaPath, []byte(test.schemaContent), 0644); err != nil {
				t.Fatalf("failed to write cenv.schema.json file: %v", err)
			}

			if err := cenv.Fix(envPath, schemaPath); err != nil {
				t.Fatalf("Fix failed: %v", err)
			}

			data, err := os.ReadFile(envPath)
			if err != nil {
				t.Fatalf("failed to read .env file: %v", err)
			}

			if !strings.Contains(string(data), "FOO=") {
				t.Errorf("env file is missing FOO")
			}
		})
	}
}
