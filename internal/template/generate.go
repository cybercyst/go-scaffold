package template

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/flosch/pongo2/v6"
)

func (t *Template) Generate() error {
	templateFilesDir := filepath.Join(t.LocalPath, "template")
	err := filepath.WalkDir(templateFilesDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		templatedPathName, err := executeTemplate([]byte(path), t.Input)
		if err != nil {
			return err
		}

		finalPathName, err := filepath.Rel(templateFilesDir, templatedPathName)
		if err != nil {
			return err
		}

		if finalPathName == "" {
			return nil
		}

		if d.IsDir() {
			os.MkdirAll(finalPathName, 0755)
			return nil
		}

		err = os.MkdirAll(filepath.Dir(finalPathName), 0755)
		if err != nil {
			return err
		}

		rawContents, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		contents, err := executeTemplate(rawContents, t.Input)
		if err != nil {
			return err
		}

		err = os.WriteFile(finalPathName, []byte(contents), 0644)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func executeTemplate(templateBytes []byte, input map[string]interface{}) (string, error) {
	template, err := pongo2.FromBytes(templateBytes)
	if err != nil {
		return "", err
	}

	output, err := template.Execute(input)
	if err != nil {
		return "", err
	}

	return output, nil
}
