package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadEnv(t *testing.T) {
	tests := []struct {
		name         string
		envContent   string
		expectedKeys map[string]CenvField
		expectError  bool
	}{
		{
			name:         "empty file",
			envContent:   "",
			expectedKeys: map[string]CenvField{},
			expectError:  false,
		},
		{
			name:       "simple",
			envContent: "FOO=bar",
			expectedKeys: map[string]CenvField{
				"FOO": {
					Key:            "FOO",
					value:          "bar",
					Required:       false,
					Public:         false,
					LengthRequired: false,
					Length:         0,
					Format:         "",
				},
			},
			expectError: false,
		},
		{
			name:       "multiple",
			envContent: "FOO=bar\nBAZ=qux",
			expectedKeys: map[string]CenvField{
				"FOO": {
					Key:            "FOO",
					value:          "bar",
					Required:       false,
					Public:         false,
					LengthRequired: false,
					Length:         0,
					Format:         "",
				},
				"BAZ": {
					Key:            "BAZ",
					value:          "qux",
					Required:       false,
					Public:         false,
					LengthRequired: false,
					Length:         0,
					Format:         "",
				},
			},
		},
		{
			name: "required",
			envContent: `# @required
FOO=bar`,
			expectedKeys: map[string]CenvField{
				"FOO": {
					Key:            "FOO",
					value:          "bar",
					Required:       true,
					Public:         false,
					LengthRequired: false,
					Length:         0,
					Format:         "",
				},
			},
		},
		{
			name: "public",
			envContent: `# @public
FOO=bar`,
			expectedKeys: map[string]CenvField{
				"FOO": {
					Key:            "FOO",
					value:          "bar",
					Required:       false,
					Public:         true,
					LengthRequired: false,
					Length:         0,
					Format:         "",
				},
			},
		},
		{
			name: "length",
			envContent: `# @length 3
FOO=bar`,
			expectedKeys: map[string]CenvField{
				"FOO": {
					Key:            "FOO",
					value:          "bar",
					Required:       false,
					Public:         false,
					LengthRequired: true,
					Length:         3,
					Format:         "",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tempDir := t.TempDir()
			envFilePath := filepath.Join(tempDir, ".env")
			if err := os.WriteFile(envFilePath, []byte(test.envContent), 0644); err != nil {
				t.Fatalf("failed to write temp .env file: %v", err)
			}

			gotEnv, err := ReadEnv(envFilePath)

			if (err != nil) != test.expectError {
				t.Errorf("unexpected error status: got %v, want error=%v", err, test.expectError)
			}

			if len(gotEnv) != len(test.expectedKeys) {
				t.Fatalf("unexpected number of keys: got %d, want %d", len(gotEnv), len(test.expectedKeys))
			}

			for key, expected := range test.expectedKeys {
				got, exists := gotEnv[key]
				if !exists {
					t.Errorf("missing key in parsed env: %s", key)
					continue
				}

				if got.Key != expected.Key || got.value != expected.value || got.Required != expected.Required || got.Public != expected.Public || got.LengthRequired != expected.LengthRequired || got.Length != expected.Length {
					t.Errorf("mismatch for key %s: got %+v, want %+v", key, got, expected)
				}
			}
		})
	}
}
