package cookiecutter

import (
	g "github.com/cybercyst/go-cookiecutter/internal/generate"
	t "github.com/cybercyst/go-cookiecutter/internal/template"
	"github.com/cybercyst/go-cookiecutter/internal/utils"
)

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
