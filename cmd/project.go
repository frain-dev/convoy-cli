package main

import (
	"bytes"
	"errors"
	"fmt"
	convoyCli "github.com/frain-dev/convoy-cli"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"

	"github.com/frain-dev/convoy/util"
	"github.com/spf13/cobra"
)

func addProjectCommand() *cobra.Command {
	var list bool
	var refresh bool
	var projectId string

	cmd := &cobra.Command{
		Use:               "project",
		Short:             "Switch, List or Refresh projects",
		SilenceUsage:      true,
		PersistentPreRun:  func(cmd *cobra.Command, args []string) {},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {},
		Run: func(cmd *cobra.Command, args []string) {
			if list {
				err := listProjects()
				if err != nil {
					log.Fatal(err)
				}
				return
			}

			if refresh {
				err := login("", "", false)
				if err != nil {
					log.Fatal(err)
				}
				return
			}

			err := switchProject(projectId)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().BoolVar(&list, "list", false, "List all projects")
	cmd.Flags().BoolVar(&refresh, "refresh", false, "Refresh the project list")
	cmd.Flags().StringVar(&projectId, "switch-to", "", "Switch to specified project")

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
		return fmt.Errorf("project with id: %s not found\nRun `convoy-cli project --refresh` to refresh the project list", projectId)
	}

	c.ActiveProjectID = project.UID

	err = c.WriteToDisk()
	if err != nil {
		return err
	}

	fmt.Printf("Successfully switched to %s\n", project.Name)
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

		v := "- ID: %s\n Name: %s\n Type: %s\n Host: %s\n Status: %s\n\n"
		formated := fmt.Sprintf(v, p.UID, p.Name, p.Type, p.Host, status)

		buf.WriteString(formated)
	}

	_, err = buf.WriteTo(os.Stdout)
	return err
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
