package main

import (
	cli "github.com/frain-dev/convoy-cli"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func addLogoutCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "logout",
		Short:             "Logs out of your Convoy instance",
		SilenceUsage:      true,
		PersistentPreRun:  func(cmd *cobra.Command, args []string) {},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {},
		Run: func(cmd *cobra.Command, args []string) {
			err := cli.DeleteConfigFile()
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	return cmd
}
