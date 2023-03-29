package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"os"
	"time"

	convoyCli "github.com/frain-dev/convoy-cli"
	convoyNet "github.com/frain-dev/convoy-cli/net"
	"github.com/frain-dev/convoy-cli/util"
	"github.com/spf13/cobra"
)

func addLoginCommand() *cobra.Command {
	var apiKey string
	var host string

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Logs into your Convoy instance using a Personal API Key",
		Run: func(cmd *cobra.Command, args []string) {
			err := login(host, apiKey, true)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().StringVar(&apiKey, "api-key", "", "API Key")
	cmd.Flags().StringVar(&host, "host", "https://cli.getconvoy.io", "Host")

	return cmd
}

func login(host, apiKey string, isLogin bool) error {
	c, err := convoyCli.NewConfig(host, apiKey)
	if err != nil {
		return err
	}

	if util.IsStringEmpty(c.Host) {
		return errors.New("host is required")
	}

	if util.IsStringEmpty(c.ActiveApiKey) {
		return errors.New("api key is required")
	}

	hostName, err := generateDeviceHostName()
	if err != nil {
		return err
	}

	loginRequest := &convoyCli.LoginRequest{HostName: hostName}
	body, err := json.Marshal(loginRequest)
	if err != nil {
		return err
	}

	var response *convoyCli.LoginResponse

	dispatch, err := convoyNet.NewDispatcher(time.Second*10, "")
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/stream/login", c.Host)
	resp, err := dispatch.SendCliRequest(url, http.MethodPost, c.ActiveApiKey, body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New(string(resp.Body))
	}

	err = json.Unmarshal(resp.Body, &response)
	if err != nil {
		return err
	}

	err = c.UpdateConfig(response, isLogin)
	if err != nil {
		return err
	}

	if isLogin {
		fmt.Println("Login Success!")
		fmt.Println("Name:", response.UserName)
		fmt.Println("Host:", host)
		return nil
	}

	fmt.Println("Refresh Project list Successful!")
	return nil
}

// generateDeviceHostName uses the machine's host name and the mac address to generate a predictable unique id per device
func generateDeviceHostName() (string, error) {
	name, err := os.Hostname()
	if err != nil {
		return "", err
	}

	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	var mac uint64
	for _, i := range interfaces {
		if i.Flags&net.FlagUp != 0 && !bytes.Equal(i.HardwareAddr, nil) {

			// Skip virtual MAC addresses (Locally Administered Addresses).
			if i.HardwareAddr[0]&2 == 2 {
				continue
			}

			for j, b := range i.HardwareAddr {
				if j >= 8 {
					break
				}
				mac <<= 8
				mac += uint64(b)
			}
		}
	}

	return fmt.Sprintf("%v-%v", name, mac), nil
}
