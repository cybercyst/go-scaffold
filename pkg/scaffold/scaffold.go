package scaffold

import (
	g "github.com/cybercyst/go-scaffold/internal/generate"
	s "github.com/cybercyst/go-scaffold/internal/schema"
	t "github.com/cybercyst/go-scaffold/internal/template"
	"github.com/cybercyst/go-scaffold/internal/utils"
	"github.com/spf13/afero"
)

func ReadSchemaBytes(templateUri string) ([]byte, error) {
	template, err := t.NewTemplate(templateUri)
	if err != nil {
		return nil, err
	}

	fs := afero.NewBasePathFs(afero.NewOsFs(), template.LocalPath)
	schemaBytes, err := s.ReadSchemaBytes(fs)
	if err != nil {
		return nil, err
	}

	return schemaBytes, nil
}

func Generate(templateUri string, templateInput *map[string]interface{}, outputPath string) (*g.GeneratedMetadata, error) {
	if err := utils.EnsurePathExists(outputPath); err != nil {
		return nil, err
	}

	template, err := t.NewTemplate(templateUri)
	if err != nil {
		return nil, err
	}

	if err := template.ValidateInput(templateInput); err != nil {
		return nil, err
	}

	createdFiles, err := template.Execute(templateInput, outputPath)
	if err != nil {
		return nil, err
	}

	return &g.GeneratedMetadata{
		Uri:          template.Uri,
		Version:      template.Version,
		Input:        templateInput,
		CreatedFiles: &createdFiles,
	}, nil
}
