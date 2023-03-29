package main

import (
	"fmt"
	convoyCli "github.com/frain-dev/convoy-cli"
	"github.com/frain-dev/convoy-cli/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func addStatusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Checks status of the cli login",
		Run: func(cmd *cobra.Command, args []string) {
			c, err := convoyCli.LoadConfig()
			if err != nil {
				log.Fatal("Error loading config file:", err)
			}

			hasAPIKey := !util.IsStringEmpty(c.ActiveApiKey)
			hasProjects := len(c.Projects) > 0

			if hasProjects && hasAPIKey {
				fmt.Printf("You are logged in with %d Projects\n", len(c.Projects))
				p := FindProjectById(c.Projects, c.ActiveProjectID)

				if p == nil {
					fmt.Println("You have no active project\nRun `convoy-cli project --list` to list your projects")
					fmt.Println("Then run `convoy-cli project --switch-to {project_id}` to choose a project")
				} else {
					fmt.Printf("Your current active project is %s\n", p.Name)
				}

				return
			}

			fmt.Println("You are not logged in, run `convoy-cli login --api-key {api-key} --host {host}` to login")

		},
	}

	return cmd
}
