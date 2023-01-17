package template

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/qri-io/jsonschema"
	"github.com/spf13/afero"
	"sigs.k8s.io/yaml"
)

func (t *Template) ValidateInput() error {
	fs := afero.NewBasePathFs(afero.NewOsFs(), t.LocalPath)
	schema, err := loadSchemaFromFile(fs, "schema.yaml")
	if err != nil {
		return err
	}

	t.Schema = schema

	err = validateInput(t.Schema, &t.Input)
	if err != nil {
		return err
	}

	return nil
}

func validateInput(schema *jsonschema.Schema, input *map[string]interface{}) error {
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

func loadSchemaFromFile(fs afero.Fs, schemaFile string) (*jsonschema.Schema, error) {
	schemaBytes, err := afero.ReadFile(fs, schemaFile)
	if err != nil {
		return nil, err
	}

	schemaJson, err := yaml.YAMLToJSON(schemaBytes)
	if err != nil {
		return nil, err
	}

	rs := &jsonschema.Schema{}
	err = json.Unmarshal(schemaJson, rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}
