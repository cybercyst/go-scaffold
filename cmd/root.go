package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/cybercyst/go-cookiecutter/internal"
	"github.com/cybercyst/go-cookiecutter/internal/template"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   fmt.Sprintf("%s [TEMPLATE]", internal.ProgramName),
	Short: "A cookiecutter-like templating CLI written in Go",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		templateUri := args[0]

		path, err := template.Download(templateUri)
		if err != nil {
			log.Fatal("Error while preparing template: ", err)
		}

		err = template.Generate(os.DirFS(path), map[string]interface{}{})
		if err != nil {
			log.Fatal("Error while generating template: ", err)
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
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.go-cookiecutter.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(fmt.Sprintf(".%s", internal.ProgramName))
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
