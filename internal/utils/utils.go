package utils

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/cybercyst/go-scaffold/internal/consts"
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

func EnsurePathExists(path string) error {
	info, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
		return EnsurePathExists(path)
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
