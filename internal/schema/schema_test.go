package schema

import (
	"testing"
)

var templateFile = "template.yaml"

func TestValidateSchemaShouldMergeSchemaEntries(t *testing.T) {
	rawSchema := []interface{}{
		map[string]interface{}{
			"required": []string{
				"project_name",
			},
			"properties": map[string]interface{}{
				"project_name": map[string]interface{}{
					"title":       "Project Name",
					"type":        "string",
					"description": "Give your project a name to dazzle",
				},
			},
		},
		map[string]interface{}{
			"required": []string{
				"owner",
			},
			"properties": map[string]interface{}{
				"owner": map[string]interface{}{
					"title":       "Owner",
					"type":        "string",
					"description": "The owner of this project. This will go in the CODEOWNERS file",
				},
			},
		},
	}

	schema, err := LoadSchema(rawSchema)
	if err != nil {
		t.Fatalf("got unexpected error while parsing schema: %s", err)
	}

	input := map[string]interface{}{
		"project_name": "my-project",
		"owner":        "test guy",
	}

	err = ValidateInput(schema, &input)
	if err != nil {
		t.Error("Got error when validating test input", err)
	}
}

func TestValidateSchemaShouldThrowNoErrorWhenInputMatchesSchema(t *testing.T) {
	rawSchema := map[string]interface{}{
		"required": []string{
			"project_name",
		},
		"properties": map[string]interface{}{
			"project_name": map[string]interface{}{
				"title":       "Project Name",
				"type":        "string",
				"description": "Give your project a name to dazzle",
			},
		},
	}

	schema, err := LoadSchema(rawSchema)
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
	rawSchema := map[string]interface{}{
		"required": []string{
			"project_name",
		},
		"properties": map[string]interface{}{
			"project_name": map[string]interface{}{
				"title":       "Project Name",
				"type":        "string",
				"description": "Give your project a name to dazzle",
			},
		},
	}

	schema, err := LoadSchema(rawSchema)
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
