package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

const fileName = "~/.stash/config.json"

// Config deals with the stash configuration file
type Config struct {
	// Name of the configuration file.
	FileName string
}

type ConfigEntry struct {
	Name        string
	Folders     []string
	Type        string
	Credentials map[string]string
}

func NewConfig() *Config {
	config := Config{FileName: fileName}
	// Create config file if it doesn't exist:
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		config.createConfigFile()
	}
	return &config
}

func (self *Config) createConfigFile() {
	log.Println("Creating config file")
	args := []string{"-p", filepath.Dir(self.FileName)}
	if err := exec.Command("mkdir", args...).Run(); err != nil {
		panic(err)
	}
	if err := exec.Command("touch", self.FileName).Run(); err != nil {
		panic(err)
	}
}

func (self *Config) AddDestination(dest Destination) {
	fmt.Fprintf(os.Stdout,
		"Adding destination [%s] to config file [%s]\n",
		dest.Name(),
		self.FileName)
}

func (self *Config) LoadConfigFile() {
	fmt.Fprintf(os.Stdout, "Loading config file [%s]\n", self.FileName)
	content, err := ioutil.ReadFile(self.FileName)
	if err != nil {
		panic(err)
	}
	var entry ConfigEntry
	if err := json.Unmarshal(content, &entry); err != nil {
		panic(err)
	}
}
