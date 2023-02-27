package template

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/cybercyst/go-scaffold/internal/generate"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

type Step struct {
	Id          string                 `json:"id" yaml:"id"`
	Description string                 `json:"description" yaml:"description"`
	Action      string                 `json:"action" yaml:"action"`
	Target      string                 `json:"target" yaml:"target"`
	Command     []string               `json:"command" yaml:"command"`
	Source      string                 `json:"source" yaml:"source"`
	Inputs      map[string]interface{} `json:"inputs" yaml:"inputs"`
}

func mergeStepsFromConfigAndTemplateFilesystem(templateFs afero.Fs, config *TemplateConfig) ([]Step, []error) {
	validSteps := []Step{}
	allErrors := []error{}

	if config.Steps != nil {
		for idx, step := range config.Steps {
			if step.Id == "" {
				step.Id = fmt.Sprintf("Step %d", idx)
			}
			step = *parseStep(&step)

			if err := validateStep(step); err != nil {
				allErrors = append(allErrors, errors.New(fmt.Sprintf("step %+v is not valid: %s", step, err)))
			}

			validSteps = append(validSteps, step)
		}
	}

	stepsDirExists, _ := afero.IsDir(templateFs, "steps")
	if !stepsDirExists {
		return validSteps, allErrors
	}

	stepsFs := afero.NewBasePathFs(templateFs, "steps")
	files, err := afero.ReadDir(stepsFs, ".")
	if err != nil {
		allErrors = append(allErrors, errors.New(fmt.Sprintf("problem reading steps directory: %s", err)))
		return validSteps, allErrors
	}

	for _, file := range files {
		path := file.Name()
		step, err := loadStep(stepsFs, path)
		if err != nil {
			allErrors = append(allErrors, errors.New(fmt.Sprintf("step %s is not valid: %s", path, err)))
			continue
		}

		if step.Id == "" {
			base := filepath.Base(path)
			ext := filepath.Ext(path)
			step.Id = strings.TrimRight(base, ext)
		}
		step = parseStep(step)

		err = validateStep(*step)
		if err != nil {
			allErrors = append(allErrors, errors.New(fmt.Sprintf("step %s is not valid: %s", path, err)))
			continue
		}

		validSteps = append(validSteps, *step)
	}

	return validSteps, allErrors
}

func readStepFile(stepsFs afero.Fs, fname string, fileType FileType) (*Step, error) {
	var step Step

	bytes, err := afero.ReadFile(stepsFs, fname)
	if err != nil {
		return nil, err
	}

	switch fileType {
	case JSON:
		err = json.Unmarshal(bytes, &step)
		if err != nil {
			return nil, err
		}
	case YAML:
		err = yaml.Unmarshal(bytes, &step)
		if err != nil {
			return nil, err
		}
	}

	return &step, nil
}

func loadStep(stepsFs afero.Fs, stepPath string) (*Step, error) {
	var fileType FileType

	ext := filepath.Ext(stepPath)
	switch ext {
	case ".yaml":
		fileType = YAML
	case ".yml":
		fileType = YAML
	case ".json":
		fileType = JSON
	default:
		return nil, errors.New(fmt.Sprintf("unrecognized file extension %s", ext))
	}

	step, err := readStepFile(stepsFs, stepPath, fileType)
	if err != nil {
		return nil, err
	}

	err = validateStep(*step)
	if err != nil {
		return nil, err
	}

	return step, nil
}

func parseStep(step *Step) *Step {
	if step.Action == "" {
		step.Action = "template"
	}

	if step.Target == "" {
		step.Target = "."
	}

	return step
}

func validateStep(step Step) error {
	if step.Action == "template" {
		if step.Source == "" {
			return errors.New(fmt.Sprintf("required field source not set on step ID %s", step.Id))
		}
	}

	return nil
}

func (t *Template) ExecuteSteps(input *map[string]interface{}, fs afero.Fs, outputFs afero.Fs) ([]string, []error) {
	createdFiles := []string{}
	errs := []error{}
	for _, step := range t.Steps {
		templateFs := afero.NewBasePathFs(fs, t.LocalPath)
		if isDependency(step, templateFs) {
			continue
		}

		switch step.Action {
		case "template":
			newFiles, err := t.executeTemplateStep(step, templateFs, outputFs, input)
			if err != nil {
				errs = append(errs, err)
			}
			createdFiles = append(createdFiles, newFiles...)
		default:
			err := t.executeActionStep(step, templateFs)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	return createdFiles, errs
}

func (t *Template) executeTemplateStep(step Step, templateFs afero.Fs, outputFs afero.Fs, input *map[string]interface{}) ([]string, error) {
	sourceFs := afero.NewBasePathFs(templateFs, step.Source)
	var targetFs afero.Fs
	if step.Target != "." {
		targetFs = afero.NewBasePathFs(outputFs, step.Target)
	} else {
		targetFs = outputFs
	}
	return generate.GenerateTemplateFiles(sourceFs, targetFs, input)
}

func (t *Template) executeActionStep(step Step, targetFs afero.Fs) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}

	reader, err := cli.ImagePull(ctx, step.Action, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer reader.Close()
	io.Copy(os.Stdout, reader)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: step.Action,
		Tty:   false,
		Cmd:   step.Command,
	}, nil, nil, nil, "")
	if err != nil {
		return err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	return nil
}
