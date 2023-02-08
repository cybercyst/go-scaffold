package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/cybercyst/go-scaffold/pkg/scaffold"
	"github.com/spf13/cobra"
)

// schemaCmd represents the schema command
var schemaCmd = &cobra.Command{
	Use:   "schema [URI]",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		uri := args[0]

		template, err := scaffold.Download(uri)
		if err != nil {
			panic(err)
		}

		templateBytes, err := json.Marshal(template)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(templateBytes))
	},
}

func init() {
	rootCmd.AddCommand(schemaCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// schemaCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// schemaCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
