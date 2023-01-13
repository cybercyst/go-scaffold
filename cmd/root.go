package cmd

import (
	"fmt"
	"log"

	"github.com/cybercyst/go-cookiecutter/internal"
	"github.com/cybercyst/go-cookiecutter/internal/template"
	"github.com/spf13/cobra"
)

var (
	inputFile string
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

		err = template.Generate()
		if err != nil {
			log.Fatal("Unable to generate template", err)
		}

		// fmt.Printf("%+v\n", template)
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
	rootCmd.MarkFlagRequired("input-file")
}
