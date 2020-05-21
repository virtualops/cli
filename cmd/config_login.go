package cmd

import (
	"fmt"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/virtualops/cli/pkg/config"
	"os"
)

var loginCmd = &cobra.Command{
	Use: "login",
	PersistentPreRun: loadGlobalConfig,
	Run: func(cmd *cobra.Command, args []string) {
		if config.GlobalConfig.AuthToken == "" {
			fmt.Printf("%s You are already logged in\n", promptui.IconGood)
			os.Exit(0)
		}
	},
}
