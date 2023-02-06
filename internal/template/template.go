package template

import (
	"path/filepath"

	"github.com/cybercyst/go-scaffold/internal/download"
	"github.com/cybercyst/go-scaffold/internal/generate"
	"github.com/cybercyst/go-scaffold/internal/schema"
	"github.com/qri-io/jsonschema"
	"github.com/spf13/afero"
)

type Template struct {
	Uri       string             `json:"uri" yaml:"uri"`
	LocalPath string             `json:"localPath" yaml:"localPath"`
	Version   string             `json:"version" yaml:"version"`
	Config    *TemplateConfig    `json:"config" yaml:"config"`
	Schema    *jsonschema.Schema `json:"schema" yaml:"schema"`
}

func NewTemplate(uri string) (*Template, error) {
	downloadInfo, err := download.Download(uri)
	if err != nil {
		return nil, err
	}

	templateFs := afero.NewBasePathFs(afero.NewOsFs(), downloadInfo.LocalPath)
	config, err := LoadConfig(templateFs)
	if err != nil {
		return nil, err
	}

	schema, err := schema.LoadSchema(config.RawSchema)
	if err != nil {
		return nil, err
	}

	return &Template{
		Uri:       uri,
		LocalPath: downloadInfo.LocalPath,
		Version:   downloadInfo.Version,
		Config:    config,
		Schema:    schema,
	}, nil
}

func (t *Template) ValidateInput(input *map[string]interface{}) error {
	if err := schema.ValidateInput(t.Schema, input); err != nil {
		return err
	}

	return nil
}

func (t *Template) Execute(input *map[string]interface{}, outputPath string) ([]string, error) {
	templateFilesDir := filepath.Join(t.LocalPath, "template")
	templateFs := afero.NewBasePathFs(afero.NewOsFs(), templateFilesDir)
	outputFs := afero.NewBasePathFs(afero.NewOsFs(), outputPath)

	createdFiles, err := generate.GenerateTemplateFiles(templateFs, outputFs, input)
	if err != nil {
		return []string{}, err
	}

	return createdFiles, nil
}
