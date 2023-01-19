package cmd

import (
	"fmt"
	"os"

	"github.com/bclicn/color"
	"github.com/cybercyst/go-cookiecutter/internal/consts"
	"github.com/cybercyst/go-cookiecutter/internal/utils"
	"github.com/cybercyst/go-cookiecutter/pkg/cookiecutter"
	"github.com/spf13/cobra"
)

var (
	inputFile       string
	outputDirectory string
)

var rootCmd = &cobra.Command{
	Use:   fmt.Sprintf("%s [TEMPLATE]", consts.ProgramName),
	Short: "A cookiecutter-like templating CLI written in Go",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		uri := args[0]

		input, err := utils.ReadTemplateInput(inputFile)
		if err != nil {
			fmt.Fprint(os.Stderr, color.Red("[ERROR]"), " problem reading input: ", err)
			os.Exit(1)
		}

		fmt.Println()
		fmt.Println("Generating template", color.BBlue(uri))

		metadata, err := cookiecutter.Generate(uri, &input, outputDirectory)
		if err != nil {
			fmt.Fprint(os.Stderr, color.Red("[ERROR]"), " problem generating template: ", err)
			os.Exit(1)
		}

		fmt.Println()
		for _, file := range *metadata.CreatedFiles {
			fmt.Println(color.Green("CREATE\t"), file)
		}

		fmt.Println()
		fmt.Println("Template", color.BBlue(uri), "successfully generated")
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
