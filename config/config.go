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

func (self *ConfigFile) SetDefaults() {
	self.FileName = "$HOME/.stash/config.json"
}

func (self *ConfigFile) AddDestination(dest destination.Destination) {
	fmt.Println(dest.Name())
}
