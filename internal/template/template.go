package template

import (
	"path/filepath"

	"github.com/cybercyst/go-cookiecutter/internal/download"
	"github.com/cybercyst/go-cookiecutter/internal/generate"
	"github.com/cybercyst/go-cookiecutter/internal/schema"
	"github.com/qri-io/jsonschema"
	"github.com/spf13/afero"
)

type Template struct {
	Uri       string
	LocalPath string
	Version   string
	Schema    *jsonschema.Schema
}

func NewTemplate(uri string) (*Template, error) {
	downloadInfo, err := download.Download(uri)
	if err != nil {
		return nil, err
	}

	templateFs := afero.NewBasePathFs(afero.NewOsFs(), downloadInfo.LocalPath)
	schema, err := schema.LoadSchema(templateFs)
	if err != nil {
		return nil, err
	}

	return &Template{
		Uri:       uri,
		LocalPath: downloadInfo.LocalPath,
		Version:   downloadInfo.Version,
		Schema:    schema,
	}, nil
}

func (t *Template) ValidateInput(input *map[string]interface{}) error {
	if err := schema.ValidateInput(t.Schema, input); err != nil {
		return err
	}

	return nil
}

func (t *Template) Execute(input *map[string]interface{}, outputPath string) error {
	templateFilesDir := filepath.Join(t.LocalPath, "template")
	templateFs := afero.NewBasePathFs(afero.NewOsFs(), templateFilesDir)
	outputFs := afero.NewBasePathFs(afero.NewOsFs(), outputPath)
	err := generate.GenerateTemplateFiles(templateFs, outputFs, input)
	if err != nil {
		return err
	}

	return nil
}
