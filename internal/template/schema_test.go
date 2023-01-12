package template

import (
	"testing"
	"testing/fstest"
)

func TestValidateSchemaShouldThrowNoErrorWhenInputMatchesSchema(t *testing.T) {
	m := fstest.MapFS{
		"schema.yaml": {
			Data: []byte(`
name: My Template
type: object
schema:
  project_name:
    type: string
required:
  - project_name
`),
		},
	}

	schema, err := loadSchemaFromFile(m, "schema.yaml")
	if err != nil {
		t.Error("got unexpected error while parsing schema")
	}

	input := map[string]interface{}{
		"project_name": "my-project",
	}

	err = validateInput(schema, input)
	if err != nil {
		t.Error("Got error when validating test input", err)
	}
}

func TestValidateSchemaShouldThrowErrorWhenInputDoesntMatchSchema(t *testing.T) {
	m := fstest.MapFS{
		"schema.yaml": {
			Data: []byte(`
name: My Template
type: object
schema:
  project_name:
    type: string
required:
  - project_name
`),
		},
	}

	schema, err := loadSchemaFromFile(m, "schema.yaml")
	if err != nil {
		t.Error("got unexpected error while parsing schema")
	}

	input := map[string]interface{}{
		"invalid_key": "No One Cares About This Value",
	}

	err = validateInput(schema, input)
	if err == nil {
		t.Error("Did not get an expected error when passing bad user input")
	}
}
