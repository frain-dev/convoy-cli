package main

import (
	"errors"
	"fmt"
	convoyCli "github.com/frain-dev/convoy-cli"
	"strings"

	"github.com/frain-dev/convoy/util"
	"github.com/spf13/cobra"

	"github.com/frain-dev/convoy/pkg/log"
)

func addSwitchCommand() *cobra.Command {
	var appName string
	var projectId string

	cmd := &cobra.Command{
		Use:               "switch",
		Short:             "Switches the current project context",
		PersistentPreRun:  func(cmd *cobra.Command, args []string) {},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := convoyCli.NewConfig("", "")
			if err != nil {
				return err
			}

			if !c.HasDefaultConfigFile() {
				return errors.New("login with your personal access key to be able to use the switch command")
			}

			if util.IsStringEmpty(appName) && util.IsStringEmpty(projectId) {
				return errors.New("one of app name or app id is required")
			}

			var project *convoyCli.ConfigProject
			if !util.IsStringEmpty(appName) {
				project = FindProjectByName(c.Projects, appName)
				if project == nil {
					return fmt.Errorf("app with name: %s not found", appName)
				}
			}

			if !util.IsStringEmpty(projectId) {
				project = FindProjectById(c.Projects, projectId)
				if project == nil {
					return fmt.Errorf("project with id: %s not found", projectId)
				}
			}

			c.ActiveProjectID = project.UID

			err = c.WriteToDisk()
			if err != nil {
				return err
			}

			log.Infof("%s is now the active project", c.ActiveProjectID)
			return nil
		},
	}

	cmd.Flags().StringVar(&projectId, "id", "", "Endpoint Id")

	return cmd
}

func FindProjectByName(endpoints []convoyCli.ConfigProject, endpointName string) *convoyCli.ConfigProject {
	var project *convoyCli.ConfigProject

	for _, endpoint := range endpoints {
		if strings.TrimSpace(strings.ToLower(endpoint.Name)) == strings.TrimSpace(strings.ToLower(endpointName)) {
			return &endpoint
		}
	}

	return project
}

func FindProjectById(projects []convoyCli.ConfigProject, projectId string) *convoyCli.ConfigProject {
	var project *convoyCli.ConfigProject

	for _, endpoint := range projects {
		if strings.TrimSpace(strings.ToLower(endpoint.UID)) == strings.TrimSpace(strings.ToLower(projectId)) {
			return &endpoint
		}
	}

	return project
}
