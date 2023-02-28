package template

import (
	"errors"
	"fmt"

	"github.com/cybercyst/go-scaffold/internal/download"
	"github.com/cybercyst/go-scaffold/internal/schema"
	"github.com/qri-io/jsonschema"
	"github.com/spf13/afero"
)

type MetaTemplate struct {
	Templates []*Template        `json:"templates" yaml:"templates"`
	Schema    *jsonschema.Schema `json:"schema" yaml:"schema"`
}

type Template struct {
	Uri       string          `json:"uri" yaml:"uri"`
	LocalPath string          `json:"localPath" yaml:"localPath"`
	Version   string          `json:"version" yaml:"version"`
	Config    *TemplateConfig `json:"config" yaml:"config"`
	Steps     []Step          `json:"steps" yaml:"steps"`
}

func NewTemplate(fs afero.Fs, uri string) (*MetaTemplate, error) {
	template, err := downloadTemplate(fs, uri)
	if err != nil {
		return nil, err
	}

	deps := []*Template{
		template,
	}
	err = template.fetchDependencies(fs, &deps)
	if err != nil {
		return nil, err
	}

	schema, err := schema.LoadSchema(template.Config.RawSchema)

	return &MetaTemplate{
		Templates: deps,
		Schema:    schema,
	}, nil
}

func downloadTemplate(fs afero.Fs, uri string) (*Template, error) {
	downloadInfo, err := download.Download(fs, uri)
	if err != nil {
		return nil, err
	}

	templateFs := afero.NewBasePathFs(fs, downloadInfo.LocalPath)
	config, err := LoadConfig(templateFs)
	if err != nil {
		return nil, err
	}

	if _, err = schema.LoadSchema(config.RawSchema); err != nil {
		return nil, err
	}

	steps, stepErrors := mergeStepsFromConfigAndTemplateFilesystem(templateFs, config)
	if len(stepErrors) > 0 {
		err := errors.New("error loading steps")
		err = errors.Join(stepErrors...)
		return nil, err
	}

	return &Template{
		Uri:       uri,
		LocalPath: downloadInfo.LocalPath,
		Version:   downloadInfo.Version,
		Config:    config,
		Steps:     steps,
	}, nil
}

func (t *Template) fetchDependencies(fs afero.Fs, deps *[]*Template) error {
	remoteDependencies := []Step{}
	for _, step := range t.Steps {
		templateFs := afero.NewBasePathFs(fs, t.LocalPath)
		if isDependency(step, templateFs) {
			if err := checkCircularDependency(*deps, step); err != nil {
				return err
			}

			remoteDependencies = append(remoteDependencies, step)
		}
	}

	if len(remoteDependencies) == 0 {
		return nil
	}

	for _, step := range remoteDependencies {
		template, err := downloadTemplate(fs, step.Source)
		if err != nil {
			return err
		}

		*deps = append(*deps, template)

		if err := template.fetchDependencies(fs, deps); err != nil {
			return err
		}
	}

	return nil
}

func isDependency(step Step, templateFs afero.Fs) bool {
	// A template step that is not a directory within the
	// template's filesystem is a remote dependency and
	// needs to be fetched

	if step.Action != "template" {
		return false
	}

	isDir, _ := afero.DirExists(templateFs, step.Source)

	return !isDir
}

func checkCircularDependency(templates []*Template, step Step) error {
	for _, t := range templates {
		if t.Uri == step.Source {
			return errors.New(fmt.Sprintf("using template at %s would cause a circular dependency", step.Source))
		}
	}

	return nil
}

func (t *MetaTemplate) ValidateInput(input *map[string]interface{}) error {
	return schema.ValidateInput(t.Schema, input)
}
