package main

import (
	cli "github.com/frain-dev/convoy-cli"
	"github.com/spf13/cobra"
)

func addLogoutCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "logout",
		Short:             "Logs out of your Convoy instance",
		SilenceUsage:      true,
		PersistentPreRun:  func(cmd *cobra.Command, args []string) {},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cli.DeleteConfigFile()
		},
	}

	return cmd
}
