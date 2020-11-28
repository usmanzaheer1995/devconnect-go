package config

import (
	"fmt"
)

// PostgresConfig stores all the postgres config
type PostgresConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password,omitempty"`
	Name     string `yaml:"name"`
}

// Dialect returns the dialect used by Gorm
func (c PostgresConfig) Dialect() string {
	return "postgres"
}

// ConnectionInfo returns the connection string for gorm connection to postgres
func (c PostgresConfig) ConnectionInfo() string {
	if c.Password == "" {
		return fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable",
			c.Host, c.Port, c.User, c.Name,
		)
	}
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.Name,
	)
}

// DefaultPostgresConfig returns a default postgres config in case none is provided
func DefaultPostgresConfig() PostgresConfig {
	return PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "embrace123",
		Name:     "devconnect_dev",
	}
}
