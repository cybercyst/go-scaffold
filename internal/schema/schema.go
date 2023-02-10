package schema

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/barkimedes/go-deepcopy"
	"github.com/imdario/mergo"
	"github.com/qri-io/jsonschema"
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

func LoadSchema(schemaOrig interface{}) (*jsonschema.Schema, error) {
	var schema map[string]interface{}

	schemaRaw, err := deepcopy.Anything(schemaOrig)
	if err != nil {
		return nil, err
	}

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
		// If we can't detect the schema, we don't have one
		return nil, nil
	}

	schemaBytes, err := json.Marshal(schema)
	if err != nil {
		return nil, err
	}

	rs := &jsonschema.Schema{}
	err = json.Unmarshal(schemaBytes, rs)
	if err != nil {
		return nil, err
	}

	return rs, nil
}
