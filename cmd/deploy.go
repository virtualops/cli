package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use: "deploy",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("running deploy")
	},
}
