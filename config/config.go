package config

import (
	"fmt"
	"os"

	"github.com/cameronsparr/stash/destination"
)

// Config deals with the stash configuration file ONLY
type Config struct {
	// Name of the configuration file.
	FileName string
}

func NewConfig() *Config {
	return &Config{FileName: "$HOME/.stash/config.json"}
}

func (self *Config) AddDestination(dest destination.Destination) {
	fmt.Fprintf(os.Stdout,
		"Adding destination [%s] to config file [%s]\n",
		dest.Name(),
		self.FileName)
}
