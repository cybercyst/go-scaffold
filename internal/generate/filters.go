package generate

import (
	"github.com/extemporalgenome/slug"
	"github.com/flosch/pongo2/v6"
)

func initFilters() {
	pongo2.RegisterFilter("slugify", filterSlugify)
}

func filterSlugify(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	return pongo2.AsValue(slug.Slug(in.String())), nil
}
