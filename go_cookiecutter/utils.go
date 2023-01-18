package go_cookiecutter

import (
	"errors"
	"fmt"
	"log"
	"os"
)

const ProgramName = "go-cookiecutter"

func ensurePathExists(path string) error {
	info, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
		return ensurePathExists(path)
	}

	if !info.IsDir() {
		return err
	}

	return nil
}

func createTempDir() string {
	path, err := os.MkdirTemp(os.TempDir(), fmt.Sprintf("%s-", ProgramName))
	if err != nil {
		log.Fatal("Unable to create temporary directory: ", err)
	}

	return path
}
