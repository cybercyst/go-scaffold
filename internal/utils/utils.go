package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cybercyst/go-scaffold/internal/consts"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

func ReadTemplateInput(inputFilePath string) (map[string]interface{}, error) {
	inputBytes, err := os.ReadFile(inputFilePath)
	if err != nil {
		return nil, err
	}

	inputJson := make(map[string]interface{})
	err = yaml.Unmarshal(inputBytes, &inputJson)
	if err != nil {
		return nil, err
	}

	return inputJson, nil
}

func EnsurePathExists(fs afero.Fs, path string) error {
	info, err := fs.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		err = fs.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
		return EnsurePathExists(fs, path)
	}

	if !info.IsDir() {
		return err
	}

	return nil
}

func CreateTempDir() string {
	path, err := os.MkdirTemp(os.TempDir(), fmt.Sprintf("%s-", consts.ProgramName))
	if err != nil {
		log.Fatal("Unable to create temporary directory: ", err)
	}

	return path
}

func CreateTestFs(fsys map[string]string) afero.Fs {
	fs := afero.NewMemMapFs()

	for path, contents := range fsys {
		fs.MkdirAll(path, 0755)
		afero.WriteFile(fs, path, []byte(strings.TrimSpace(contents)), 0644)
	}

	return fs
}
