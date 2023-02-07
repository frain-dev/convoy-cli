package main

import (
	"bytes"
	"errors"
	"fmt"
	convoyCli "github.com/frain-dev/convoy-cli"
	"os"
	"strings"

	"github.com/frain-dev/convoy/util"
	"github.com/spf13/cobra"

	"github.com/frain-dev/convoy/pkg/log"
)

func addProjectCommand() *cobra.Command {
	var list bool
	var projectId string

	cmd := &cobra.Command{
		Use:               "project",
		Short:             "Switches the current project context or List all projects",
		PersistentPreRun:  func(cmd *cobra.Command, args []string) {},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {},
		RunE: func(cmd *cobra.Command, args []string) error {
			if list {
				return listProjects()
			}

			return switchProject(projectId)
		},
	}

	cmd.LocalFlags().BoolVar(&list, "list", false, "List all projects")
	cmd.LocalFlags().StringVar(&projectId, "switch-to", "", "Switch to specified project")

	return cmd
}

func switchProject(projectId string) error {
	if util.IsStringEmpty(projectId) {
		return errors.New("project id is required")
	}

	c, err := convoyCli.NewConfig("", "")
	if err != nil {
		return err
	}

	if !c.HasDefaultConfigFile() {
		return errors.New("login with your personal access key to be able to use the switch command")
	}

	project := FindProjectById(c.Projects, projectId)
	if project == nil {
		return fmt.Errorf("project with id: %s not found", projectId)
	}

	c.ActiveProjectID = project.UID

	err = c.WriteToDisk()
	if err != nil {
		return err
	}

	log.Infof("Successfully switched to %s", project.Name)
	return nil
}

func listProjects() error {
	c, err := convoyCli.NewConfig("", "")
	if err != nil {
		return err
	}

	buf := bytes.Buffer{}
	for _, p := range c.Projects {
		status := "Inactive"
		if c.ActiveProjectID == p.UID {
			status = "Active"
		}

		v := "- ID: %s\nName: %s\nType: %s\nHost: %s\nStatus: %s\n"
		formated := fmt.Sprintf(v, p.UID, p.Name, p.Type, p.Host, status)

		buf.WriteString(formated)
	}

	buf.WriteTo(os.Stdout)
	return nil
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
