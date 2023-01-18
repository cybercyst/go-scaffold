package schema

import (
	"testing"

	"github.com/spf13/afero"
)

func TestValidateSchemaShouldThrowNoErrorWhenInputMatchesSchema(t *testing.T) {
	fs := afero.NewMemMapFs()
	afero.WriteFile(fs, "schema.yaml", []byte(`
name: My Template
type: object
schema:
  project_name:
    type: string
required:
  - project_name
`), 0644)

	schema, err := LoadSchema(fs)
	if err != nil {
		t.Fatalf("got unexpected error while parsing schema: %s", err)
	}

	input := map[string]interface{}{
		"project_name": "my-project",
	}

	err = ValidateInput(schema, &input)
	if err != nil {
		t.Error("Got error when validating test input", err)
	}
}

func TestValidateSchemaShouldThrowErrorWhenInputDoesntMatchSchema(t *testing.T) {
	fs := afero.NewMemMapFs()
	afero.WriteFile(fs, "schema.yaml", []byte(`
name: My Template
type: object
schema:
  project_name:
    type: string
required:
  - project_name
`), 0644)

	schema, err := LoadSchema(fs)
	if err != nil {
		t.Error("got unexpected error while parsing schema")
	}

	input := map[string]interface{}{
		"invalid_key": "No One Cares About This Value",
	}

	err = ValidateInput(schema, &input)
	if err == nil {
		t.Error("Did not get an expected error when passing bad user input")
	}
}
