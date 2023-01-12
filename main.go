package main

import (
	"os"

	"github.com/cybercyst/go-cookiecutter/cmd"
)

func cmdMain() int {
	return cmd.Execute()
}

func main() {
	os.Exit(cmdMain())
}
