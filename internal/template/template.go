package template

import (
	"context"
	"errors"
	"io"
	"os"

	"github.com/cybercyst/go-scaffold/internal/download"
	"github.com/cybercyst/go-scaffold/internal/generate"
	"github.com/cybercyst/go-scaffold/internal/schema"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/qri-io/jsonschema"
	"github.com/spf13/afero"
)

type Template struct {
	Uri       string             `json:"uri" yaml:"uri"`
	LocalPath string             `json:"localPath" yaml:"localPath"`
	Version   string             `json:"version" yaml:"version"`
	Config    *TemplateConfig    `json:"config" yaml:"config"`
	Schema    *jsonschema.Schema `json:"schema" yaml:"schema"`
	Steps     []Step             `json:"steps" yaml:"steps"`
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

	steps, stepErrors := loadSteps(templateFs, config)
	// TODO: Treat errors as warnings here?
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
		Schema:    schema,
		Steps:     steps,
	}, nil
}

func (t *Template) ValidateInput(input *map[string]interface{}) error {
	if t.Schema == nil {
		// We have no schema defined, so we're assuming everything is A-OK
		return nil
	}

	if err := schema.ValidateInput(t.Schema, input); err != nil {
		return err
	}

	return nil
}

func (t *Template) ExecuteSteps(input *map[string]interface{}, outputPath string) ([]string, []error) {
	allCreatedFiles := []string{}
	allErrors := []error{}

	templateFs := afero.NewBasePathFs(afero.NewOsFs(), t.LocalPath)
	outputFs := afero.NewBasePathFs(afero.NewOsFs(), outputPath)

	for _, step := range t.Steps {
		var createdFiles []string
		var err error

		var targetFs afero.Fs
		if step.Target == "." {
			targetFs = outputFs
		} else {
			targetFs = afero.NewBasePathFs(outputFs, step.Target)
		}

		switch step.Action {
		case "template":
			sourceFs := afero.NewBasePathFs(templateFs, step.Source)
			isDir, _ := afero.IsDir(sourceFs, ".")

			if isDir {
				createdFiles, err = t.executeLocalTemplateStep(input, sourceFs, targetFs)
			} else {
				createdFiles, err = t.executeRemoteTemplateStep(input, step.Source, targetFs)
			}
		default:
			createdFiles, err = t.executeActionStep(step, targetFs)
		}

		if err != nil {
			allErrors = append(allErrors, err)
			continue
		}
		allCreatedFiles = append(allCreatedFiles, createdFiles...)
	}

	return allCreatedFiles, nil
}

func (t *Template) executeRemoteTemplateStep(input *map[string]interface{}, uri string, targetFs afero.Fs) ([]string, error) {
	nextTemplate, err := NewTemplate(uri)
	if err != nil {
		return nil, err
	}

	if err := nextTemplate.ValidateInput(input); err != nil {
		return nil, err
	}

	info, err := targetFs.Stat(".")
	if err != nil {
		return nil, err
	}
	createdFiles, stepErrors := nextTemplate.ExecuteSteps(input, info.Name())
	if len(stepErrors) > 0 {
		err := errors.New("problem running steps")
		err = errors.Join(stepErrors...)
		return createdFiles, err
	}

	return createdFiles, nil
}

func (t *Template) executeLocalTemplateStep(input *map[string]interface{}, sourceFs afero.Fs, targetFs afero.Fs) ([]string, error) {
	return generate.GenerateTemplateFiles(sourceFs, targetFs, input)
}

func (t *Template) executeActionStep(step Step, targetFs afero.Fs) ([]string, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	reader, err := cli.ImagePull(ctx, step.Action, types.ImagePullOptions{})
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	io.Copy(os.Stdout, reader)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: step.Action,
		Tty:   false,
		Cmd:   step.Command,
	}, nil, nil, nil, "")
	if err != nil {
		return nil, err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return nil, err
	}

	return nil, nil
}
