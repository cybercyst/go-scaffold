package scaffold

import (
	g "github.com/cybercyst/go-scaffold/internal/generate"
	t "github.com/cybercyst/go-scaffold/internal/template"
	"github.com/cybercyst/go-scaffold/internal/utils"
)

type Template = t.Template
type GeneratedMetadata = g.GeneratedMetadata

func Download(templateUri string) (*Template, error) {
	template, err := t.NewTemplate(templateUri)
	if err != nil {
		return nil, err
	}

	return template, nil
}

func Generate(template *Template, templateInput *map[string]interface{}, outputPath string) (*g.GeneratedMetadata, error) {
	if err := utils.EnsurePathExists(outputPath); err != nil {
		return nil, err
	}

	if err := template.ValidateInput(templateInput); err != nil {
		return nil, err
	}

	createdFiles, err := template.ExecuteSteps(templateInput, outputPath)
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
