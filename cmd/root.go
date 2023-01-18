package cmd

import (
	"fmt"
	"os"

	"github.com/cybercyst/go-cookiecutter/go_cookiecutter"
	"github.com/cybercyst/go-cookiecutter/internal/utils"
	"github.com/spf13/cobra"
)

var (
	inputFile       string
	outputDirectory string
)

const ProgramName = "go-cookiecutter"

var rootCmd = &cobra.Command{
	Use:   fmt.Sprintf("%s [TEMPLATE]", ProgramName),
	Short: "A cookiecutter-like templating CLI written in Go",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		uri := args[0]

		input, err := utils.ReadTemplateInput(inputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: problem reading input: %s", err)
			os.Exit(1)
		}

		// TODO: get metadata and write to file
		_, err = go_cookiecutter.Generate(uri, &input, outputDirectory)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: problem generating template: %s", err)
			os.Exit(1)
		}
	},
}

func Execute() int {
	err := rootCmd.Execute()
	if err != nil {
		return 1
	}

	return 0
}

func init() {
	rootCmd.Flags().StringVarP(&inputFile, "input-file", "i", "", "File containing variables used as input for the template")
	rootCmd.Flags().StringVarP(&outputDirectory, "output-directory", "o", ".", "Directory where template will be generated")
	rootCmd.MarkFlagRequired("input-file")
}
