package config

import (
	"fmt"
	"os"
)

// Config deals with the stash configuration file ONLY
type Config struct {
	// Name of the configuration file.
	FileName string
}

type ConfigFile struct {
	Destinations []Destination
}

func NewConfig() *Config {
	return &Config{FileName: "$HOME/.stash/config.json"}
}

func (self *Config) AddDestination(dest Destination) {
	fmt.Fprintf(os.Stdout,
		"Adding destination [%s] to config file [%s]\n",
		dest.Name(),
		self.FileName)
}
