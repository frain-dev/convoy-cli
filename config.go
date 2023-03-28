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
	Host                 string          `yaml:"host"`
	ActiveApiKey         string          `yaml:"active_api_key"`
	ActiveProjectID      string          `yaml:"active_project_id"`
	Projects             []ConfigProject `yaml:"projects"`
	path                 string
	hasDefaultConfigFile bool
	isNewHost            bool
}

type ConfigProject struct {
	UID  string `yaml:"uid"`
	Name string `yaml:"name"`
	Host string `yaml:"host"`
	Type string `yaml:"type"`
	//ApiKey   string `yaml:"api_key"`
	DeviceID string `yaml:"device_id"`
}

func DeleteConfigFile() error {
	homedir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	path := filepath.Join(homedir, defaultConfigDir)

	return os.Remove(path)
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
			c.ActiveApiKey = apiKey
		}
		return c, nil
	} else {
		err = os.Mkdir(homedir+"/.convoy", 0777)
		if err != nil && !os.IsExist(err) {
			return nil, fmt.Errorf("failed to create config directory: %v", err)
		}
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

	if err = os.WriteFile(c.path, d, 0644); err != nil {
		return fmt.Errorf("failed to write config to disk: %v", err)
	}

	return nil
}

func (c *Config) HasDefaultConfigFile() bool {
	return c.hasDefaultConfigFile
}

func (c *Config) UpdateConfig(response *LoginResponse, isLogin bool) error {
	if len(response.Projects) > 0 && isLogin {
		c.ActiveProjectID = response.Projects[0].Project.UID
	}

	c.Projects = make([]ConfigProject, 0, len(response.Projects))

	for i := range response.Projects {
		rp := &response.Projects[i]

		c.Projects = append(c.Projects, ConfigProject{
			UID:      rp.Project.UID,
			Name:     rp.Project.Name,
			Host:     c.Host,
			Type:     rp.Project.Type,
			DeviceID: rp.Device.UID,
		})
	}

	//if c.hasDefaultConfigFile {
	//	if c.isNewHost {
	//		// if the host is different from the current host in the config file,
	//		// the data in the config file is overwritten
	//		c.Projects = []ConfigProject{
	//			{
	//				UID:      response.Endpoint.UID,
	//				Name:     name,
	//				ApiKey:   c.ActiveApiKey,
	//				DeviceID: response.Device.UID,
	//			},
	//		}
	//	}
	//
	//	if c.isNewApiKey {
	//		if doesEndpointExist(c, response.Endpoint.UID) {
	//			return fmt.Errorf("endpoint with ID (%s) has been added already", response.Endpoint.UID)
	//		}
	//
	//		// If the api key provided is different from the active api key,
	//		// we append the project returned to the list of projects within the config
	//		c.Projects = append(c.Projects, ConfigProject{
	//			UID:      response.Endpoint.UID,
	//			Name:     name,
	//			ApiKey:   c.ActiveApiKey,
	//			DeviceID: response.Device.UID,
	//		})
	//	}
	//
	//} else {
	//	// Make sure the directory holding our config exists
	//	if err := os.MkdirAll(filepath.Dir(c.path), 0o755); err != nil {
	//		return err
	//	}
	//	c.Projects = []ConfigProject{
	//		{
	//			UID:      response.Endpoint.UID,
	//			Name:     name,
	//			ApiKey:   c.ActiveApiKey,
	//			DeviceID: response.Device.UID,
	//		},
	//	}
	//}

	err := c.WriteToDisk()
	if err != nil {
		return err
	}

	return nil
}

func (c Config) FindPro() {

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
