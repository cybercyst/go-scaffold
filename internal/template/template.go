package template

import "github.com/qri-io/jsonschema"

type Template struct {
	Uri       string
	LocalPath string
	Input     map[string]interface{}
	Version   string
	Schema    *jsonschema.Schema
}
