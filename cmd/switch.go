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
	var endpointId string

	cmd := &cobra.Command{
		Use:               "switch",
		Short:             "Switches the current application context",
		PersistentPreRun:  func(cmd *cobra.Command, args []string) {},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := convoyCli.NewConfig("", "")
			if err != nil {
				return err
			}

			if !c.HasDefaultConfigFile() {
				return errors.New("login with your cli token to be able to use the switch command")
			}

			if util.IsStringEmpty(appName) && util.IsStringEmpty(endpointId) {
				return errors.New("one of app name or app id is required")
			}

			var endpoint *convoyCli.ConfigEndpoint
			if !util.IsStringEmpty(appName) {
				endpoint = FindEndpointByName(c.Endpoints, appName)
				if endpoint == nil {
					return fmt.Errorf("app with name: %s not found", appName)
				}
			}

			if !util.IsStringEmpty(endpointId) {
				endpoint = FindEndpointById(c.Endpoints, endpointId)
				if endpoint == nil {
					return fmt.Errorf("endpoint with id: %s not found", endpointId)
				}
			}

			c.ActiveEndpoint = endpoint.Name
			c.ActiveDeviceID = endpoint.DeviceID
			c.ActiveApiKey = endpoint.ApiKey

			err = c.WriteToDisk()
			if err != nil {
				return err
			}

			log.Infof("%s is now the active endpoint", c.ActiveEndpoint)
			return nil
		},
	}

	cmd.Flags().StringVar(&appName, "name", "", "Endpoint Name")
	cmd.Flags().StringVar(&endpointId, "id", "", "Endpoint Id")

	return cmd
}

func FindEndpointByName(endpoints []convoyCli.ConfigEndpoint, endpointName string) *convoyCli.ConfigEndpoint {
	var endpoint *convoyCli.ConfigEndpoint

	for _, endpoint := range endpoints {
		if strings.TrimSpace(strings.ToLower(endpoint.Name)) == strings.TrimSpace(strings.ToLower(endpointName)) {
			return &endpoint
		}
	}

	return endpoint
}

func FindEndpointById(endpoints []convoyCli.ConfigEndpoint, endpointId string) *convoyCli.ConfigEndpoint {
	var endpoint *convoyCli.ConfigEndpoint

	for _, endpoint := range endpoints {
		if strings.TrimSpace(strings.ToLower(endpoint.UID)) == strings.TrimSpace(strings.ToLower(endpointId)) {
			return &endpoint
		}
	}

	return endpoint
}
