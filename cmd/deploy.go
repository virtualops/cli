package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use: "deploy",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// for now we can ignore the breeze.yaml check since the command doesn't exist
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Unfortunately, Breeze isn't yet ready for production releases.")
		fmt.Println("You can sign up for updates on Breeze for production at https://breeze.sh")
	},
}
