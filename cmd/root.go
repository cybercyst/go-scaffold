package cmd

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/cybercyst/go-cookiecutter/internal"
	"github.com/cybercyst/go-cookiecutter/internal/template"
	"github.com/spf13/cobra"
)

var (
	inputFile       string
	outputDirectory string
)

var rootCmd = &cobra.Command{
	Use:   fmt.Sprintf("%s [TEMPLATE]", internal.ProgramName),
	Short: "A cookiecutter-like templating CLI written in Go",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		template := &template.Template{}

		uri := args[0]
		err := template.Download(uri)
		if err != nil {
			log.Fatal("Error while preparing template: ", err)
		}

		input, err := internal.ReadTemplateInput(inputFile)
		if err != nil {
			log.Fatal("Error while getting template input: ", err)
		}
		template.Input = input

		err = template.ValidateInput()
		if err != nil {
			log.Fatal("Error while validating input: ", err)
		}

		absOutputDirectory, err := filepath.Abs(outputDirectory)
		if err != nil {
			log.Fatal("Error setting output directory: ", err)
		}
		template.OutputPath = absOutputDirectory

		err = template.Generate()
		if err != nil {
			log.Fatal("Unable to generate template: ", err)
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
