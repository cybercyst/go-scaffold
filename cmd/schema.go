package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/cybercyst/go-scaffold/pkg/scaffold"
	"github.com/spf13/cobra"
)

var schemaCmd = &cobra.Command{
	Use:   "schema [URI]",
	Short: "Get schema for the provided template",
	Long: `Return the schema for the defined template in JSON format.
This can be used by other tooling in order to provide a rich user
interface to collect the required input.`,
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
}
