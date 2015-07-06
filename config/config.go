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

// Config deals with the stash configuration file
type ConfigMngr struct {
	// Name of the configuration file.
	FileName string
}

type ConfigEntry struct {
	Name        string
	Folders     []string
	Type        string
	Credentials map[string]string
}

func NewConfigMngr() *ConfigMngr {
	filename := filepath.Join(os.Getenv("HOME"), ".stash", "config.json")
	config := ConfigMngr{FileName: filename}
	// Create config file if it doesn't exist:
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		config.createConfigFile()
	}
	return &config
}

func (self *ConfigMngr) createConfigFile() {
	log.Println("Creating config file")
	args := []string{"-p", filepath.Dir(self.FileName)}
	if err := exec.Command("mkdir", args...).Run(); err != nil {
		panic(err)
	}
	if err := exec.Command("touch", self.FileName).Run(); err != nil {
		panic(err)
	}
}

func (self *ConfigMngr) AddDestination(configEntry ConfigEntry) {
	// TODO: implement
	fmt.Fprintf(os.Stdout,
		"Adding destination [%s] to config file [%s]\n",
		configEntry.Name,
		self.FileName)

	jsonEntry, err := json.MarshalIndent(configEntry, "", "  ")
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile(self.FileName, jsonEntry, 0644)
}

func (self *ConfigMngr) LoadConfigFile() []ConfigEntry {
	fmt.Fprintf(os.Stdout, "Loading config file [%s]\n", self.FileName)
	content, err := ioutil.ReadFile(self.FileName)
	if err != nil {
		panic(err)
	}
	var entries []ConfigEntry
	if err := json.Unmarshal(content, &entries); err != nil {
		panic(err)
	}
	return entries
}
