package internal

import (
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
