package schema

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/barkimedes/go-deepcopy"
	"github.com/imdario/mergo"
	"github.com/qri-io/jsonschema"
)

func ValidateInput(schema *jsonschema.Schema, input *map[string]interface{}) error {
	// no schema defined, we accept any input
	if schema == nil {
		return nil
	}

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
		err := errors.New("invalid user input provided")
		for _, validationError := range validationErrors {
			err = errors.Join(err, errors.New(validationError.Message))
		}
		return err
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
		if err = Merge(&schema, schemaRaw); err != nil {
			return nil, err
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

func Merge(dst *map[string]interface{}, schemas []interface{}) error {
	for _, schema := range schemas {
		if err := mergo.Merge(dst, schema, mergo.WithAppendSlice); err != nil {
			return err
		}
	}

	return nil
}
