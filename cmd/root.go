package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bclicn/color"
	"github.com/cybercyst/go-scaffold/internal/consts"
	"github.com/cybercyst/go-scaffold/internal/utils"
	"github.com/cybercyst/go-scaffold/pkg/scaffold"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
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

		template, err := scaffold.Download(uri)
		if err != nil {
			fmt.Fprint(os.Stderr, color.Red("[ERROR]"), " problem preparing template: ", err)
			os.Exit(1)
		}

		metadata, err := scaffold.Generate(template, &input, outputDirectory)
		if err != nil {
			fmt.Fprint(os.Stderr, color.Red("[ERROR]"), " problem generating template: ", err)
			os.Exit(1)
		}

		metadataYaml, err := yaml.Marshal(metadata)
		if err != nil {
			fmt.Fprint(os.Stderr, color.Red("[ERROR]"), " unable to parse generated artifact metadata: ", err)
			os.Exit(1)
		}
		err = os.WriteFile(filepath.Join(outputDirectory, ".metadata.yaml"), metadataYaml, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, color.Red("[ERROR]"), " unable to save generated artifact metadata: ", err)
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
