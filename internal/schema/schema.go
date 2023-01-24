package schema

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/imdario/mergo"
	"github.com/qri-io/jsonschema"
	"github.com/spf13/afero"
	"sigs.k8s.io/yaml"
)

func ValidateInput(schema *jsonschema.Schema, input *map[string]interface{}) error {
	ctx := context.Background()

	userInputBytes, err := json.Marshal(input)
	if err != nil {
		return err
	}

	validationErrors, err := schema.ValidateBytes(ctx, userInputBytes)
	if err != nil {
		return err
	}

	if len(validationErrors) > 0 {
		fmt.Println("The following validation errors were discovered while attempting to generate this template:")
		for _, validationError := range validationErrors {
			fmt.Println(validationError.Error())
		}
		return fmt.Errorf("the provided user input did not pass this template's schema")
	}

	return nil
}

func LoadSchema(fs afero.Fs) (*jsonschema.Schema, error) {
	schemaBytes, err := afero.ReadFile(fs, "schema.yaml")
	if err != nil {
		return nil, err
	}

	var schemaRaw interface{}
	if err := yaml.Unmarshal(schemaBytes, &schemaRaw); err != nil {
		return nil, err
	}

	var schema map[string]interface{}

	switch schemaRaw := schemaRaw.(type) {
	case []interface{}:
		for _, schemaSection := range schemaRaw {
			if err := mergo.Merge(&schema, schemaSection, mergo.WithAppendSlice); err != nil {
				return nil, err
			}
		}
	case map[string]interface{}:
		schema = schemaRaw
	default:
		return nil, fmt.Errorf("got unexpected value from schema.yaml")
	}

	schemaJsonBytes, err := json.Marshal(schema)
	if err != nil {
		return nil, err
	}

	rs := &jsonschema.Schema{}
	err = json.Unmarshal(schemaJsonBytes, rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}

func ReadSchema(templatePath string) ([]byte, error) {
	fs := afero.NewBasePathFs(afero.NewOsFs(), templatePath)

	schemaYamlBytes, err := afero.ReadFile(fs, "schema.yaml")
	if err != nil {
		return nil, err
	}

	schemaJsonBytes, err := yaml.YAMLToJSON(schemaYamlBytes)
	if err != nil {
		return nil, err
	}

	return schemaJsonBytes, nil
}
