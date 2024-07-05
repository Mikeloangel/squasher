// Package config provides an app configuration
package config

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Config has all needed config fields for an app
type Config struct {
	ServerLocation string
	ServerPort     int
	HostLocation   string

	StorageFileLocation string

	DBDSN string
}

// NewConfig creates a new instance of Config
func NewConfig(
	serverLocation string,
	serverPort int,
	hostLocation string,
	storageFileLocation string,
	DBDSN string,
) *Config {
	return &Config{
		ServerLocation:      serverLocation,
		ServerPort:          serverPort,
		HostLocation:        hostLocation,
		StorageFileLocation: storageFileLocation,
		DBDSN:               DBDSN,
	}
}

// GetServerConnectionString return string of server connection host
func (c Config) GetServerConnectionString() string {
	return fmt.Sprintf("%s:%v", c.ServerLocation, c.ServerPort)
}

// GetHostLocation returns a string og host location
func (c Config) GetHostLocation() string {
	return fmt.Sprintf("%s/", c.HostLocation)
}

// ParseServerConfig parses server config for host:port string into
// ServerLocation and ServerPort variable
// return error if param was incorrect or port is not number
func (c *Config) ParseServerConfig(s string) error {
	var err error

	host := strings.Split(s, ":")

	if len(host) != 2 {
		return errors.New("bad format for server string. Expected format host:port")
	}

	c.ServerLocation = host[0]
	c.ServerPort, err = strconv.Atoi(host[1])

	if err != nil {
		return err
	}

	return nil
}
