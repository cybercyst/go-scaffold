package main

import (
	"os"

	"github.com/cybercyst/go-scaffold/cmd"
)

func cmdMain() int {
	return cmd.Execute()
}

func main() {
	os.Exit(cmdMain())
}
