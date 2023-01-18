package go_cookiecutter

import (
	"io/fs"
	"path/filepath"

	"github.com/flosch/pongo2/v6"
	"github.com/spf13/afero"
)

type GeneratedMetadata struct {
	Uri     string
	Version string
	Input   *map[string]interface{}
}

func Generate(templateUri string, templateInput *map[string]interface{}, outputPath string) (*GeneratedMetadata, error) {
	if err := ensurePathExists(outputPath); err != nil {
		return nil, err
	}

	template, err := NewTemplate(templateUri)
	if err != nil {
		return nil, err
	}

	if err := template.ValidateInput(templateInput); err != nil {
		return nil, err
	}

	if err := template.Execute(templateInput, outputPath); err != nil {
		return nil, err
	}

	return &GeneratedMetadata{
		Uri:     template.Uri,
		Version: template.Version,
		Input:   templateInput,
	}, nil
}

func generateTemplateFiles(templateFs afero.Fs, outputFs afero.Fs, input *map[string]interface{}) error {
	return afero.Walk(templateFs, ".", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		templatedPathName, err := executeTemplate([]byte(path), input)
		if err != nil {
			return err
		}

		if info.IsDir() {
			outputFs.MkdirAll(templatedPathName, 0755)
			return nil
		}

		err = outputFs.MkdirAll(filepath.Dir(templatedPathName), 0755)
		if err != nil {
			return err
		}

		rawContents, err := afero.ReadFile(templateFs, path)
		if err != nil {
			return err
		}

		contents, err := executeTemplate(rawContents, input)
		if err != nil {
			return err
		}

		err = afero.WriteFile(outputFs, templatedPathName, []byte(contents), 0644)
		if err != nil {
			return err
		}

		return nil
	})
}

func executeTemplate(templateBytes []byte, input *map[string]interface{}) (string, error) {
	template, err := pongo2.FromBytes(templateBytes)
	if err != nil {
		return "", err
	}

	output, err := template.Execute(*input)
	if err != nil {
		return "", err
	}

	return output, nil
}
