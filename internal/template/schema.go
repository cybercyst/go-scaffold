package template

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/qri-io/jsonschema"
	"sigs.k8s.io/yaml"
)

func (t *Template) ValidateInput() error {
	schema, err := loadSchemaFromFile(filepath.Join(t.LocalPath, "schema.yaml"))
	if err != nil {
		return err
	}

	t.Schema = schema

	err = t.validateInput()
	if err != nil {
		return err
	}

	return nil
}

func (t *Template) validateInput() error {
	ctx := context.Background()

	userInputBytes, err := json.Marshal(t.Input)
	if err != nil {
		return err
	}

	validationErrors, err := t.Schema.ValidateBytes(ctx, userInputBytes)
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

func loadSchemaFromFile(schemaFile string) (*jsonschema.Schema, error) {
	schemaBytes, err := os.ReadFile(schemaFile)
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
