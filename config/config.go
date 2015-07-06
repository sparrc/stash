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

func (cm *ConfigMngr) createConfigFile() {
	log.Println("Creating config file")
	args := []string{"-p", filepath.Dir(cm.FileName)}
	if err := exec.Command("mkdir", args...).Run(); err != nil {
		panic(err)
	}
	if err := exec.Command("touch", cm.FileName).Run(); err != nil {
		panic(err)
	}
}

func (cm *ConfigMngr) AddDestination(configEntry ConfigEntry) {
	fmt.Fprintf(os.Stdout,
		"Adding destination [%s] to config file [%s]\n",
		configEntry.Name,
		cm.FileName)

	allConfigEntries := cm.GetNewConfigEntries(configEntry)
	JSON := cm.JSONMarshallEntry(allConfigEntries)
	ioutil.WriteFile(cm.FileName, JSON, 0644)
}

// This function takes a new config entry, loads previous entries, and combines
func (cm *ConfigMngr) GetNewConfigEntries(configEntry ConfigEntry) []ConfigEntry {
	configFile := cm.LoadConfigFile()
	configFile = append(configFile, configEntry)
	return configFile
}

func (cm *ConfigMngr) JSONMarshallEntry(configEntries []ConfigEntry) []byte {
	JSON, err := json.MarshalIndent(configEntries, "", "  ")
	if err != nil {
		panic(err)
	}
	return JSON
}

func (cm *ConfigMngr) LoadConfigFile() []ConfigEntry {
	fmt.Fprintf(os.Stdout, "Loading config file [%s]\n", cm.FileName)
	content, err := ioutil.ReadFile(cm.FileName)
	if err != nil {
		panic(err)
	}
	var entries []ConfigEntry
	if err := json.Unmarshal(content, &entries); err != nil {
		log.Println("No config entries loaded.")
	}
	return entries
}
