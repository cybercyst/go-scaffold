package internal

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
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

func CreateTempDir() string {
	path, err := os.MkdirTemp(os.TempDir(), fmt.Sprintf("%s-", ProgramName))
	if err != nil {
		log.Fatal("Unable to create temporary directory: ", err)
	}

	return path
}
