package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	config2 "github.com/virtualops/breeze-cli/pkg/config"
	"os"
)

var initCmd = &cobra.Command{
	Use: "init",
	PreRun: func(cmd *cobra.Command, args []string) {
		// empty callback to override the global config loader
	},
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		breezeFile := fmt.Sprintf("%s/%s", cwd, config2.DefaultFileName)
		if err != nil {
			os.Exit(1)
		}

		_, err = os.Stat(breezeFile)

		if err == nil {
			fmt.Println("Breeze configuration already exists")
			os.Exit(1)
		}

		file, err := os.Create(breezeFile)

		if err != nil {
			os.Exit(1)
		}

		_, err = file.WriteString(config2.DefaultConfigFile)

		if err != nil {
			fmt.Println("Failed to create new config file")
			os.Exit(1)
		}

		file.Close()
	},
}
