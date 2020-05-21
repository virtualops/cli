package cmd

import (
	"github.com/spf13/cobra"
	"github.com/virtualops/cli/pkg/config"
)

func loadGlobalConfig(command *cobra.Command, args []string) {
	if err := config.GlobalConfig.Load(); err != nil {
		config.GlobalConfig.Persist()
	}
}
