package generate

import (
	"encoding/json"

	"github.com/extemporalgenome/slug"
	"github.com/flosch/pongo2/v6"
	"gopkg.in/yaml.v3"
)

func initFilters() {
	pongo2.RegisterFilter("slugify", filterSlugify)
	pongo2.RegisterFilter("yaml", filterYaml)
	pongo2.RegisterFilter("json", filterJSON)
}

func filterSlugify(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	return pongo2.AsValue(slug.Slug(in.String())), nil
}

func filterYaml(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	yamlBytes, err := yaml.Marshal(in.Interface())
	if err != nil {
		return nil, &pongo2.Error{
			Sender:    "filter:yaml",
			OrigError: err,
		}
	}

	return pongo2.AsValue(string(yamlBytes)), nil
}

func filterJSON(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	jsonBytes, err := json.MarshalIndent(in.Interface(), "", "  ")
	if err != nil {
		return nil, &pongo2.Error{
			Sender:    "filter:json",
			OrigError: err,
		}
	}

	return pongo2.AsSafeValue(string(jsonBytes)), nil
}
