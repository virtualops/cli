package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use: "deploy",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Unfortunately, Breeze isn't yet ready for production releases.")
		fmt.Println("You can sign up for updates on Breeze for production at https://breeze.sh")
	},
}
