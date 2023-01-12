package template

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/flosch/pongo2/v6"
)

func Generate(path fs.FS, input map[string]interface{}) error {
	schema, err := loadSchemaFromFile(path, "schema.yaml")
	if err != nil {
		return err
	}

	err = validateInput(schema, input)
	if err != nil {
		return err
	}

	fmt.Println("Got here!")

	err = fs.WalkDir(path, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		finalPathName, err := executeTemplate([]byte(path), input)
		if err != nil {
			return err
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

		contents, err := executeTemplate(rawContents, input)
		if err != nil {
			return err
		}

		err = os.WriteFile(finalPathName, []byte(contents), 0644)
		if err != nil {
			return err
		}

		return nil
	})
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
