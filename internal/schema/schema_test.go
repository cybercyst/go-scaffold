package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateInput(t *testing.T) {
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
	assert.Nil(t, err)

	input := map[string]interface{}{
		"project_name": "my-project",
	}

	err = ValidateInput(schema, &input)
	assert.Nil(t, err)
}

func TestValidateInputShouldErrorWithMissingInput(t *testing.T) {
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
	assert.Nil(t, err)

	emptyInput := map[string]interface{}{}

	err = ValidateInput(schema, &emptyInput)
	assert.Error(t, err, "invalid user input provided\n\"project_name\" value is required")
}

func TestLoadSchema(t *testing.T) {
	singleSchema := map[string]interface{}{
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

	_, err := LoadSchema(singleSchema)
	assert.Nil(t, err)

	multipleSchemas := []interface{}{
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
	_, err = LoadSchema(multipleSchemas)
	assert.Nil(t, err)
}

func TestMergeSchemaShouldMergeSchemaEntries(t *testing.T) {
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

	got := map[string]interface{}{}
	err := Merge(&got, rawSchema)
	assert.Nil(t, err)

	want := map[string]interface{}{
		"required": []string{
			"project_name",
			"owner",
		},
		"properties": map[string]interface{}{
			"project_name": map[string]interface{}{
				"title":       "Project Name",
				"type":        "string",
				"description": "Give your project a name to dazzle",
			},
			"owner": map[string]interface{}{
				"title":       "Owner",
				"type":        "string",
				"description": "The owner of this project. This will go in the CODEOWNERS file",
			},
		},
	}

	assert.Equal(t, got, want)
}
