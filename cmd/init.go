package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	config2 "github.com/virtualops/breeze-cli/pkg/config"
	"os"
)

var initCmd = &cobra.Command{
	Use: "init",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// empty callback to override the global Config loader
	},
	Run: func(cmd *cobra.Command, args []string) {
		cwd, err := os.Getwd()
		if err != nil {
			os.Exit(1)
		}

		breezeFile := fmt.Sprintf("%s/%s", cwd, config2.DefaultFileName)

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
			fmt.Println("Failed to create new Config file")
			os.Exit(1)
		}

		fmt.Println("Created Breeze configuration file at breeze.yaml")

		file.Close()
	},
}
