package main

import (
	"os"

	convoyCli "github.com/frain-dev/convoy-cli"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	err := os.Setenv("TZ", "") // Use UTC by default :)
	if err != nil {
		log.Fatal("failed to set TZ env - ", err)
	}

	cmd := &cobra.Command{
		Use:     "convoy-cli",
		Version: convoyCli.GetVersion(),
		Short:   "Client CLI for debugging your events locally",
	}

	cmd.AddCommand(addListenCommand())
	cmd.AddCommand(addLoginCommand())
	cmd.AddCommand(addProjectCommand())
	cmd.AddCommand(addLogoutCommand())
	cmd.AddCommand(addStatusCommand())

	err = cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
