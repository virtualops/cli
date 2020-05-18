package cmd

import "github.com/spf13/cobra"

var shellCmd = &cobra.Command{
	Use:     "shell",
	Short:   "Start an interactive shell session for the application",
	Long:    "Start an interactive shell session for the current application, a shorthand for breeze run bash",
	Aliases: []string{"sh", "bash"},
	Run: func(cmd *cobra.Command, args []string) {
		runCmd.Run(cmd, []string{"bash"})
	},
}
