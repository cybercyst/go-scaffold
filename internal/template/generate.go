package template

import (
	"io/fs"
	"path/filepath"

	"github.com/flosch/pongo2/v6"
	"github.com/spf13/afero"
)

func (t *Template) Generate() error {
	templateFilesDir := filepath.Join(t.LocalPath, "template")
	templateFs := afero.NewBasePathFs(afero.NewOsFs(), templateFilesDir)
	outputFs := afero.NewBasePathFs(afero.NewOsFs(), t.OutputPath)
	err := generateTemplateFiles(templateFs, outputFs, &t.Input)
	if err != nil {
		return err
	}

	return nil
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
