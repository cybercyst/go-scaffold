package generate

import (
	"io/fs"
	"path/filepath"

	"github.com/flosch/pongo2/v6"
	"github.com/spf13/afero"
)

type GeneratedMetadata struct {
	Uri          string                  `json:"uri"`
	Version      string                  `json:"version"`
	Input        *map[string]interface{} `json:"input"`
	CreatedFiles *[]string               `json:"-"`
}

func GenerateTemplateFiles(templateFs afero.Fs, outputFs afero.Fs, input *map[string]interface{}) ([]string, error) {
	createdFiles := []string{}
	err := afero.Walk(templateFs, ".", func(path string, info fs.FileInfo, err error) error {
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

		createdFiles = append(createdFiles, templatedPathName)

		return nil
	})

	return createdFiles, err
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
