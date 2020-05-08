package cmd

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	config2 "github.com/virtualops/breeze-cli/pkg/config"
	"io/ioutil"
	"os"
)

var rootCmd = &cobra.Command{
	Use:              "breeze",
	Short:            "Breeze CLI lets you manage and deploy your apps on Breeze",
	Long:             "Breeze CLI lets you manage and deploy your apps on Breeze",
	PersistentPreRun: loadBreezeConfig,
	PreRun: func(cmd *cobra.Command, args []string) {
		// override the persistent pre run on the root
	},
}

var Config = &config2.BreezeConfiguration{}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(
		initCmd,
		deployCmd,
	)
}

func loadBreezeConfig(cmd *cobra.Command, args []string) {
	cwd, err := os.Getwd()

	if err != nil {
		os.Exit(1)
	}

	var content []byte
	// we use a loop to trigger the os.Open call a second time after running init
	for {
		content, err = ioutil.ReadFile(fmt.Sprintf("%s/breeze.yaml", cwd))
		if err == nil {
			break
		}

		prompt := promptui.Prompt{
			Label:     "No breeze configuration was found in the current directory, would you like to create one now",
			IsConfirm: true,
		}
		_, err = prompt.Run()
		if err == nil {
			initCmd.Run(cmd, args)
		} else {
			fmt.Println("You need a Breeze configuration to run this command")
			os.Exit(1)
		}
	}

	err = yaml.Unmarshal(content, Config)

	if err != nil {
		os.Exit(1)
	}
}
