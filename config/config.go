package config

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Config struct {
	ServerLocation string
	ServerPort     int
	HostLocation   string
}

func NewConfig() *Config {
	return &Config{
		ServerLocation: "localhost",
		ServerPort:     8080,
		HostLocation:   "http://localhost:8080",
	}
}

func (c Config) GetServerConnectionString() string {
	return fmt.Sprintf("%s:%v", c.ServerLocation, c.ServerPort)
}

func (c Config) GetHostLocation() string {
	return fmt.Sprintf("%s/", c.HostLocation)
}

func (c *Config) SetServerFromString(s string) error {
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
