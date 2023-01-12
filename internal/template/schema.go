package template

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"

	"github.com/qri-io/jsonschema"
	"sigs.k8s.io/yaml"
)

func loadSchemaFromFile(fsys fs.FS, schemaFile string) (*jsonschema.Schema, error) {
	file, err := fsys.Open("schema.yaml")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	schemaBytes, err := io.ReadAll(file)
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

func validateInput(schema *jsonschema.Schema, input map[string]interface{}) error {
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
