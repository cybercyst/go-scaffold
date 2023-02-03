package generate

import (
	"encoding/json"
	"fmt"

	"github.com/extemporalgenome/slug"
	"github.com/flosch/pongo2/v6"
	"gopkg.in/yaml.v3"
)

func initFilters() {
	pongo2.RegisterFilter("slugify", filterSlugify)
	pongo2.RegisterFilter("yaml", filterYaml)
	pongo2.RegisterFilter("json", filterJson)
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

func filterJson(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	test := in.Interface()
	fmt.Println(test)

	jsonBytes, err := json.Marshal(in.Interface())
	if err != nil {
		return nil, &pongo2.Error{
			Sender:    "filter:json",
			OrigError: err,
		}
	}

	jsonStr := string(jsonBytes)

	return pongo2.AsValue(jsonStr), nil
}
