package config

import (
	"os"
	"os/user"

	"gopkg.in/yaml.v3"
)

const (
	SYSTEM_CONFIG_PATH string = "/etc/confman.conf"
	HOME_CONFIG_PATH   string = ".config/confman.conf"
)

type Config struct {
	Paths map[string]string
}

func (c *Config) AddPath(source, name string) error {
	c.Paths[source] = name

	return c.Save()
}

func (c *Config) RemovePath(source string) error {
	delete(c.Paths, source)

	return c.Save()
}

func (c *Config) Save() error {
	output, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	info, err := os.Stat("/confman/.confman.yaml")
	if err != nil {
		return err
	}

	return os.WriteFile("/confman/.confman.yaml", output, info.Mode())
}

func NewConfig() *Config {
	return &Config{
		Paths: make(map[string]string),
	}
}

func GetConfigFrom(path string) (*Config, error) {
	// Read the yaml file
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := NewConfig()

	err = yaml.Unmarshal(file, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func GetConfigFile() (string, error) {
	user, err := user.Current()
	if err != nil {
		return "", err
	}

	if user.Username == "root" {
		return SYSTEM_CONFIG_PATH, nil
	}

	return user.HomeDir + "/" + HOME_CONFIG_PATH, nil
}
