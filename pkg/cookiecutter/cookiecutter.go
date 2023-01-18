package cookiecutter

import (
	"github.com/cybercyst/go-cookiecutter/internal/generate"
	"github.com/cybercyst/go-cookiecutter/internal/template"
	"github.com/cybercyst/go-cookiecutter/internal/utils"
)

func Generate(templateUri string, templateInput *map[string]interface{}, outputPath string) (*generate.GeneratedMetadata, error) {
	if err := utils.EnsurePathExists(outputPath); err != nil {
		return nil, err
	}

	t, err := template.NewTemplate(templateUri)
	if err != nil {
		return nil, err
	}

	if err := t.ValidateInput(templateInput); err != nil {
		return nil, err
	}

	if err := t.Execute(templateInput, outputPath); err != nil {
		return nil, err
	}

	return &generate.GeneratedMetadata{
		Uri:     t.Uri,
		Version: t.Version,
		Input:   templateInput,
	}, nil
}
