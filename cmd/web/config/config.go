package config

import (
	"fmt"
	"github.com/usmanzaheer1995/devconnect-go-v2/pkg/utils"
	"gopkg.in/yaml.v2"
	"os"
)

// Config contains the system wide configuration
type Config struct {
	Port     string         `yaml:"port"`
	Env      string         `yaml:"env"`
	Secret   string         `yaml:"secret"`
	Database PostgresConfig `yaml:"postgres"` // Change this in case of different db
}

// DefaultConfig provides default config in development environment
func DefaultConfig() Config {
	return Config{
		Port:     "5000",
		Env:      "dev",
		Database: DefaultPostgresConfig(),
	}
}

// LoadConfig is used to load all the app wide configuration
func LoadConfig(configReq bool) Config {
	f, err := os.Open("./cmd/web/config/config.yaml")
	if err != nil {
		if configReq {
			panic(err)
		}
		fmt.Println("Using the default config...")
		return DefaultConfig()
	}

	var c Config
	dec := yaml.NewDecoder(f)
	err = dec.Decode(&c)
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully loaded config.yaml")
	c.Port = utils.EnvString("PORT", c.Port)
	if err = utils.SetEnvString("SECRET", c.Secret); err != nil {
		panic(err)
	}
	return c
}
