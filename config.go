package convoy_cli

import (
	"errors"
	"fmt"
	"github.com/frain-dev/convoy-cli/util"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

const (
	defaultConfigDir = ".convoy/config"
)

type Config struct {
	Host                 string           `yaml:"host"`
	ActiveDeviceID       string           `yaml:"active_device_id"`
	ActiveApiKey         string           `yaml:"active_api_key"`
	ActiveEndpoint       string           `yaml:"active_endpoint"`
	Endpoints            []ConfigEndpoint `yaml:"endpoints"`
	path                 string
	hasDefaultConfigFile bool
	isNewApiKey          bool
	isNewHost            bool
}

type ConfigEndpoint struct {
	UID      string `yaml:"uid"`
	Name     string `yaml:"name"`
	ApiKey   string `yaml:"api_key"`
	DeviceID string `yaml:"device_id"`
}

func LoadConfig() (*Config, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(homedir, defaultConfigDir)

	c := &Config{path: path}
	c.hasDefaultConfigFile = HasDefaultConfigFile(path)

	if !c.hasDefaultConfigFile {
		return nil, errors.New("config file not found")
	}

	if c.hasDefaultConfigFile {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		err = yaml.Unmarshal(data, &c)
		if err != nil {
			return nil, err
		}

		return c, nil
	}

	return nil, nil
}

func NewConfig(host, apiKey string) (*Config, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(homedir, defaultConfigDir)

	c := &Config{path: path}
	c.hasDefaultConfigFile = HasDefaultConfigFile(path)

	if c.hasDefaultConfigFile {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		err = yaml.Unmarshal(data, &c)
		if err != nil {
			return nil, err
		}

		if !util.IsStringEmpty(host) {
			c.isNewHost = IsNewHost(c.Host, host)
			c.Host = host
		}

		if !util.IsStringEmpty(apiKey) {
			c.isNewApiKey = IsNewApiKey(c, apiKey)
			c.ActiveApiKey = apiKey
		}
		return c, nil
	}

	c.Host = host
	c.ActiveApiKey = apiKey

	return c, nil
}

func (c *Config) WriteToDisk() error {
	d, err := yaml.Marshal(&c)
	if err != nil {
		return err
	}

	if err := os.WriteFile(c.path, []byte(d), 0o644); err != nil {
		return err
	}

	return nil
}

func (c *Config) HasDefaultConfigFile() bool {
	return c.hasDefaultConfigFile
}

func (c *Config) UpdateConfig(response *LoginResponse) error {
	name := fmt.Sprintf("%s (%s)", response.Endpoint.Title, response.Project.Name)
	c.ActiveEndpoint = name
	c.ActiveDeviceID = response.Device.UID

	if c.hasDefaultConfigFile {
		if c.isNewHost {
			// if the host is different from the current host in the config file,
			// the data in the config file is overwritten
			c.Endpoints = []ConfigEndpoint{
				{
					UID:      response.Endpoint.UID,
					Name:     name,
					ApiKey:   c.ActiveApiKey,
					DeviceID: response.Device.UID,
				},
			}
		}

		if c.isNewApiKey {
			if doesEndpointExist(c, response.Endpoint.UID) {
				return fmt.Errorf("endpoint with ID (%s) has been added already", response.Endpoint.UID)
			}

			// If the api key provided is different from the active api key,
			// we append the project returned to the list of projects within the config
			c.Endpoints = append(c.Endpoints, ConfigEndpoint{
				UID:      response.Endpoint.UID,
				Name:     name,
				ApiKey:   c.ActiveApiKey,
				DeviceID: response.Device.UID,
			})
		}

	} else {
		// Make sure the directory holding our config exists
		if err := os.MkdirAll(filepath.Dir(c.path), 0o755); err != nil {
			return err
		}
		c.Endpoints = []ConfigEndpoint{
			{
				UID:      response.Endpoint.UID,
				Name:     name,
				ApiKey:   c.ActiveApiKey,
				DeviceID: response.Device.UID,
			},
		}
	}

	err := c.WriteToDisk()
	if err != nil {
		return err
	}

	return nil
}

func doesEndpointExist(c *Config, endpointId string) bool {
	for _, endpoint := range c.Endpoints {
		if endpoint.UID == endpointId {
			return true
		}
	}

	return false
}

func HasDefaultConfigFile(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}

func IsNewHost(currentHost, newHost string) bool {
	return currentHost != newHost
}

// The api key is considered new if it doesn't already
// exist within the config file
func IsNewApiKey(c *Config, apiKey string) bool {
	for _, project := range c.Endpoints {
		if project.ApiKey == apiKey {
			return false
		}
	}

	return true
}
