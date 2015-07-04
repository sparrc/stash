package config

import (
	"fmt"
	"os"

	"github.com/cameronsparr/stash/destination"
)

// ConfigFile deals with the stash configuration file ONLY
type ConfigFile struct {
	// Name of the confiuration file.
	FileName string
}

func NewConfigFile() *ConfigFile {
	return &ConfigFile{FileName: "$HOME/.stash/config.json"}
}

func (self *ConfigFile) AddDestination(dest destination.Destination) {
	fmt.Fprintf(os.Stdout,
		"Adding destination [%s] to config file [%s]\n",
		dest.Name(),
		self.FileName)
}
