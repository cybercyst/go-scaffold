package scaffold

import (
	"errors"

	g "github.com/cybercyst/go-scaffold/internal/generate"
	t "github.com/cybercyst/go-scaffold/internal/template"
	"github.com/cybercyst/go-scaffold/internal/utils"
	"github.com/spf13/afero"
)

type (
	Template          = t.MetaTemplate
	GeneratedMetadata = g.GeneratedMetadata
)

func Download(templateURI string) (*Template, error) {
	fs := afero.NewOsFs()
	template, err := t.NewTemplate(fs, templateURI)
	if err != nil {
		return nil, err
	}

	return template, nil
}

func Generate(meta *Template, input *map[string]interface{}, outputPath string) (*g.GeneratedMetadata, error) {
	fs := afero.NewOsFs()
	if err := utils.EnsurePathExists(fs, outputPath); err != nil {
		return nil, err
	}
	outputFs := afero.NewBasePathFs(fs, outputPath)

	if err := meta.ValidateInput(input); err != nil {
		return nil, err
	}

	createdFiles := []string{}
	templateMetadata := []g.TemplateMetadata{}

	for _, template := range meta.Templates {
		stepCreatedFiles, stepErrors := template.ExecuteSteps(input, fs, outputFs)
		if len(stepErrors) > 0 {
			// err := errors.New("problem running steps")
			err := errors.Join(stepErrors...)
			return nil, err
		}
		createdFiles = append(createdFiles, stepCreatedFiles...)
		templateMetadata = append(templateMetadata, g.TemplateMetadata{
			URI:     template.URI,
			Version: template.Version,
		})
	}

	return &g.GeneratedMetadata{
		Templates:    &templateMetadata,
		Input:        input,
		CreatedFiles: &createdFiles,
	}, nil
}
