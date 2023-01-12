package main

import (
	"os"
	"testing"

	"github.com/cybercyst/go-cookiecutter/cmd"
	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	os.Exit(testscript.RunMain(m, map[string]func() int{
		"go-cookiecutter": cmd.Execute(),
	}))
}
