package config

import "fmt"

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
