package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var clusterSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Install base components into the current cluster",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Installing base components...")
	},
}
